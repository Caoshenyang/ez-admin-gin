# EZ Admin Gin Brand Assets

基于你提供的四张品牌板重建并导出的项目物料包。包含可直接用于 Web、文档站、GitHub、App Icon、favicon、社交头像和 Open Graph 的常用格式。

## 目录

- `svg/`：主 logo、横版 logo、上下堆叠 logo、深色版、单色版、白色版、badge、OG card、sprite。
- `png/logo-mark/`：透明背景 logo mark，含 16 / 20 / 24 / 32 / 48 / 64 / 128 / 256 / 512 / 1024。
- `png/lockup-horizontal/`：横版 logo PNG，含 1x / 2x / 3x、dark、mono、white-on-dark。
- `png/lockup-stacked/`：上下堆叠版 logo PNG。
- `png/app-icons/`：圆角白底 App Icon，含 64 / 128 / 180 / 192 / 256 / 512 / 1024。
- `png/favicons/`：favicon PNG 尺寸集合。
- `ico/favicon.ico`：多尺寸 ICO，含 16 / 24 / 32 / 48 / 64 / 128 / 256。
- `png/social/`：头像、Open Graph light/dark cards。
- `png/badges/`：README badge PNG。
- `brand/`：颜色与字体 token：CSS / SCSS / JSON。
- `web/`：`site.webmanifest`、`favicon.ico`、`apple-touch-icon.png`、Android Chrome icons。
- `preview/`：预览页面和预览图。
- `source-reference/`：本次导出使用的参考图与裁剪 mark。

## 推荐使用

- 文档站 Header：`svg/logo-horizontal.svg` 或 `png/lockup-horizontal/logo-horizontal-2x.png`
- 深色 Header：`svg/logo-horizontal-dark.svg`
- GitHub README 徽章：`svg/readme-badge.svg`
- Favicon：`web/favicon.ico` 或 `svg/favicon.svg`
- PWA / Web App：复制 `web/` 内文件到站点 public 根目录
- App / 桌面图标：`png/app-icons/app-icon-512.png` 或 `app-icon-1024.png`
- 社交分享图：`png/social/open-graph-light.png` 或 `open-graph-dark.png`

## 品牌颜色

| Token | Hex |
|---|---|
| Primary Blue | #2563FF |
| Primary Blue Deep | #1D4ED8 |
| Accent Blue | #38D0F8 |
| Deep Navy | #0D1B2A |
| Navy Gray | #1F2937 |
| Slate 600 | #475569 |
| Slate 400 | #94A3B8 |
| Slate 200 | #E2E8F0 |
| Slate 100 | #F1F5F9 |

## 说明

原始文件是位图品牌板，因此这里采用“参考图裁剪 + 矢量轮廓重建”的方式导出，视觉风格、配色和常用尺寸已尽量贴近你发来的设计图。后续如果你有真正的源矢量文件，可以再替换 `svg/logo-mark.svg` 里的 mark path，以获得完全一致的边缘和圆角。
