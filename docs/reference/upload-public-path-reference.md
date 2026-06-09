---
title: 上传与公开路径参考
description: "集中说明 upload.dir、upload.public_path、sys_file.path、sys_file.url 和 Nginx /uploads 代理之间的关系。"
---

# 上传与公开路径参考

这页专门用来查一个很容易混淆的问题：

> 文件上传后，磁盘路径、数据库路径、公开 URL 和 Nginx 代理路径，到底分别是什么关系？

::: tip 快速结论
当前主线里有两条路径一定要分开记：

- `path`：后端磁盘相对路径，例如 `uploads/20260501/xxx.pdf`
- `url`：给前端访问的公开路径，例如 `/uploads/20260501/xxx.pdf`
:::

## 先看真实配置

当前上传配置来自：

- `server/configs/config.yaml`
- `server/internal/config/config.go`

默认值是：

| 配置项 | 当前默认值 | 作用 |
| --- | --- | --- |
| `upload.dir` | `uploads` | 文件落盘目录 |
| `upload.public_path` | `/uploads` | 前端公开访问前缀 |
| `upload.max_size_mb` | `10` | 单文件大小限制 |
| `upload.allowed_exts` | 图片、文档白名单 | 文件类型限制 |

## 上传服务真正做了什么

当前真实上传逻辑在：

- `server/internal/module/system/file/service.go`

服务里会做下面几件事：

1. 校验大小和后缀
2. 按日期创建子目录
3. 生成后端文件名
4. 落盘到 `upload.dir/date/filename`
5. 同时写入数据库里的 `path` 和 `url`

其中最关键的两行是：

```go
relativePath := filepath.ToSlash(filepath.Join(s.cfg.Dir, dateDir, fileName))
url := publicPath + "/" + dateDir + "/" + fileName
```

也就是说：

- `relativePath` 会变成 `uploads/20260501/a.pdf`
- `url` 会变成 `/uploads/20260501/a.pdf`

## `sys_file` 里到底存什么

文件记录最容易混的是这两个字段：

| 字段 | 例子 | 含义 |
| --- | --- | --- |
| `path` | `uploads/20260501/20260501123000_abcd1234.pdf` | 服务端相对磁盘路径 |
| `url` | `/uploads/20260501/20260501123000_abcd1234.pdf` | 前端公开访问地址 |

可以直接把它们理解成：

- `path` 给服务端管理和排查用
- `url` 给前端展示和复制链接用

## 为什么 `upload.public_path` 一定要带前导 `/`

当前服务里会做一次规范化：

```go
if !strings.HasPrefix(publicPath, "/") {
  publicPath = "/" + publicPath
}
return strings.TrimRight(publicPath, "/")
```

所以当前推荐始终写成：

```yaml
upload:
  public_path: /uploads
```

而不是：

```yaml
upload:
  public_path: uploads/
```

前者更符合整个前端和 Nginx 的访问习惯。

## Nginx 为什么还要单独代理 `/uploads/`

当前服务器原生 Nginx 配置里，`/uploads/` 会继续反代到后端：

```nginx
location /uploads/ {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

这说明当前架构把上传文件访问仍然视为后端资源入口的一部分，而不是独立静态文件站。

这样做的好处是：

- 前端永远只记 `/uploads/...`
- Nginx、后端、本地开发都能沿同一条公开路径工作
- 部署时不需要前端额外知道宿主机真实磁盘目录

## 当前一条完整上传链怎么走

```text
前端上传文件
  ↓
POST /api/v1/system/files
  ↓
file/service.go 落盘到 uploads/YYYYMMDD/filename
  ↓
sys_file.path = uploads/YYYYMMDD/filename
sys_file.url  = /uploads/YYYYMMDD/filename
  ↓
前端拿 url 访问
  ↓
Nginx /uploads/ 反代到后端
```

## 本地开发时要对齐什么

本地开发最值得确认 3 件事：

1. `server/` 下 `uploads/` 目录可写
2. `upload.public_path` 仍然保持 `/uploads`
3. 前端展示文件时使用的是接口返回的 `url`，不是自己拼磁盘路径

如果这三件事都对齐，本地和线上行为会非常接近。

## 最常见的 4 个误区

### 1. 把 `path` 直接给前端用

错误示例：

```text
uploads/20260501/test.pdf
```

这是后端文件路径，不是浏览器公开 URL。

前端应优先使用：

```text
/uploads/20260501/test.pdf
```

### 2. 改了 `upload.dir`，但没同步理解 `url` 不会跟着变目录名

当前 `url` 取决于 `upload.public_path`，不是直接取决于磁盘目录名字。

所以：

- 磁盘目录可以是 `uploads`
- 公开路径也可以统一保持 `/uploads`

它们有关联，但不是同一个字段。

### 3. Nginx 只配了 `/api/`，没配 `/uploads/`

这会导致：

- 文件上传成功
- 数据库也有 `url`
- 但浏览器访问文件仍然 404

### 4. 页面复制链接时自己手拼域名和目录

当前文件页直接复制的就是接口返回的 `url`。如果后续要拼成绝对地址，也应该基于这个 `url` 再做外层站点前缀处理，而不是重造目录规则。

## 新项目复用时最稳的默认值

如果你暂时不想在文件访问路径上做额外设计，当前最稳的默认组合就是：

```yaml
upload:
  dir: uploads
  public_path: /uploads
```

配套 Nginx：

```nginx
location /uploads/ {
    proxy_pass http://127.0.0.1:8080;
}
```

这样最容易和当前仓库全部现成页面、接口和部署脚本保持一致。

## 和哪些页一起看最顺

- [环境变量参考](./environment-variables-reference)
- [数据库建表语句](./database-ddl)
- [Nginx 配置参考](./nginx-config-reference)
- [前端概览](/frontend/overview)
- [部署概览](/deployment/overview)
