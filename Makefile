# 设置变量
APP_NAME := crawler
VERSION ?= v1.0.0
DOCKER_IMAGE := $(APP_NAME):$(VERSION)

# 构建镜像
build:
	docker build -t $(DOCKER_IMAGE) -f deploy/Dockerfile .

# 运行容器
run:
	docker run -d \
		--name $(APP_NAME) \
		-p 8080:8080 \
		-v $(PWD)/config:/app/config \
		-v $(PWD)/data/logs:/app/data/logs \
		-v $(PWD)/data/cookies:/app/data/cookies \
		$(DOCKER_IMAGE)

# 停止容器
stop:
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

# 清理容器和镜像
clean: stop
	docker rmi $(DOCKER_IMAGE) || true

# 查看日志
logs:
	docker logs -f $(APP_NAME)

# 重启容器
restart: stop run

# 创建必要的目录结构
init:
	mkdir -p config data/logs data/cookies

# 声明伪目标
.PHONY: help build run stop clean logs restart init