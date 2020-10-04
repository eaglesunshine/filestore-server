package handler

import (
	cfg "filestore-server/config"
	dblayer "filestore-server/db"
	"filestore-server/meta"
	"filestore-server/store/ceph"
	"filestore-server/store/oss"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// DownloadURLHandler : 生成文件的下载地址
func DownloadURLHandler(w http.ResponseWriter, r *http.Request) {
	filehash := r.Form.Get("filehash")
	// 从文件表查找记录
	row, _ := dblayer.GetFileMeta(filehash)

	fmt.Println("fileAddr: " + row.FileAddr.String)
	// 判断文件存在OSS，还是Ceph，还是在本地
	if strings.HasPrefix(row.FileAddr.String, cfg.MergeLocalRootDir) ||
		strings.HasPrefix(row.FileAddr.String, "/ceph") {
		username := r.Form.Get("username")
		token := r.Form.Get("token")
		tmpURL := fmt.Sprintf(
			"http://%s/file/download?filehash=%s&username=%s&token=%s",
			r.Host,
			filehash,
			username,
			token)
		w.Write([]byte(tmpURL))
	} else if strings.HasPrefix(row.FileAddr.String, "oss/") {
		// oss下载url
		signedURL := oss.DownloadURL(row.FileAddr.String)
		w.Write([]byte(signedURL))
	} else {
		w.Write([]byte("Error: 下载链接暂时无法生成"))
	}
}

// DownloadHandler : 文件下载接口
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	username := r.Form.Get("username")

	fm, _ := meta.GetFileMetaDB(fsha1)
	userFile, err := dblayer.QueryUserFileMeta(username, fsha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var fileData []byte
	if strings.HasPrefix(fm.Location, cfg.MergeLocalRootDir) {
		f, err := os.Open(fm.Location)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		fileData, err = ioutil.ReadAll(f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(fm.Location, "/ceph") {
		fmt.Println("to download file from ceph...")
		bucket := ceph.GetCephBucket("userfile")
		fileData, err = bucket.Get(fm.Location)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if strings.HasPrefix(fm.Location, "oss") {
		fmt.Println("to download file from oss...")
		// TODO: to verify the code in this block
		fd, err := oss.Bucket().GetObject(fm.Location)
		if err == nil {
			fileData, err = ioutil.ReadAll(fd)
		}
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.Write([]byte("File not found."))
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+userFile.FileName+"\"")
	w.Write(fileData)
}

// RangeDownloadHandler : 支持断点的文件下载接口
func RangeDownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	username := r.Form.Get("username")

	fm, _ := meta.GetFileMetaDB(fsha1)
	userFile, err := dblayer.QueryUserFileMeta(username, fsha1)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 使用本地目录文件
	fpath := cfg.MergeLocalRootDir + fm.FileSha1
	fmt.Println("range-download-fpath: " + fpath)

	f, err := os.Open(fpath)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+userFile.FileName+"\"")
	http.ServeFile(w, r, fpath)
}
