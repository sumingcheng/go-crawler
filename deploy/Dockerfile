FROM golang:1.22 AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
ENV HTTP_PROXY=http://host.docker.internal:7890
ENV HTTPS_PROXY=http://host.docker.internal:7890

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=linux go build -o crawler cmd/main.go

FROM mcr.microsoft.com/playwright:v1.42.0-jammy

ENV TZ=Asia/Shanghai
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /app

COPY --from=builder /app/crawler .

# 创建必要的目录
RUN mkdir -p /app/data/logs /app/data/cookies /app/config

# 配置清华源
RUN sed -i 's/archive.ubuntu.com/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list \
  && sed -i 's/security.ubuntu.com/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list

# 安装系统依赖
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  tzdata \
  curl \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

# 安装 Playwright 依赖
RUN playwright install chromium \
  && playwright install-deps chromium

# 配置时区
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
  && echo $TZ > /etc/timezone

ENV GIN_MODE=release
ENV CONFIG_PATH=/app/config/config

VOLUME ["/app/config", "/app/data/logs", "/app/data/cookies"]

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./crawler"]