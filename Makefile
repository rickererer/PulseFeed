# PulseFeed 开发命令统一入口
# 用法：make <target>，例如 make test、make build

SHELL := /bin/bash

.PHONY: dev test build lint clean help

## 启动本地开发环境（API + Web）
dev:
	./scripts/start.sh

## 运行全部测试（后端 go test + 前端构建校验）
test:
	cd apps/api && go test ./...
	cd apps/web && npm run build

## 生产构建（后端编译 + 前端打包）
build:
	cd apps/api && go build ./...
	cd apps/web && npm run build

## 静态检查（后端 vet + 前端类型/构建）
lint:
	cd apps/api && go vet ./...
	cd apps/web && npm run build

## 清理构建产物
clean:
	rm -rf apps/web/dist

## 帮助
help:
	@echo "PulseFeed 命令："
	@echo "  make dev   启动本地开发环境"
	@echo "  make test  运行全部测试"
	@echo "  make build 生产构建"
	@echo "  make lint  静态检查"
	@echo "  make clean 清理产物"
