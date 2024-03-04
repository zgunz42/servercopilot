package stream

import (
	"context"
	"strconv"

	_mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/reactivex/rxgo/v2"
	"github.com/zgunz42/servercopilot/internal/server/mqtt"
)

type HumSensor struct {
	Client *mqtt.MqttClient
	Agent  *Agent `optional:"true"`
}

func CreateHumSens(Client *mqtt.MqttClient) *HumSensor {
	agent := &Agent{}
	agent.stream = make(chan rxgo.Item)
	agent.obs = rxgo.FromChannel(agent.stream)

	return &HumSensor{
		Client: Client,
		Agent:  agent,
	}
}

func (s *HumSensor) Sub() rxgo.Observable {
	err := s.Client.Subscribe("device/humidity", func(client _mqtt.Client, msg _mqtt.Message) {
		// convert to float64
		data := msg.Payload()
		dataStr := string(data)
		temp, err := strconv.ParseFloat(dataStr, 64)
		if err != nil {
			println(err)
			return
		}

		select {
		case s.Agent.stream <- rxgo.Of(temp):
		default:
			// do nothing
		}
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
