package mqtt_cfg

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Subscribe(client mqtt.Client) {
	topic := "topic/test"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s\n", topic)
}

func Publish(client mqtt.Client, message string) {
	client.Publish("topic/test", 0, false, message)
}
