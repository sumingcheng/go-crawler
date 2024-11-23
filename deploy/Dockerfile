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
ENV HTTP_PROXY=http://host.docker.internal:7890
ENV HTTPS_PROXY=http://host.docker.internal:7890
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright

WORKDIR /app

COPY --from=builder /app/crawler .

RUN mkdir -p /app/data/logs /app/data/cookies /app/config \
  && apt-get update \
  && apt-get install -y --no-install-recommends \
  tzdata \
  curl \
  ca-certificates \
  && playwright install chromium \
  && playwright install-deps chromium \
  && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
  && echo $TZ > /etc/timezone \
  && rm -rf /var/lib/apt/lists/*

ENV GIN_MODE=release
ENV CONFIG_PATH=/app/config/config.yaml

VOLUME ["/app/config", "/app/data/logs", "/app/data/cookies"]

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./crawler"]