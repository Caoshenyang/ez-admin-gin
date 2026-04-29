#!/bin/bash
# EZ Admin 首次部署脚本（服务器端）
# 负责初始化 Docker 环境、生成密钥、启动服务、创建管理员。
# 后续更新请使用 update-server.sh。
#
# 使用方法：sudo bash /opt/ez-admin/setup-server.sh

set -e

BASE="/opt/ez-admin"

# ---- 1. 整理文件 ----

echo ">>> 创建目录..."
mkdir -p "$BASE/nginx" "$BASE/ssl" "$BASE/data/postgres" "$BASE/data/redis"

echo ">>> 整理文件..."
if [ -d "$BASE/dist" ]; then
  rm -rf "$BASE/web"
  mv "$BASE/dist" "$BASE/web"
  echo "    前端文件已就位"
fi

[ -f "$BASE/nginx-native.conf" ] && mv "$BASE/nginx-native.conf" "$BASE/nginx/nginx-native.conf"

if [ -f "$BASE/.env.example" ] && [ ! -f "$BASE/.env" ]; then
  mv "$BASE/.env.example" "$BASE/.env"
  echo "    .env 已创建"
fi

if [ -f "$BASE/ez-admin.service" ]; then
  sudo cp "$BASE/ez-admin.service" /etc/systemd/system/
  sudo systemctl daemon-reload
fi

[ -f "$BASE/server" ] && chmod +x "$BASE/server"

# ---- 2. 生成 JWT 密钥（仅首次）----

if grep -q "change-me-to-a-random-string" "$BASE/.env" 2>/dev/null; then
  SECRET=$(openssl rand -hex 32)
  sed -i "s/change-me-to-a-random-string-at-least-32-chars/$SECRET/" "$BASE/.env"
  echo "    JWT 密钥已自动生成"
fi

# ---- 3. 启动基础服务 ----

echo ">>> 启动 PostgreSQL + Redis + Nginx..."
cd "$BASE"
docker compose -f compose.server.yml up -d

echo "    等待数据库就绪..."
for i in $(seq 1 30); do
  if docker compose -f compose.server.yml exec -T postgres pg_isready -U ez_admin >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

# ---- 4. 启动后端 ----

echo ">>> 启动后端..."
sudo systemctl enable --now ez-admin

echo "    等待后端就绪..."
for i in $(seq 1 15); do
  if curl -sf http://localhost/health >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

# ---- 5. 初始化管理员（仅首次）----

STATUS=$(curl -sf -o /dev/null -w "%{http_code}" http://localhost/api/v1/setup/init \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"Admin@123456","nickname":"管理员"}' 2>/dev/null || true)

if [ "$STATUS" = "200" ]; then
  echo "    管理员已创建（admin / Admin@123456，请登录后修改密码）"
elif [ "$STATUS" = "409" ]; then
  echo "    管理员已存在，跳过初始化"
else
  echo "    ⚠️ 管理员初始化返回 $STATUS，请手动执行："
  echo "    curl -X POST http://localhost/api/v1/setup/init -H 'Content-Type: application/json' -d '{\"username\":\"admin\",\"password\":\"Admin@123456\",\"nickname\":\"管理员\"}'"
fi

# ---- 6. 结果 ----

echo ""
echo "========================================="
echo "✅ 部署完成！"
echo ""
echo "  访问地址：http://$(hostname -I | awk '{print $1}')"
echo "  默认账号：admin / Admin@123456"
echo ""
echo "  查看后端日志：sudo journalctl -u ez-admin -f"
echo "  查看容器状态：docker compose -f $BASE/compose.server.yml ps"
echo "========================================="
