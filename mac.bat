SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -gcflags "-N -l" -ldflags "-s -w" -o darwin_wechat_root_amd64