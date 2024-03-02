package mqtt

import (
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	opts   *mqtt.ClientOptions
	client mqtt.Client
	topics []string
}

func NewMqttClient() *MqttClient {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(os.Getenv("MQTT_BROKER_URL"))
	opts.Username = os.Getenv("MQTT_USERNAME")
	opts.Password = os.Getenv("MQTT_PASSWORD")
	return &MqttClient{
		opts: opts,
	}
}

func generateRandomAlphaNum(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func (m *MqttClient) Connect(attemp int) error {

	if m.client != mqtt.Client(nil) && m.client.IsConnected() {
		return nil
	}

	clientID := os.Getenv("MQTT_CLIENT_ID") + "-" + generateRandomAlphaNum(5)
	m.opts.SetClientID(clientID)

	client := mqtt.NewClient(m.opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		if attemp > 0 {
			return m.Connect(attemp - 1)
		}

		return token.Error()
	}

	m.client = client
	return nil
}

func (m *MqttClient) Disconnect() error {

	for _, topic := range m.topics {
		err := m.Unsubscribe(topic)
		if err != nil {
			return err
		}
	}

	if m.client == mqtt.Client(nil) {
		return nil
	}

	m.client.Disconnect(0)
	return nil
}

func (m *MqttClient) Publish(topic string, payload interface{}) error {
	token := m.client.Publish(topic, 0, false, payload)
	return token.Error()
}

func (m *MqttClient) Subscribe(topic string, callback mqtt.MessageHandler) error {
	m.topics = append(m.topics, topic)
	token := m.client.Subscribe(topic, 0, callback)
	return token.Error()
}

func (m *MqttClient) Unsubscribe(topic string) error {
	token := m.client.Unsubscribe(topic)
	return token.Error()
}
