build:
	CGO_ENABLED=1 GOARCH=386 go build -buildmode=c-shared -o goHttpServer.dll httpServer.go routers.go channels.go logger.go

build_cmd:
	set CGO_ENABLED=1&&set GOARCH=386&&go build -buildmode=c-shared -o goHttpServer.dll httpServer.go routers.go channels.go logger.go