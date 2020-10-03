## 划重点

- ceph io 流程示意图(网摘)

<img src='/doc/ceph_io.jpg'></img>

## 关于 go modules

在 go1.11 之后，官方推出了 Go Modules 这原生的包管理工具；相对于 vendor, go mod 可以更有效的进行包版本管理。

- 在 1.12 及之前， go modules 需要配置环境变量来开启此功能:

```bash
# 分别有 on off auto
go env -w GO111MODULE=on
```

- 配置代理。因为众所周知的原因，有些包我们国内无法访问，一般需要通过代理(如 goproxy.cn):

```bash
go env -w GOSUMDB=sum.golang.google.cn
go env -w GOPROXY=https://goproxy.cn,direct
# 查看是否成功(go env的输出中包含代理信息)
go env
```

- go mod 初始化

```
go mod init <指定一个module名，如工程名>
```

在项目根目录下成功初始化后，会生成一个 go.mod 和 go.sum 文件。在之后执行 go build 或 go run 时，会触发依赖解析，并且下载对应的依赖包。

更具体的用法可以参考网上其他教程哦(如https://github.com/golang/go/wiki/Modules)。
