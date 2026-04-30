#!/bin/bash
# EZ Admin 一键部署脚本 (macOS / Linux)
# 使用方法：在项目根目录执行
#   bash scripts/deploy.sh ubuntu@1.2.3.4

set -e

SERVER="${1:?用法: bash scripts/deploy.sh ubuntu@SERVER_IP}"
PACK_DIR="deploy-package"
TMP_TAR="deploy-package.tar.gz"

# ---- 1. 构建 ----

echo ">>> 编译后端 (Linux amd64)..."
(cd server && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "../$PACK_DIR/server" .)

echo ">>> 构建前端..."
(cd admin && pnpm install 2>/dev/null; pnpm build && rm -rf "../$PACK_DIR/dist" && cp -r dist "../$PACK_DIR/dist")

# ---- 2. 打包配置 ----

echo ">>> 打包配置文件..."
cp deploy/compose.server.yml deploy/.env.example deploy/ez-admin.service scripts/setup-server.sh "$PACK_DIR/"
mkdir -p "$PACK_DIR/nginx"
cp deploy/nginx/nginx-native.conf "$PACK_DIR/nginx/"

# ---- 3. 打包 + 上传 ----

echo ">>> 打包..."
rm -f "$TMP_TAR"
tar czf "$TMP_TAR" -C "$PACK_DIR" .

echo ">>> 上传到 $SERVER ..."
scp "$TMP_TAR" "$SERVER:/tmp/"
rm -f "$TMP_TAR"

# ---- 4. 远端初始化 ----

echo ">>> 远端初始化..."
ssh "$SERVER" "mkdir -p /opt/ez-admin && cd /opt/ez-admin && tar xzf /tmp/$TMP_TAR && rm /tmp/$TMP_TAR && bash setup-server.sh"
