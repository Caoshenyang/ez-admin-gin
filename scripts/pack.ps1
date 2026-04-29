# 一键构建并打包部署文件
# 使用方法：在项目根目录执行 .\scripts\pack.ps1
chcp 65001 > $null
$ErrorActionPreference = "Stop"
$PackDir = "deploy-package"

Write-Host ">>> 清理旧打包..." -ForegroundColor Cyan
if (Test-Path $PackDir) { Remove-Item -Recurse -Force $PackDir }

# 构建后端
Write-Host ">>> 编译后端 (Linux amd64)..." -ForegroundColor Cyan
Push-Location server
$env:GOOS = "linux"; $env:GOARCH = "amd64"
go build -ldflags="-s -w" -o "../$PackDir/server" .
Pop-Location

# 构建前端
Write-Host ">>> 构建前端..." -ForegroundColor Cyan
Push-Location admin
pnpm install 2>$null; pnpm build
Copy-Item -Recurse -Force "dist" "../$PackDir/dist"
Pop-Location

# 复制部署配置
Write-Host ">>> 打包配置文件..." -ForegroundColor Cyan
Copy-Item "deploy/compose.server.yml" "$PackDir/"
Copy-Item "deploy/.env.example" "$PackDir/"
Copy-Item "deploy/ez-admin.service" "$PackDir/"
Copy-Item "scripts/setup-server.sh" "$PackDir/"
Copy-Item "scripts/update-server.sh" "$PackDir/"
New-Item -ItemType Directory -Path "$PackDir/nginx" -Force | Out-Null
Copy-Item "deploy/nginx/nginx-native.conf" "$PackDir/nginx/"
New-Item -ItemType Directory -Path "$PackDir/ssl" -Force | Out-Null
New-Item -ItemType Directory -Path "$PackDir/configs" -Force | Out-Null
Copy-Item "server\configs\config.yaml" "$PackDir/configs/"
Copy-Item "server\configs\rbac_model.conf" "$PackDir/configs/"

Write-Host ">>> 生成压缩包..." -ForegroundColor Cyan
if (Test-Path "deploy-package.zip") { Remove-Item -Force "deploy-package.zip" }
Compress-Archive -Path "$PackDir\*" -DestinationPath "deploy-package.zip"

Write-Host ""
Write-Host "✅ 打包完成！上传 deploy-package.zip 到服务器即可。" -ForegroundColor Green
Write-Host "   上传目标：/opt/ez-admin/"
Write-Host "   然后执行："
Write-Host "     mkdir -p /opt/ez-admin && cd /opt/ez-admin"
Write-Host "     unzip ~/deploy-package.zip"
Write-Host "     bash setup-server.sh"
