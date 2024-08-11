package appWebsocketHTTP

import (
	"SystemgeSpawnerTest/topics"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
	"github.com/neutralusername/Systemge/Spawner"
	"github.com/neutralusername/Systemge/Tools"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		topics.ASYNC: func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			err := node.AsyncMessage(topics.ASYNC, "")
			if err != nil {
				panic(err)
			}
			println("sent async message")
			return nil
		},
		topics.SYNC: func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			started := time.Now()
			if responseChannel, err := node.SyncMessage(topics.SYNC, ""); err != nil {
				panic(err)
			} else {
				responseCount := 0
				for {
					_, err := responseChannel.ReceiveResponse()
					if err != nil {
						println("received", responseCount, "sync responses", "in", time.Since(started).Milliseconds(), "ms")
						break
					} else {
						responseCount++
					}
				}
			}
			return nil
		},
		"spawn": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			successfulSpawns := atomic.Uint32{}
			taskGroup := Tools.NewTaskGroup()
			for i := 0; i < 1000; i++ {
				port := app.ports.Add(1)
				if port > 65535 {
					break
				}
				responseChannel, err := node.SyncMessage(Spawner.SPAWN_AND_START_NODE_SYNC, Helpers.JsonMarshal(&Config.NewNode{
					NodeConfig: &Config.Node{
						Name:              Helpers.Uint32ToString(port),
						RandomizerSeed:    Tools.GetSystemTime(),
						InfoLoggerPath:    "logs.log",
						WarningLoggerPath: "logs.log",
						ErrorLoggerPath:   "logs.log",
					},
					SystemgeConfig: &Config.Systemge{
						ProcessMessagesOfEachConnectionSequentially: true,
						ProcessAllMessagesSequentially:              false,

						SyncRequestTimeoutMs:            10000,
						TcpTimeoutMs:                    5000,
						MaxConnectionAttempts:           1,
						ConnectionAttemptDelayMs:        1000,
						StopAfterOutgoingConnectionLoss: true,
						ServerConfig: &Config.TcpServer{
							Port:        uint16(port),
							TlsCertPath: "MyCertificate.crt",
							TlsKeyPath:  "MyKey.key",
						},
						Endpoint: &Config.TcpEndpoint{
							Address: "127.0.0.1:" + Helpers.IntToString(int(port)),
							TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
							Domain:  "example.com",
						},
						EndpointConfigs: []*Config.TcpEndpoint{
							{
								Address: "localhost:60002",
								TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
								Domain:  "example.com",
							},
						},
						IncomingMessageByteLimit: 0,
						MaxPayloadSize:           0,
						MaxTopicSize:             0,
						MaxSyncTokenSize:         0,
						MaxNodeNameSize:          0,
					},
				}))
				if err != nil {
					continue
				}
				taskGroup.AddTask(func() {
					response, err := responseChannel.ReceiveResponse()
					if err != nil {
						return
					}
					if response.GetTopic() == Message.TOPIC_FAILURE {
						return
					}
					err = node.ConnectToNode(&Config.TcpEndpoint{
						Address: "localhost:" + Helpers.IntToString(int(port)),
						TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
						Domain:  "example.com",
					}, true)
					if err != nil {
						panic(err)
					}
					successfulSpawns.Add(1)
				})
			}
			taskGroup.ExecuteTasks()
			println("spawned", successfulSpawns.Load(), "nodes")
			return nil
		},
		"despawn": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			responseChannel, err := node.SyncMessage(Spawner.DESPAWN_ALL_NODES_SYNC, "")
			if err != nil {
				return err
			}
			response, err := responseChannel.ReceiveResponse()
			if err != nil {
				return err
			}
			if response.GetTopic() == Message.TOPIC_FAILURE {
				return fmt.Errorf(response.GetPayload())
			}
			println("despawned all nodes")
			app.ports.Store(32768)
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {

}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
}
