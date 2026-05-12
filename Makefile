# EZ Admin Gin — 统一开发入口
# 使用方法: make <target>
# Windows 用户: 安装 make (choco install make 或 scoop install make)，或直接运行对应命令。

.DEFAULT_GOAL := help

# ---------- 配置 ----------

SERVER_DIR  := server
ADMIN_DIR   := admin
DOCS_DIR    := docs
DEPLOY_DIR  := deploy

GO          := go
PNPM        := pnpm
DOCKER      := docker

# ---------- 帮助 ----------

.PHONY: help
help: ## 显示所有可用命令
	@echo "EZ Admin Gin — 可用命令："
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "提示: Windows 用户可安装 make (choco install make)，或直接查看 Makefile 中的对应命令。"

# ---------- 开发 ----------

.PHONY: dev
dev: ## 启动后端 + 前端 (需要先 make docker-config 启动数据库)
	@echo ">>> 启动后端..."
	$(MAKE) server-dev &
	@echo ">>> 启动前端..."
	$(MAKE) admin-dev

.PHONY: server-dev
server-dev: ## 启动后端 (go run)
	cd $(SERVER_DIR) && $(GO) run .

.PHONY: admin-dev
admin-dev: ## 启动前端 (pnpm dev)
	cd $(ADMIN_DIR) && $(PNPM) dev

.PHONY: docs-dev
docs-dev: ## 启动文档站
	cd $(DOCS_DIR) && $(PNPM) docs:dev

# ---------- 测试 ----------

.PHONY: test
test: server-test test-contract ## 运行所有测试 (后端 unit + 契约)
	@echo ">>> 后端测试完成"

.PHONY: server-test
server-test: ## 运行后端测试 (go test ./..., 不含集成测试)
	cd $(SERVER_DIR) && $(GO) test ./... -timeout 60s

.PHONY: test-api
test-api: ## 运行 API 黑盒集成测试 (需要 DB + Redis)
	cd $(SERVER_DIR) && $(GO) test -tags integration -v -timeout 120s ./tests/api/...

.PHONY: test-rbac
test-rbac: ## 运行 RBAC 权限流程测试 (需要 DB + Redis)
	cd $(SERVER_DIR) && $(GO) test -tags integration -v -timeout 120s ./tests/rbac/...

.PHONY: test-contract
test-contract: ## 运行 OpenAPI 契约测试 (不需要 DB/Redis)
	cd $(SERVER_DIR) && $(GO) test -v -timeout 60s ./tests/contract/...

.PHONY: test-integration
test-integration: ## 运行所有集成测试 (API + RBAC, 需要 DB + Redis)
	cd $(SERVER_DIR) && $(GO) test -p 1 -tags integration -v -timeout 180s ./tests/api/... ./tests/rbac/...

.PHONY: test-e2e
test-e2e: ## 运行 E2E 测试 (需要前端 + 后端运行中)
	cd $(ADMIN_DIR) && pnpm exec playwright test

# ---------- 代码检查 ----------

.PHONY: lint
lint: server-vet admin-check ## 运行所有 lint (后端 vet + 前端检查)
	@echo ">>> 所有 lint 完成"

.PHONY: server-vet
server-vet: ## 后端 go vet
	cd $(SERVER_DIR) && $(GO) vet ./...

.PHONY: server-mod
server-mod: ## 后端 go mod tidy (检查依赖一致性)
	cd $(SERVER_DIR) && $(GO) mod tidy

.PHONY: admin-check
admin-check: ## 前端类型检查 + lint
	cd $(ADMIN_DIR) && $(PNPM) type-check && $(PNPM) lint

# ---------- 构建 ----------

.PHONY: build
build: server-build admin-build ## 构建后端 + 前端
	@echo ">>> 构建完成"

.PHONY: server-build
server-build: ## 编译后端二进制
	cd $(SERVER_DIR) && $(GO) build -ldflags="-s -w" -o server .

.PHONY: admin-build
admin-build: ## 构建前端产物
	cd $(ADMIN_DIR) && $(PNPM) build

# ---------- Docker ----------

.PHONY: docker-config
docker-config: ## 验证 Docker Compose 配置 (本地开发)
	$(DOCKER) compose -f $(DEPLOY_DIR)/compose.local.yml config -q
	@echo ">>> compose.local.yml 配置有效"

.PHONY: docker-build
docker-build: ## 构建后端 + 前端 Docker 镜像
	$(DOCKER) build -t ez-admin-server $(SERVER_DIR)
	$(DOCKER) build -t ez-admin-admin $(ADMIN_DIR)

.PHONY: docker-up
docker-up: ## 启动本地 PostgreSQL + Redis
	$(DOCKER) compose -f $(DEPLOY_DIR)/compose.local.yml up -d

.PHONY: docker-down
docker-down: ## 停止本地 PostgreSQL + Redis
	$(DOCKER) compose -f $(DEPLOY_DIR)/compose.local.yml down

# ---------- 安装 ----------

.PHONY: install
install: ## 安装前端依赖
	cd $(ADMIN_DIR) && $(PNPM) install

# ---------- 清理 ----------

.PHONY: clean
clean: ## 清理构建产物
	rm -f $(SERVER_DIR)/server
	rm -rf $(ADMIN_DIR)/dist
