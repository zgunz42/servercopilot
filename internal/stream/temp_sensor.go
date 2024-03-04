package stream

import (
	"context"
	"strconv"

	_mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/reactivex/rxgo/v2"
	"github.com/zgunz42/servercopilot/internal/server/mqtt"
)

type Agent struct {
	obs    rxgo.Observable
	stream chan rxgo.Item
}

// TempSensor is an interface for temperature sensors
type TempSensor struct {
	Client *mqtt.MqttClient
	Agent  *Agent `optional:"true"`
}

func CreateTempSens(client *mqtt.MqttClient) *TempSensor {
	agent := &Agent{}
	agent.stream = make(chan rxgo.Item)
	agent.obs = rxgo.FromChannel(agent.stream)

	return &TempSensor{
		Client: client,
		Agent:  agent,
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
