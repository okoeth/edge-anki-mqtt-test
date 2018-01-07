package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	anki "github.com/okoeth/edge-anki-base"
)

var mode = flag.String("M", "W", "mode is read or write")
var msgs = flag.Int("N", 5, "number of messages to write")

/*
//define a function for the default message handler
var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TIME: %v\n", time.Now())
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}
*/

func main() {
	flag.Parse()
	fmt.Printf("Mode: %s\n", *mode)
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("go-simple" + *mode)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if *mode == "W" {
		for i := 0; i < *msgs; i++ {
			s := anki.Status{
				MsgID:        39,
				MsgTimestamp: time.Now(),
				PosOptions: []anki.PosOption{
					{
						OptTileNo:      1,
						OptProbability: 10,
					},
					{
						OptTileNo:      4,
						OptProbability: 90,
					},
				},
			}
			sj, err := json.Marshal(s)
			if err != nil {
				fmt.Printf("ERROR: %v", err)
			}
			if token := c.Publish("go-mqtt/sample", 0, false, sj); token.Wait() && token.Error() != nil {
				fmt.Printf("ERROR: %v", token.Error())
			}
		}
	}

	if *mode == "R" {
		ch := make(chan anki.Status)
		if token := c.Subscribe("go-mqtt/sample", 0, func(client mqtt.Client, msg mqtt.Message) {
			s := anki.Status{}
			err := json.Unmarshal(msg.Payload(), &s)
			if err != nil {
				fmt.Printf("ERROR: %v", err)
			}
			ch <- s
		}); token.Wait() && token.Error() != nil {
			fmt.Printf("%v", token.Error())
		}
		for {
			t := <-ch
			fmt.Printf("Message Received with latency of %dms\n", (time.Now().UnixNano()-t.MsgTimestamp.UnixNano())/1000000)
		}
	}
	c.Disconnect(250)
}
