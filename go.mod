module SystemgeSpawnerTest

go 1.23

toolchain go1.23.0

//replace github.com/neutralusername/Systemge => ../Systemge

require (
	github.com/gorilla/websocket v1.5.3
	github.com/neutralusername/Systemge v0.0.0-20240920150811-762a862539cc
)

require (
	github.com/neutralusername/systemge v0.0.0-20241024160533-a0b1b7996b37 // indirect
	golang.org/x/oauth2 v0.21.0 // indirect
)
