@echo off

:: 设置环境变量
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0

:: 设置项目信息
set PROJECT_DIR=%~dp0
set BINARY_NAME=report
set RELEASE_DIR=%PROJECT_DIR%bin

:: 创建发布目录（如果不存在）
if not exist "%RELEASE_DIR%" (
    mkdir "%RELEASE_DIR%"
)

:: 清理旧的二进制文件
if exist "%RELEASE_DIR%\%BINARY_NAME%.exe" (
    del /F /Q "%RELEASE_DIR%\%BINARY_NAME%.exe"
)
if exist "%RELEASE_DIR%\config.sample.yaml" (
    del /F /Q "%RELEASE_DIR%\config.sample.yaml"
)

:: 执行编译命令
echo Building release for Windows...
go build -o %RELEASE_DIR%\%BINARY_NAME%.exe -ldflags="-s -w"
copy %PROJECT_DIR%config.sample.yaml %RELEASE_DIR%

:: 检查编译结果
if exist "%RELEASE_DIR%\%BINARY_NAME%.exe" (
    echo Build successful.
) else (
    echo Error: Build failed.
)

pause
