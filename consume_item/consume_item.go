package main

import (
    "fmt"
    "github.com/edebernis/sizematch-protobuf/go/items"
    "github.com/golang/protobuf/proto"
    "github.com/streadway/amqp"
    "os"
    "strconv"
    "strings"
    "time"
)

type serviceConfig struct {
    queueName string
}

var serviceConfigs = map[string]*serviceConfig{
    "parser": {
        queueName: "sizematch-items-parser",
    },
    "normalizer": {
        queueName: "sizematch-items-normalizer",
    },
    "saver": {
        queueName: "sizematch-items-saver",
    },
}

func main() {
    if len(os.Args) < 2 {
        fmt.Printf("USAGE : %s <sizematch_service> <quantity> [requeue] [source] \n", os.Args[0])
        os.Exit(0)
    }

    service := os.Args[1]
    config := serviceConfigs[service]

    quantity, err := strconv.Atoi(os.Args[2])
    if err != nil {
        panic("Invalid quantity integer argument: " + err.Error())
    }

    requeue := true
    if len(os.Args) > 3 {
        requeue, err = strconv.ParseBool(os.Args[3])
        if err != nil {
            panic("Invalid requeue bool argument: " + err.Error())
        }
    }

    if len(os.Args) > 4 && service == "parser" {
        source := os.Args[4]
        queueName := strings.Join(
            []string{serviceConfigs["parser"].queueName, source},
            "-",
        )
        serviceConfigs["parser"].queueName = queueName
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

    err = channel.Qos(quantity, 0, false)
    if err != nil {
        panic("could not set qos on channel: " + err.Error())
    }

    msgs, err := channel.Consume(config.queueName, "", false, false, false, false, nil)
    if err != nil {
        panic("could not consume item: " + err.Error())
    }

    go func(service string, quantity int, requeue bool) {
        consumedMsgs := []amqp.Delivery{}
        for msg := range msgs {
            consumedMsgs = append(consumedMsgs, msg)
            if len(consumedMsgs) == quantity {
                break
            }
        }

        for _, msg := range consumedMsgs {
            err := msg.Nack(false, requeue)
            if err != nil {
                panic("could not nack message: " + err.Error())
            }

            if service == "saver" {
                item := items.NormalizedItem{}
                err = proto.Unmarshal(msg.Body, &item)
                if err != nil {
                    panic("could not unmarshal normalized item: " + err.Error())
                }
                fmt.Printf("%+v\n", item)
            } else {
                item := items.Item{}
                err = proto.Unmarshal(msg.Body, &item)
                if err != nil {
                    panic("could not unmarshal item: " + err.Error())
                }
                fmt.Printf("%+v\n", item)
            }
        }

    }(service, quantity, requeue)

    time.Sleep(100 * time.Millisecond)
}
