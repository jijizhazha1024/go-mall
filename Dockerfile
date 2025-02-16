# 第一阶段：构建
FROM golang:1.21 as builder

WORKDIR /project/go-mall

# 安装基础工具


# 设置Go环境
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

# 复制源代码
COPY . .

# 下载依赖
RUN go mod download

# 执行构建脚本
RUN bash ./scripts/build.sh
# 第二阶段：运行
FROM alpine:3.19

# 使用阿里云镜像
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk add --no-cache tzdata ca-certificates

# 设置时区
ENV TZ=Asia/Shanghai

# 创建必要的目录
WORKDIR /app
RUN mkdir -p apis/logs services/logs

# 从构建阶段复制编译好的文件
COPY --from=builder /app /app/ 