# 2501YTC
2025-01-14 Winter Vacation Youth Training Camp

## 技术栈

## 业务逻辑

## 如何开发
1. 准备
必备清单
- Go
- Docker
- [cwgo](https://github.com/cloudwego/cwgo)
- kitex `go install github.com/cloudwego/kitex/tool/cmd/kitex@latest`
- ...

2. 克隆该项目并切换到`dev`分支，然后切出自己的分支进行开发

```shell
git clone https://github.com/1850298154/2501YTC.git
cd 2501YTC

git checkout dev
git checkout -b <姓名>/<日期>/<分支功能>
```
在开发过程中，请确保自己分支的进度与`dev`分支同步
```shell
git pull origin dev
```
3. 拷贝 `.env` 文件,设置配置文件并完善依赖
```shell
make init
make tidy
```

4. 启动环境容器(所需要的数据库等)
```shell
make env-start
```
如果你想停止他们的docker应用程序，可以运行`make env-stop`。

5. 依据idl生成代码，并在工作区同步
```shell
#创建文件夹
mkdir app/{app}
#生成客户端或服务端的代码
make gen-client svc={app}
make gen-server svc={app}
#格式化一下并完善依赖
make fmt
make tidy
# 同步工作区模块
go work sync
```
6. 每次提交 commit 前，先运行以下命令来格式化和检查静态代码（需要安装 [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) 和 [gofmt](https://pkg.go.dev/cmd/gofmt)和 [gofumpt](https://github.com/mvdan/gofumpt)）
```shell
make fmt
make lint
```


### Make 其他用法
```shell
make
```

