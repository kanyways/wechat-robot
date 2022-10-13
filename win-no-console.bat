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
@rem 不输入命令行，直接后台运行了
go build -gcflags "-N -l" -ldflags "-s -w -H windowsgui" -o bin/windows_wechat_robot_amd64.exe
upx --brute bin/windows_wechat_robot_amd64.exe