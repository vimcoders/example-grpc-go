@echo off
chcp 65001 >nul
cls

for /f "delims=" %%a in ('powershell Get-Date -Format "yyyy-MM-dd"') do set BUILD_DATE=%%a
for /f "delims=" %%i in ('git rev-parse --short HEAD') do set GIT_COMMIT=%%i
set IMAGE_ID=%BUILD_DATE%-%GIT_COMMIT%

:: ===========================================
:: 优雅输出（带框 + 对齐 + 干净清爽）
:: ===========================================
echo.
echo  ========================================
echo          📦 镜像构建信息
echo  ========================================
echo  构建日期   :  %BUILD_DATE%
echo  Git 版本   :  %GIT_COMMIT%
echo  镜像标签   :  %IMAGE_ID%
echo  ========================================
echo.

docker compose build
pause