## 关于应用启动

- 启动上传应用程序:
```bash
# cd $GOPATH/src/<你的工程目录>
> cd $GOPATH/src/filestore-server
> go run service/upload/main.go
```

- 启动转移应用程序:
```bash
# cd $GOPATH/src/<你的工程目录>
> cd $GOPATH/src/filestore-server
> go run service/transfer/main.go
```

## 关于go版本更新(升级至1.13或以上)

go升级到1.13及以上版本时，默认开启go modules包管理工具，推荐此方式管理依赖库。

- 配置代理。因为众所周知的原因，有些包我们国内无法访问，一般需要通过代理(如goproxy.cn):

```bash
go env -w GOSUMDB=sum.golang.google.cn
go env -w GOPROXY=https://goproxy.cn,direct 
# 查看是否成功(go env的输出中包含代理信息)
go env
```

- go mod初始化

```
go mod init <指定一个modeul名，如工程名>
# 如 go mod init filestore-server
```

在项目根目录下成功初始化后，会生成一个go.mod和go.sum文件。在之后执行go build或go run时，会触发依赖解析，并且下载对应的依赖包。

更具体的用法可以参考网上其他教程哦(如https://github.com/golang/go/wiki/Modules)。