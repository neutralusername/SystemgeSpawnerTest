export class root extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
                WS_CONNECTION: new WebSocket("ws://localhost:8443/ws"),
                constructMessage: (topic, payload) => {
                    return JSON.stringify({
                        topic: topic,
                        payload: payload,
                    });
                },
                setStateRoot: (state) => {
                    this.setState(state)
                }
            },
            (this.state.WS_CONNECTION.onmessage = (event) => {
                let message = JSON.parse(event.data);
                switch (message.topic) {
                    case "error":
                        let errorMessage = message.payload.split("->").reverse()[0]
                        console.log(errorMessage)
                        break;
                    default:
                        console.log("Unknown message topic: " + event.data);
                        break;
                }
            });
        this.state.WS_CONNECTION.onclose = () => {
            setTimeout(() => {
                if (this.state.WS_CONNECTION.readyState === WebSocket.CLOSED) {}
                window.location.reload();
            }, 2000);
        };
        this.state.WS_CONNECTION.onopen = () => {
            let myLoop = () => {
                this.state.WS_CONNECTION.send(this.state.constructMessage("heartbeat", ""));
                setTimeout(myLoop, 15 * 1000);
            };
            setTimeout(myLoop, 15 * 1000);
        };
    }

    render() {
        return React.createElement(
            "div", {
                id: "root",
                onContextMenu: (e) => {
                    e.preventDefault();
                },
                style: {
                    fontFamily: "sans-serif",
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "center",
                    alignItems: "center",
                    touchAction: "none",
                    userSelect: "none",
                },
            },
            React.createElement("button", {
                onClick: () => {
                    this.state.WS_CONNECTION.send(this.state.constructMessage("async", ""));
                },
            }, "async test"),
            React.createElement("button", {
                onClick: () => {
                    this.state.WS_CONNECTION.send(this.state.constructMessage("sync", ""));
                },
            }, "sync test"),
            React.createElement("button", {
                onClick: () => {
                    this.state.WS_CONNECTION.send(this.state.constructMessage("spawn", ""));
                },
            }, "spawn"),
            React.createElement("button", {
                onClick: () => {
                    this.state.WS_CONNECTION.send(this.state.constructMessage("despawn", ""));
                },
            }, "despawn")
        );
    }
}