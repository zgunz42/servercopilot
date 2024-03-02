package stream

import (
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
	if s.Agent == nil {

	}

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

func (s HumSensor) GetHum() float64 {
	if s.Agent == nil {
		return 0.0
	}

	val := s.Agent.obs.Take(5).Observe()
	avrg := 0.0
	for va := range val {
		avrg += va.V.(float64)
	}

	return avrg / 5
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
