# picture-oss-proxy

项目简介：一个基于gin框架，对接口增加redis限流，转发请求的项目(类似网关)。

### 一、gin相关

中文学习文档

https://gin-gonic.com/zh-cn/docs/

go mod 管理

```
# 初始化go mod
go mod init
# 拉取依赖 依赖包会自动下载到$GOPATH/pkg/mod，多个项目可以共享缓存的mod
go mod download
# 整理依赖关系
go mod tidy
# 缓存到vendor目录 从mod中拷贝到项目的vendor目录下，这样IDE就可以识别了！
go mod vendor
```

### 二、项目相关

###### 1、启动项目：

go run main.go 



### 三、随便记录

删除.DS_Store文件
 find . -name .DS_Store -print0 | xargs -0 git rm -f --ignore-unmatch

### 四、docker相关

镜像启动命令

```
docker run -d -p 3000:3000  --name gintool --network host gintool
```

镜像编译命令

```
docker build  -t gintool . 
```



在 Dockerfile.multistage 中使用了两次 FROM 指令，分别对应两个构建阶段。第一个阶段构建的 FROM 指令依然使用 golang:1.16-alpine 作为基础镜像，并将该阶段命名为 build。第二个构建阶段的 FROM 指令使用 scratch 作为基础镜像，告诉 Docker 接下来从一个全新的基础镜像开始构建，scratch 镜像是 Docker 项目预定义的最小的镜像。第二阶段构建主要是将上个阶段中编译好的二进制文件复制到新的镜像中。
在 Go 应用中，多阶段构建非常常见，可以减小镜像的体积、节省大量的存储空间。
在 Dockerfile.multistage 中需要额外关注的是 RUN 指令，这里使用到了交叉编译。
交叉编译
交叉编译是指在一个平台上生成另一个平台的可执行程序。

在其他编程语言中进行交叉编译可能要借助第三方工具，但 Go 内置了交叉编译工具，使用起来非常方便，通常设置 CGO_ENABLED、GOOS 和 GOARCH 这几个环境变量就够了。

##### CGO_ENABLED

默认值是 1，即默认开启 cgo，允许在 Go 代码中调用 C 代码。

当 CGO_ENABLED=1 进行编译时，会将文件中引用 libc 的库（比如常用的 net 包）以动态链接的方式生成目标文件；
当 CGO_ENABLED=0 进行编译时，则会把在目标文件中未定义的符号（如外部函数）一起链接到可执行文件中。

所以交叉编译时，我们需要将 CGO_ENABLED 设置为 0。

##### GOOS 和 GOARCH

GOOS 是目标平台的操作系统，如 linux、windows，注意 macOS 的值是 darwin。默认是当前操作系统。
GOARCH 是目标平台的 CPU 架构，如 amd64、arm、386 等。默认值是当前平台的 CPU 架构。
Go 支持的所有操作系统和 CPU 架构可以查看 syslist.go 。
我们可以使用 go env 命令获取当前 GOOS 和 GOARCH 的值。例如我当前的操作系统是 macOS：

```
$ go env GOOS GOARCH
darwin
amd64
```

所以在本文的多阶段构建 Dockerfile.multistage 中，构建命令是：

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /picture-oss-proxy
```

由于我们现在有两个 Dockerfile，所以我们必须告诉 Docker 我们要使用新的 Dockerfile 进行构建。

```
docker build -t gintool:multistage -f Dockerfile.multistage . 
```

构建完成后，你会发现 gintool:multistage只有不到 8MB，比 gintool:latest 小了几十倍。

```
$ docker image ls
REPOSITORY                   TAG            IMAGE ID       CREATED         SIZE
gintool                                      multistage                01f235c37309   13 seconds ago   18MB
gintool                                      latest                    4243757a181b   2 days ago       652MB
```

运行 Go 镜像
现在我们有了 Go 应用的镜像，接下来就可以运行 Go 镜像查看应用程序是否正常运行。
要在 Docker 容器中运行镜像，我们可以使用 docker run 命令，参数是镜像名称：

```
 docker run -p 3000:3000 gintool:multistage
```



### 项目笔记

   把github.com/dgrijalva/jwt-go 替换为github.com/golang-jwt/jwt/v4
jwt-go是个人开发者的一个Go语言的JWT实现。jwt-go 4.0.0-preview1之前版本存在安全漏洞。攻击者可利用该漏洞在使用[]string{} for m[\"aud\"](规范允许)的情况下绕过预期的访问限制



### 测试命令
ab -c 2 -n 1000 http://localhost:3000/api/v1/ping



### TODO:

1⃣️增加黑白ip名单

2⃣️增加转发请求

