module SystemgeSpawnerTest

go 1.23

toolchain go1.23.0

//replace github.com/neutralusername/Systemge => ../Systemge

require (
	github.com/gorilla/websocket v1.5.3
	github.com/neutralusername/Systemge v0.0.0-20240912114254-9d79bc5d9c59
)

require golang.org/x/oauth2 v0.21.0 // indirect
