## 开发命令

### get

```shell
go env -w GOPROXY=https://goproxy.cn,direct
# go env -w GOPROXY=https://proxy.golang.org,direct
# go env -w GOPROXY=https://goproxy.io,direct
# go env -w GOPROXY=https://mirrors.aliyun.com/goproxy,direct
# go env -w GOPROXY=https://mirrors.cloud.tencent.com/go,direct
go get -u github.com/xuxiaowei-com-cn/git-go@main
```

### mod

```shell
go mod tidy
```

```shell
go mod download
```

### run

```shell
go run main.go
```

### run help

```shell
go run main.go help
```
