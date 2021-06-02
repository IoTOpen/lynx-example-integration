package main

import (
	"encoding/json"
	"fmt"
	"github.com/IoTOpen/go-lynx"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var client *lynx.Client

func configure() {
	viper.SetConfigName("lynx-integration")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	viper.SetDefault("api.base", "https://domain.tld")
	viper.SetDefault("api.key", "secret")
	viper.SetDefault("api.broker", "tcp://domain.tld:port")
	viper.SetDefault("lynx.installation_id", 1)

	if err := viper.ReadInConfig(); err != nil {
		_ = viper.SafeWriteConfig()
		log.Fatalln("Config:", err)
	}
}

func lynxClientSetup() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(viper.GetString("api.broker"))
	opts.SetCleanSession(true)
	opts.SetClientID("lynx-integration-example")
	opts.SetConnectTimeout(time.Second)

	client = lynx.NewClient(&lynx.Options{
		Authenticator: lynx.AuthApiKey{
			Key: viper.GetString("api.key"),
		},
		ApiBase:     viper.GetString("api.base"),
		MqttOptions: opts,
	})

	for err := client.MQTTConnect(); err != nil; err = client.MQTTConnect() {
		log.Println("MQTT error connecting:", err)
		time.Sleep(time.Second * 5)
	}
}

func subscribe(fn *lynx.Function, clientID int64) {
	topic := fmt.Sprintf("%d/%s", clientID, fn.Meta["topic_read"])
	if token := client.Mqtt.Subscribe(topic, 2, messageHandler); token.WaitTimeout(time.Second) && token.Error() != nil {
		log.Fatalf("MQTT: Could not subscribe to topic: %s, error: %s", topic, token.Error())
	}
}

func messageHandler(_ mqtt.Client, message mqtt.Message) {
	m := &lynx.Message{}
	if err := json.Unmarshal(message.Payload(), m); err != nil {
		log.Println("Failed to unmarshal payload:", err)
		return
	}
	log.Printf("%s: %v", time.Unix(m.Timestamp, 0), m.Value)
}

func main() {
	configure()
	lynxClientSetup()

	installationID := viper.GetInt64("lynx.installation_id")
	installation, err := client.GetInstallation(installationID)
	if err != nil {
		log.Fatalln("Could not fetch installation:", err)
	}
	devices, err := client.GetDevices(installationID, map[string]string{
		"example.type": "go-lynx",
	})
	if err != nil {
		log.Fatalln("Could not fetch devices:", err)
	} else if len(devices) == 0 {
		log.Fatal("No device found")
	}
	functions, err := client.GetFunctions(installationID, map[string]string{
		"device_id":    fmt.Sprintf("%d", devices[0].ID),
		"example.type": "go-lynx",
	})
	if err != nil {
		log.Fatalln("Could not fetch function:", err)
	}
	subscribe(functions[0], installation.ClientID)
	sigc := make(chan os.Signal)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
}
