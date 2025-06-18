package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

type AlarmClient interface {
	StartAlarm() error
	StopAlarm() error
}

func ConsumeMessages(sub message.Subscriber, alarmClient AlarmClient) {
	messages, err := sub.Subscribe(context.Background(), "smoke_sensor")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		payload := string(msg.Payload)
		var alarmErr error = nil

		if payload == "0" {
			alarmErr = alarmClient.StopAlarm()
		} else if payload == "1" {
			alarmErr = alarmClient.StartAlarm()
		}

		if alarmErr == nil {
			msg.Ack()
		} else {
			msg.Nack()
		}

	}
}
