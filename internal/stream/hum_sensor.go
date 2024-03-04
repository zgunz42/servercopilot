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

type HumSensor struct {
	Client *mqtt.MqttClient
	Agent  *Agent `optional:"true"`
	pubSub *gochannel.GoChannel
	msg    <-chan *message.Message
}

func CreateHumSens(Client *mqtt.MqttClient, PubSub *gochannel.GoChannel) *HumSensor {
	agent := &Agent{}
	agent.stream = make(chan rxgo.Item)
	agent.obs = rxgo.FromChannel(agent.stream)

	msg, err := PubSub.Subscribe(context.Background(), "sensor.humidity")
	if err != nil {
		panic(err)
	}

	return &HumSensor{
		Client: Client,
		Agent:  agent,
		msg:    msg,
		pubSub: PubSub,
	}
}

func (s *HumSensor) GetMsg() <-chan *message.Message {
	return s.msg
}

func (s *HumSensor) Sub() rxgo.Observable {
	err := s.Client.Subscribe("device/humidity", func(client _mqtt.Client, msg _mqtt.Message) {
		// convert to float64
		data := msg.Payload()
		// temp, err := strconv.ParseFloat(dataStr, 64)
		// if err != nil {
		// 	println(err)
		// 	return
		// }

		err := s.pubSub.Publish("sensor.humidity", message.NewMessage(watermill.NewUUID(), message.Payload(data)))
		if err != nil {
			log.Error(err)
		}

		log.Debug("sensor humidity: ", string(data))
		msg.Ack()
	})
	if err != nil {
		panic(err)
	}

	return s.Agent.obs
}

func (s HumSensor) GetHum(ctx context.Context) float64 {
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

func (s HumSensor) GetHumObs() rxgo.Observable {
	if s.Agent == nil {
		return nil
	}

	return s.Agent.obs
}

func (s HumSensor) Close() error {
	if s.Agent == nil {
		return nil
	}

	close(s.Agent.stream)
	return nil
}
