package stream

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	_mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2/log"
	"github.com/reactivex/rxgo/v2"
	"github.com/zgunz42/servercopilot/internal/server/mqtt"
)

type Agent struct {
	obs    rxgo.Observable
	stream chan rxgo.Item
	pubSub *gochannel.GoChannel
	msg    <-chan *message.Message
}

// TempSensor is an interface for temperature sensors
type TempSensor struct {
	Client *mqtt.MqttClient
	Agent  *Agent `optional:"true"`
	pubSub *gochannel.GoChannel
	msg    <-chan *message.Message
}

func CreateTempSens(client *mqtt.MqttClient, PubSub *gochannel.GoChannel) *TempSensor {
	agent := &Agent{}
	agent.stream = make(chan rxgo.Item)
	agent.obs = rxgo.FromChannel(agent.stream)

	msg, err := PubSub.Subscribe(context.Background(), "sensor.temperature")
	if err != nil {
		panic(err)
	}

	return &TempSensor{
		Client: client,
		Agent:  agent,
		msg:    msg,
		pubSub: PubSub,
	}
}

func (s *TempSensor) Sub() rxgo.Observable {
	if s.Agent == nil {
		s.Agent = &Agent{}
		s.Agent.stream = make(chan rxgo.Item)
		s.Agent.obs = rxgo.FromChannel(s.Agent.stream)
	}

	err := s.Client.Subscribe("device/temperature", func(client _mqtt.Client, msg _mqtt.Message) {
		// convert to float64
		data := msg.Payload()

		err := s.pubSub.Publish("sensor.temperature", message.NewMessage(watermill.NewUUID(), message.Payload(data)))
		if err != nil {
			log.Error(err)
		}
		log.Debug("sensor temperature: ", string(data))
		msg.Ack()
	})
	if err != nil {
		panic(err)
	}

	return s.Agent.obs
}

func (s *TempSensor) GetMsg() <-chan *message.Message {
	return s.msg
}

func (s TempSensor) GetTemp(ctx context.Context) float64 {
	if s.Agent == nil {
		return 0.0
	}

	val := s.Agent.obs.Take(1).Observe()

	for {
		select {
		case va := <-val:
			return va.V.(float64)
		case <-ctx.Done():
			return 0
		}
	}

}

func (s TempSensor) GetTempObs() rxgo.Observable {
	if s.Agent == nil {
		return nil
	}

	return s.Agent.obs
}

func (s TempSensor) Close() error {
	if s.Agent == nil {
		return nil
	}

	close(s.Agent.stream)
	return nil
}
