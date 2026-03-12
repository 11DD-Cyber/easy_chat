# syntax=docker/dockerfile:1.6

FROM golang:1.22 AS builder
WORKDIR /src

# 先拉依赖，加速后续构建
COPY go.mod go.sum ./
RUN go mod download

# 拷贝全部源码并编译 devrun
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/devrun ./cmd/devrun

FROM golang:1.22
WORKDIR /app

# 带上源码以便 devrun 内部调用 go run
COPY --from=builder /src /app
COPY --from=builder /out/devrun /usr/local/bin/devrun

ENV GO111MODULE=on

# 暴露常见端口（根据需要增减）
EXPOSE 8888 9090 10001 8080 10090 10091

ENTRYPOINT ["devrun"]
