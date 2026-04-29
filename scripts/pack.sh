#!/bin/bash
# 一键构建并打包部署文件
# 使用方法：在项目根目录执行 bash scripts/pack.sh
set -e

PACK_DIR="deploy-package"

echo ">>> 清理旧打包..." && rm -rf "$PACK_DIR" && mkdir -p "$PACK_DIR"

echo ">>> 编译后端 (Linux amd64)..."
(cd server && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "../$PACK_DIR/server" .)

echo ">>> 构建前端..."
(cd admin && pnpm install 2>/dev/null; pnpm build && cp -r dist "../$PACK_DIR/dist")

echo ">>> 打包配置文件..."
cp deploy/compose.server.yml deploy/.env.example deploy/ez-admin.service scripts/setup-server.sh scripts/update-server.sh "$PACK_DIR/"
mkdir -p "$PACK_DIR/nginx" "$PACK_DIR/ssl" "$PACK_DIR/configs"
cp server/configs/config.yaml server/configs/rbac_model.conf "$PACK_DIR/configs/"
cp deploy/nginx/nginx-native.conf "$PACK_DIR/nginx/"

echo ">>> 生成压缩包..."
tar czf deploy-package.tar.gz -C "$PACK_DIR" .

echo ""
echo "✅ 打包完成！上传 deploy-package.tar.gz 到服务器即可。"
echo "   上传目标：/opt/ez-admin/"
echo "   然后执行："
echo "     mkdir -p /opt/ez-admin && cd /opt/ez-admin"
echo "     tar xzf ~/deploy-package.tar.gz"
echo "     bash setup-server.sh"
