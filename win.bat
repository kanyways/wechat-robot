@echo off
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
@rem 用于混淆
set GOROOT_FINAL=%GOROOT%
if exist "bin\windows_wechat_robot_amd64.exe" (
    del "bin\windows_wechat_robot_amd64.exe"
) else (
    mkdir bin
)
@rem 输出命令行，前台运行
go build -v -a -gcflags "-N -l" -ldflags "-s -w" -trimpath -o bin/windows_wechat_robot_amd64.exe
upx --brute bin/windows_wechat_robot_amd64.exe