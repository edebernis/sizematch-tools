package main

import (
    "fmt"
    "github.com/edebernis/sizematch-protobuf/go/items"
    "github.com/golang/protobuf/proto"
    "github.com/streadway/amqp"
    "os"
    "strconv"
    "time"
)

type serviceConfig struct {
    queueName string
}

var serviceConfigs = map[string]serviceConfig{
    "parser": {
        queueName: "sizematch-item-parser-ikea",
    },
    "normalizer": {
        queueName: "sizematch-item-normalizer",
    },
    "saver": {
        queueName: "sizematch-item-saver",
    },
}

func main() {
    if len(os.Args) < 2 {
        fmt.Printf("USAGE : %s <sizematch_service> [requeue] \n", os.Args[0])
        os.Exit(0)
    }

    config := serviceConfigs[os.Args[1]]

    requeue := true
    var err error
    if len(os.Args) == 3 {
        requeue, err = strconv.ParseBool(os.Args[2])
        if err != nil {
            panic("Invalid requeue bool argument: " + err.Error())
        }
    }

    connection, err := amqp.Dial("amqp://user:password@localhost:5672/")
    if err != nil {
        panic("could not connect to RabbitMQ: " + err.Error())
    }
    defer connection.Close()

    channel, err := connection.Channel()
    if err != nil {
        panic("could not create channel: " + err.Error())
    }
    defer channel.Close()

    err = channel.Qos(1, 0, false)
    if err != nil {
        panic("could not set qos on channel: " + err.Error())
    }

    msgs, err := channel.Consume(config.queueName, "", false, false, false, false, nil)
    if err != nil {
        panic("could not consume item: " + err.Error())
    }

    go func(requeue bool) {
        for msg := range msgs {

            err := msg.Nack(false, requeue)
            if err != nil {
                panic("could not nack message: " + err.Error())
            }

            item := items.Item{}
            err = proto.Unmarshal(msg.Body, &item)
            if err != nil {
                panic("could not unmarshal item: " + err.Error())
            }

            fmt.Printf("%+v\n", item)
        }
    }(requeue)

    time.Sleep(1 * time.Second)
}
