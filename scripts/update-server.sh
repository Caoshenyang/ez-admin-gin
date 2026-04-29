#!/bin/bash
# EZ Admin 更新脚本（服务器端）
# 仅替换文件并重启后端，不触碰 Docker 环境和已有配置。
#
# 使用方法：sudo bash /opt/ez-admin/update-server.sh

set -e

BASE="/opt/ez-admin"

# ---- 1. 整理文件 ----

echo ">>> 整理前端文件..."
if [ -d "$BASE/dist" ]; then
  rm -rf "$BASE/web"
  mv "$BASE/dist" "$BASE/web"
  echo "    前端文件已更新"
fi

echo ">>> 更新后端二进制..."
[ -f "$BASE/server" ] && chmod +x "$BASE/server"

# ---- 2. 重启后端 ----

echo ">>> 重启后端..."
sudo systemctl restart ez-admin

echo "    等待后端就绪..."
for i in $(seq 1 15); do
  if curl -sf http://localhost/health >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

# ---- 3. 结果 ----

echo ""
echo "========================================="
echo "✅ 更新完成！"
echo ""
echo "  查看后端日志：sudo journalctl -u ez-admin -f"
echo "========================================="
