# ---- 构建阶段：编译前端资源 ----
FROM node:22-alpine AS builder

WORKDIR /app

# 先复制依赖文件，利用 Docker 层缓存。
COPY package.json pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile

# 再复制源码并构建。
COPY . .
RUN pnpm build

# ---- 运行阶段：Nginx 托管静态资源 ----
FROM nginx:1.27-alpine

# 删除默认配置，使用项目自定义配置。
RUN rm /etc/nginx/conf.d/default.conf

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
