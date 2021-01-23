package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/abulo/ratel/util"
)

func main() {
	var width string
	var high string
	fmt.Println("请图片输入宽度：")
	fmt.Scan(&width)
	fmt.Println("请图片输入高度：")
	fmt.Scan(&high)

	// 获取文件列表
	list := getFilelist("./")

	for _, v := range list {
		fileExt := path.Ext(v)
		if fileExt == ".mp4" || fileExt == ".ts" || fileExt == ".m3u8" || fileExt == ".mkv" {
			dir, _ := os.Getwd()
			videoPath := dir + "\\" + v
			paths, fileName := filepath.Split(videoPath)
			filesuffix := path.Ext(fileName)
			fileprefix := fileName[0 : len(fileName)-len(filesuffix)]
			pictureSavePath := paths + util.ZhCharToFirstPinyin(fileprefix) + ".jpg"

			videoLen, _ := GenerateLength(videoPath)
			gt := strconv.Itoa(int(math.Ceil(float64(videoLen) / 3)))
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(50000)*time.Millisecond)
			cmd := exec.CommandContext(ctx, "ffmpeg",
				"-y",
				"-ss", gt,
				"-t", "1",
				"-i", videoPath,
				// "-s", "122x92",
				"-s", width+"x"+high,
				"-vframes", "1",
				pictureSavePath)
			defer cancel()
			var buffer bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Stdout = &buffer
			if cmd.Run() != nil {
				panic("could not generate frame")
			}

			// 日志
			Logfile(pictureSavePath)
			fmt.Println(pictureSavePath)
		}
	}
	fmt.Println("生成结束！")
	fmt.Scanf(" ")

}

// func getDuration(filePath string) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
// 	cmd := exec.CommandContext(ctx, "ffprobe", "-of", "json", "-show_streams", filePath)
// 	defer cancel()
// 	var buffer bytes.Buffer
// 	var stderr bytes.Buffer
// 	cmd.Stderr = &stderr
// 	cmd.Stdout = &buffer
// 	if cmd.Run() != nil {
// 		panic("could not generate frame")
// 	}
// 	fmt.Printf("%v", stderr.String())
// 	var re interface{}
// 	err := json.Unmarshal([]byte(stderr.String()), &re)
// 	if err != nil {
// 		fmt.Println(re)
// 	}

// 	fmt.Println("111")
// }

//GenerateLength 获取视频时长 秒
func GenerateLength(filename string) (int, error) {
	var length int
	videoLengthRegexp, _ := regexp.Compile(`Duration: (.*?),`)
	for i := 0; i < 2; i++ {
		//视频处理使用，延长超时时间
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		cmd := exec.CommandContext(ctx, "ffmpeg", "-i", filename)
		defer cancel()
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Run()
		str := videoLengthRegexp.FindString(stderr.String())
		str = "2006-01-02" + strings.TrimSuffix(strings.TrimPrefix(str, "Duration:"), ",")
		if videotime, err := time.Parse("2006-01-02 15:04:05", str); err != nil {
			if ctx.Err() != nil {
				fmt.Println("GenerateLength Err:", ctx.Err())
			}
			length = 0
		} else {
			length = videotime.Hour()*3600 + videotime.Minute()*60 + videotime.Second()
			break
		}
	}
	return length, nil
}

// getFilelist 获取给定目录下所有文件
func getFilelist(path string) []string {
	fileList := make([]string, 0)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		fileList = append(fileList, path)
	}
	return fileList
}

func checkFile(Filename string) bool {
	var exist = true
	if _, err := os.Stat(Filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// Logfile 写入文件
func Logfile(Log string) {
	var f1 *os.File
	var err1 error

	Filenames := "./log.log" //也可将name作为参数传进来

	if checkFile(Filenames) { //如果文件存在
		f1, _ = os.OpenFile(Filenames, os.O_APPEND|os.O_WRONLY, 0666) //打开文件,第二个参数是写入方式和权限

	} else {
		f1, _ = os.Create(Filenames) //创建文件

	}
	_, err1 = io.WriteString(f1, Log+"\n") //写入文件(字符串)
	if err1 != nil {
		fmt.Println(err1)
	}
	//fmt.Printf("写入 %d 个字节\n", n)

	return
}
