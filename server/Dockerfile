# ---- 构建阶段：编译 Go 二进制 ----
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

# 先复制依赖文件，利用 Docker 层缓存加速后续构建。
COPY go.mod go.sum ./
RUN go mod download

# 再复制源码并编译。
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

# ---- 运行阶段：只保留二进制和配置 ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/server .
COPY configs/ ./configs/

# 上传目录和日志目录需要在运行时存在。
RUN mkdir -p uploads logs

EXPOSE 8080

CMD ["./server"]
