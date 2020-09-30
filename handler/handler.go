package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

/*
UploadHandler: 文件上传
w: 向用户返回数据的http.ResponseWriter对象
r: 用于接收用户请求的http.Request对象指针
*/
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")

		if err != nil { //加载失败
			io.WriteString(w, "internel server error")
			return
		}

		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		//接收文件流并存储到本地目录
		file, head, err := r.FormFile("file") //返回文件句柄、文件头、错误信息
		if err != nil {
			fmt.Printf("Failed to get data, err:%s", err.Error())
			return
		}
		defer file.Close()

		//创建一个本地的文件来接收文件流
		newFile, err := os.Create("/tmp/" + head.Filename) //创建文件
		if err != nil {
			fmt.Printf("Failed to create file, err:%s", err.Error())
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file) //文件复制，返回复制的字节长度、错误信息
		if err != nil {
			fmt.Printf("Failed to save data into file, err:%s", err.Error())
			return
		}

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound) //上传完成，重定向到上传成功的网页url
	}
}

//UploadSuchHandler: 上传已完成
func UploadSuchHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished")
}
