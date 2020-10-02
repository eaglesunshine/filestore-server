package meta

import (
	mydb "filestore-server/db"
	"sort"
)

//文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

//init:内建init方法，在程序首次运行时将执行一次
func init() {
	fileMetas = make(map[string]FileMeta)
}

//UpdateFileMeta:新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// UpdateFileMetaDB : 新增/更新文件元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(
		fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

//GetFileMeta: 通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

//GetLastFileMetas: 获取批量的文件元信息列表
func GetLastFileMetas(count int) []FileMeta {
	var fMetaArray []FileMeta
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}

	sort.Sort(ByUploadTime(fMetaArray))
	if count > len(fMetaArray) {
		count = len(fMetaArray)
	}

	return fMetaArray[0:count]
}

//RemoveFileMeta: 是你出元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}

//GetFileMetaDb: 从mysql获取文件元信息
func GetFileMetaDB(filesha1 string) (FileMeta, error) {
	tfile, err := mydb.GetFileMeta(filesha1)
	if err != nil {
		return FileMeta{}, err
	}

	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil
}
