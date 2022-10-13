@echo off
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
set GOROOT_FINAL=%GOROOT%
if exist "bin\linux_wechat_robot_amd64" (
    del "bin\linux_wechat_robot_amd64"
) else (
    mkdir bin
)
go build -v -a -gcflags "-N -l" -ldflags "-s -w" -trimpath -o bin/linux_wechat_robot_amd64
upx --brute bin/linux_wechat_robot_amd64