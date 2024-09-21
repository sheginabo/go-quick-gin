# 構建階段
FROM --platform=linux/amd64 golang:1.23 AS builder
# 設置工作目錄
WORKDIR /app
# Install the necessary packages for building and running the Golang application
RUN apk add --no-cache \
  git \
  go \
  musl-dev
# 複製 go mod 和 sum 文件 優化緩存
COPY go.mod go.sum ./
# 下載依賴
RUN go mod download
# 複製源代碼
COPY . ./
# 構建應用
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# 最終階段
FROM --platform=linux/amd64 alpine:latest
WORKDIR /app
# 從構建階段複製二進制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080
# 運行
CMD ["/app/main"]