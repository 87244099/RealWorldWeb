# 安装go 课程里面的go-apline已经失效，改成下面这个看看
FROM golang:1.22.3-alpine
# 设置工作目录
WORKDIR /app
# 把宿主机上的当前目录，拷贝到WORKDIR下
COPY . .
# 设置go的代理镜像，方便下载代码
ENV GOPROXY=https://goproxy.cn,direct
# 对项目进行go的打包，然后会生成一个和和宿主环境同名的可执行文件，比如RealWorldWeb.exe
RUN go build .
# 镜像运行后，对外暴露8080端口：Q：go服务内部绑定了8000端口，但是这里对外变成8080，不知道行不行
EXPOSE 8080
# 执行上面打包出来 RealWorldWeb.exe，这里不带exe是因为你省略的他也是可以执行的
CMD ./RealWorldWeb