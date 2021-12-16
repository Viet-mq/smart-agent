package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"log"
	"net/http"
	"smart_agent/mqtt_cfg"
	"time"
)

var (
	broker                                = "localhost"
	port                                  = 10004
	messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Message \"%s\" received on topic  \"%s\"\n", msg.Payload(), msg.Topic())
	}

	connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		fmt.Println("Connected")
	}

	connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		fmt.Printf("Connection Lost: %s\n", err.Error())
	}
)

type Agent struct {
	ID            string
	clientOpt     *mqtt.ClientOptions
	Notifier      chan []byte
	newClient     chan chan []byte
	closingClient chan chan []byte
	clients       map[chan []byte]bool
}

func newAgent() (agent *Agent) {
	agent = &Agent{
		ID:            uuid.NewString(),
		Notifier:      make(chan []byte, 1),
		newClient:     make(chan chan []byte),
		closingClient: make(chan chan []byte),
		clients:       make(map[chan []byte]bool),
		clientOpt:     mqtt.NewClientOptions(),
	}

	go agent.listen()

	return
}

func (agent *Agent) listen() {

	agent.clientOpt.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	agent.clientOpt.SetClientID("agent")
	agent.clientOpt.SetKeepAlive(60 * time.Second)
	agent.clientOpt.SetPingTimeout(1 * time.Second)
	agent.clientOpt.SetDefaultPublishHandler(messagePubHandler)
	agent.clientOpt.OnConnect = connectHandler
	agent.clientOpt.OnConnectionLost = connectionLostHandler

	client := mqtt.NewClient(agent.clientOpt)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//mqtt_cfg.Subscribe(client)

	for {
		select {
		case s := <-agent.newClient:
			//New client connect
			agent.clients[s] = true
			event := fmt.Sprintf("Client added to %s. Clients connected: %d", agent.ID, len(agent.clients))
			log.Println(event)
			mqtt_cfg.Publish(client, event)

		case s := <-agent.closingClient:
			//Client disconnect
			delete(agent.clients, s)
			event := fmt.Sprintf("Removed client from %s. Client connected: %d", agent.ID, len(agent.clients))
			log.Println(event)
			mqtt_cfg.Publish(client, event)

		case event := <-agent.Notifier:
			log.Printf(string(event))

			mqtt_cfg.Publish(client, string(event))

			//client.Disconnect(100)
		}

	}

}

func (agent *Agent) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	messageChan := make(chan []byte)
	agent.newClient <- messageChan

	defer func() {
		agent.closingClient <- messageChan
	}()

	notify := req.Context().Done()

	go func() {
		<-notify
		agent.closingClient <- messageChan
	}()

	for {
		_, _ = fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)
		flusher.Flush()
	}
}

func main() {
	agent := newAgent()

	go func() {
		eventString := fmt.Sprintf("The time is %v", time.Now())
		agent.Notifier <- []byte(eventString)
	}()

	go log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", agent))
}
