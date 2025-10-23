# 阶段1: 打包构建阶段
FROM golang:1.24.0-alpine AS builder

# 设置工作目录
WORKDIR /build

# 设置Go环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 复制go.mod和go.sum文件，利用Docker缓存优化依赖下载
COPY go.mod go.sum ./
RUN go mod download

# 复制项目源代码
COPY . .

# 编译项目，生成静态链接的二进制文件
RUN go build -ldflags="-s -w" -o xarr-pay-merchant-xbot main.go

# 阶段2: 运行环境阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone


# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /build/xarr-pay-merchant-xbot .

# 创建必要的目录
RUN mkdir -p logs data/storage

# 启动应用
ENTRYPOINT ["./xarr-pay-merchant-xbot"]
