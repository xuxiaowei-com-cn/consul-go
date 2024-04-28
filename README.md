## 开发命令

### get

```shell
go env -w GOPROXY=https://goproxy.cn,direct
# go env -w GOPROXY=https://proxy.golang.org,direct
# go env -w GOPROXY=https://goproxy.io,direct
# go env -w GOPROXY=https://mirrors.aliyun.com/goproxy,direct
# go env -w GOPROXY=https://mirrors.cloud.tencent.com/go,direct
go get -u github.com/google/go-querystring
go get -u github.com/hashicorp/go-cleanhttp
go get -u github.com/hashicorp/go-retryablehttp
go get -u github.com/stretchr/testify
go get -u github.com/urfave/cli/v2
go get -u github.com/xuxiaowei-com-cn/git-go@main
go get -u gopkg.in/yaml.v3
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

- Windows 环境为 %xxx%
- Linux 环境为 $xxx

```shell
go run main.go --help
```

```shell
$ go run main.go --help
NAME:
   consul-go - 基于 Go 语言开发的 consul 命令行工具

USAGE:
   consul-go [global options] command [command options]

VERSION:
   dev

AUTHOR:
   徐晓伟 <xuxiaowei@xuxiaowei.com.cn>

COMMANDS:
   kv       Key / Value
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

COPYRIGHT:
   徐晓伟工作室 <xuxiaowei@xuxiaowei.com.cn>
```

- Key / Value

```shell
$ go run main.go kv --help
NAME:
   consul-go kv - Key / Value

USAGE:
   consul-go kv command [command options]

COMMANDS:
   export   导出
   import   导入
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --base-url value  Consul URL (default: "http://127.0.0.1:8500") [%CONSUL_GO_BASE_URL%]
   --dc value        dc (default: "dc1") [%CONSUL_GO_DC%]
   --help, -h        show help
```

- Key / Value：export

```shell
$ go run main.go kv export --help
NAME:
   consul-go kv export - 导出

USAGE:
   consul-go kv export [command options] [arguments...]

OPTIONS:
   --base-url value       Consul URL (default: "http://127.0.0.1:8500") [%CONSUL_GO_BASE_URL%]
   --dc value             dc (default: "dc1") [%CONSUL_GO_DC%]
   --export-folder value  导出文件夹 (default: "consul-go-export-folder") [%CONSUL_GO_DC_EXPORT_FOLDER%]
   --help, -h             show help
```

- Key / Value：import

```shell
$ go run main.go kv import --help
NAME:
   consul-go kv import - 导入

USAGE:
   consul-go kv import [command options] [arguments...]

OPTIONS:
   --base-url value       Consul URL (default: "http://127.0.0.1:8500") [%CONSUL_GO_BASE_URL%]
   --dc value             dc (default: "dc1") [%CONSUL_GO_DC%]
   --import-folder value  导入文件夹 (default: "consul-go-import-folder") [%CONSUL_GO_DC_IMPORT_FOLDER%]
   --help, -h             show help
```

### test

```shell
go test ./... -v
```

### build

```shell
go build
# GOOS=：设置构建的目标操作系统（darwin | freebsd | linux | windows）
# GOARCH=：设置构建的目标操作系统（386 | amd64 | arm | arm64）
# -v：打印编译过程中的详细信息
# -ldflags：设置在编译时传递给链接器的参数
# -ldflags "-s -w -buildid="
#                           -s: 删除符号表信息，减小可执行文件的大小。
#                           -w: 删除调试信息，使可执行文件在运行时不会打印调试信息。
#                           -buildid=: 删除构建ID，使可执行文件在运行时不会打印构建ID。
# -trimpath：去掉所有包含 go path 的路径
# -o：指定构建后输出的文件名
```

- Windows
    - amd64
        ```shell
        go build -o buildinfo/buildinfo.exe buildinfo/buildinfo.go
        GOOS=windows GOARCH=amd64 go build -v -ldflags "-s -w -buildid= -X main.BuildDate=$(buildinfo/buildinfo.exe now) -X main.Compiler= -X main.GitCommitBranch=$(buildinfo/buildinfo.exe commitBranch) -X main.GitCommitSha=$(buildinfo/buildinfo.exe commitSha) -X main.GitCommitShortSha=$(buildinfo/buildinfo.exe commitShortSha) -X main.GitCommitTag=$(buildinfo/buildinfo.exe commitTag) -X main.GitCommitTimestamp=$(buildinfo/buildinfo.exe commitTimestamp) -X main.GitTreeState=$(buildinfo/buildinfo.exe git-tree-state) -X main.GitVersion=$(buildinfo/buildinfo.exe commitTag) -X main.GoVersion=$(buildinfo/buildinfo.exe goShortVersion) -X main.Major= -X main.Minor= -X main.Revision= -X main.Platform=windows/amd64 -X main.CiPipelineId= -X main.CiJobId=" -trimpath -o consul-go-windows-amd64.exe .
        ```
    - arm64
        ```shell
        go build -o buildinfo/buildinfo.exe buildinfo/buildinfo.go
        GOOS=windows GOARCH=arm64 go build -v -ldflags "-s -w -buildid= -X main.BuildDate=$(buildinfo/buildinfo.exe now) -X main.Compiler= -X main.GitCommitBranch=$(buildinfo/buildinfo.exe commitBranch) -X main.GitCommitSha=$(buildinfo/buildinfo.exe commitSha) -X main.GitCommitShortSha=$(buildinfo/buildinfo.exe commitShortSha) -X main.GitCommitTag=$(buildinfo/buildinfo.exe commitTag) -X main.GitCommitTimestamp=$(buildinfo/buildinfo.exe commitTimestamp) -X main.GitTreeState=$(buildinfo/buildinfo.exe git-tree-state) -X main.GitVersion=$(buildinfo/buildinfo.exe commitTag) -X main.GoVersion=$(buildinfo/buildinfo.exe goShortVersion) -X main.Major= -X main.Minor= -X main.Revision= -X main.Platform=windows/arm64 -X main.CiPipelineId= -X main.CiJobId=" -trimpath -o consul-go-windows-arm64.exe .
        ```

- Linux
    - amd64
        ```shell
        go build -o buildinfo/buildinfo buildinfo/buildinfo.go
        GOOS=linux GOARCH=amd64 go build -v -ldflags "-s -w -buildid= -X main.BuildDate=$(buildinfo/buildinfo now) -X main.Compiler= -X main.GitCommitBranch=$(buildinfo/buildinfo commitBranch) -X main.GitCommitSha=$(buildinfo/buildinfo commitSha) -X main.GitCommitShortSha=$(buildinfo/buildinfo commitShortSha) -X main.GitCommitTag=$(buildinfo/buildinfo commitTag) -X main.GitCommitTimestamp=$(buildinfo/buildinfo commitTimestamp) -X main.GitTreeState=$(buildinfo/buildinfo git-tree-state) -X main.GitVersion=$(buildinfo/buildinfo commitTag) -X main.GoVersion=$(buildinfo/buildinfo goShortVersion) -X main.Major= -X main.Minor= -X main.Revision= -X main.Platform=linux/amd64 -X main.CiPipelineId= -X main.CiJobId=" -trimpath -o consul-go-linux-amd64 .
        ```
    - arm64
        ```shell
        go build -o buildinfo/buildinfo buildinfo/buildinfo.go
        GOOS=linux GOARCH=arm64 go build -v -ldflags "-s -w -buildid= -X main.BuildDate=$(buildinfo/buildinfo now) -X main.Compiler= -X main.GitCommitBranch=$(buildinfo/buildinfo commitBranch) -X main.GitCommitSha=$(buildinfo/buildinfo commitSha) -X main.GitCommitShortSha=$(buildinfo/buildinfo commitShortSha) -X main.GitCommitTag=$(buildinfo/buildinfo commitTag) -X main.GitCommitTimestamp=$(buildinfo/buildinfo commitTimestamp) -X main.GitTreeState=$(buildinfo/buildinfo git-tree-state) -X main.GitVersion=$(buildinfo/buildinfo commitTag) -X main.GoVersion=$(buildinfo/buildinfo goShortVersion) -X main.Major= -X main.Minor= -X main.Revision= -X main.Platform=linux/arm64 -X main.CiPipelineId= -X main.CiJobId=" -trimpath -o consul-go-linux-arm64 .
        ```

- LoongArch
    - 64-bit
        ```shell
        go build -o buildinfo/buildinfo buildinfo/buildinfo.go
        GOOS=linux GOARCH=loong64 go build -v -ldflags "-s -w -buildid= -X main.BuildDate=$(buildinfo/buildinfo now) -X main.Compiler= -X main.GitCommitBranch=$(buildinfo/buildinfo commitBranch) -X main.GitCommitSha=$(buildinfo/buildinfo commitSha) -X main.GitCommitShortSha=$(buildinfo/buildinfo commitShortSha) -X main.GitCommitTag=$(buildinfo/buildinfo commitTag) -X main.GitCommitTimestamp=$(buildinfo/buildinfo commitTimestamp) -X main.GitTreeState=$(buildinfo/buildinfo git-tree-state) -X main.GitVersion=$(buildinfo/buildinfo commitTag) -X main.GoVersion=$(buildinfo/buildinfo goShortVersion) -X main.Major= -X main.Minor= -X main.Revision= -X main.Platform=darwin/amd64 -X main.CiPipelineId= -X main.CiJobId=" -trimpath -o consul-go-loong64 .
        ```

- Darwin
    - amd64
        ```shell
        go build -o buildinfo/buildinfo buildinfo/buildinfo.go
        GOOS=darwin GOARCH=amd64 go build -v -ldflags "-s -w -buildid= -X main.BuildDate=$(buildinfo/buildinfo now) -X main.Compiler= -X main.GitCommitBranch=$(buildinfo/buildinfo commitBranch) -X main.GitCommitSha=$(buildinfo/buildinfo commitSha) -X main.GitCommitShortSha=$(buildinfo/buildinfo commitShortSha) -X main.GitCommitTag=$(buildinfo/buildinfo commitTag) -X main.GitCommitTimestamp=$(buildinfo/buildinfo commitTimestamp) -X main.GitTreeState=$(buildinfo/buildinfo git-tree-state) -X main.GitVersion=$(buildinfo/buildinfo commitTag) -X main.GoVersion=$(buildinfo/buildinfo goShortVersion) -X main.Major= -X main.Minor= -X main.Revision= -X main.Platform=darwin/amd64 -X main.CiPipelineId= -X main.CiJobId=" -trimpath -o consul-go-darwin-amd64 .
        ```
    - arm64
        ```shell
        go build -o buildinfo/buildinfo buildinfo/buildinfo.go
        GOOS=darwin GOARCH=arm64 go build -v -ldflags "-s -w -buildid= -X main.BuildDate=$(buildinfo/buildinfo now) -X main.Compiler= -X main.GitCommitBranch=$(buildinfo/buildinfo commitBranch) -X main.GitCommitSha=$(buildinfo/buildinfo commitSha) -X main.GitCommitShortSha=$(buildinfo/buildinfo commitShortSha) -X main.GitCommitTag=$(buildinfo/buildinfo commitTag) -X main.GitCommitTimestamp=$(buildinfo/buildinfo commitTimestamp) -X main.GitTreeState=$(buildinfo/buildinfo git-tree-state) -X main.GitVersion=$(buildinfo/buildinfo commitTag) -X main.GoVersion=$(buildinfo/buildinfo goShortVersion) -X main.Major= -X main.Minor= -X main.Revision= -X main.Platform=darwin/arm64 -X main.CiPipelineId= -X main.CiJobId=" -trimpath -o consul-go-darwin-arm64 .
        ```
