# Stage 1: 创建项目二进制文件
FROM golang:1.26-alpine AS builder

# 设置工作目录
WORKDIR /lab4

# 复制 go.mod 和 go.sum 文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制整个项目源码到容器
COPY . .

# 构建 Go 二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myvideo_lab4

# Stage 2: 创建最终的镜像
FROM alpine:3.20

# 安装运行 Go 应用所需的依赖包（ca-certificates 用于支持 HTTPS 请求）
RUN apk update && apk add --no-cache ca-certificates netcat-openbsd  

# 设置工作目录
WORKDIR /lab4

# 从构建阶段复制 Go 二进制文件到最终镜像
COPY --from=builder /lab4/myvideo_lab4 /lab4/

# 复制静态文件
COPY ./static /lab4/static/

# 复制 .env 文件
COPY .env /lab4/.env

# 复制等待脚本
COPY wait-for-it.sh /lab4/wait-for-it.sh
RUN chmod +x /lab4/wait-for-it.sh  # 给脚本添加可执行权限

# 暴露容器的端口
EXPOSE 8080

# 启动应用程序
# CMD ["/bin/sh", "/lab4/wait-for-it.sh", "mysql", "3306", "--", "/lab4/wait-for-it.sh", "redis", "6379", "--", "/lab4/myvideo_lab4"]
