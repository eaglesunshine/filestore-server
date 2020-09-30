package main

import (
	"filestore-server/handler"
	"fmt"
	"net/http"
)

func main() {
	//网页js、css等静态资源服务
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload", handler.UploadHandler) //设定路由规则
	http.HandleFunc("/file/upload/suc", handler.UploadSuchHandler)
	err := http.ListenAndServe("127.0.0.1:8888", nil) //设定监听端口，这里设置为监听所有网卡的8080端口
	if err != nil {
		fmt.Printf("Failed to start server, err:%s", err.Error())
	}
}
