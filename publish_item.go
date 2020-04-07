package main

import (
    "github.com/edebernis/sizematch-protobuf/build/go/items"
    "github.com/golang/protobuf/proto"
    "github.com/streadway/amqp"
)

func main() {
    exchangeName := "sizematch-items"
    routingKey := "items.normalize"
    queueName := "sizematch-item-normalizer"
    appID := "sizematch-item-normalizer"

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

    err = channel.ExchangeDeclare(exchangeName, "direct", false, false, false, false, nil)
    if err != nil {
        panic("could not declare exchange: " + err.Error())
    }

    _, err = channel.QueueDeclare(queueName, false, false, false, false, nil)
    if err != nil {
        panic("could not declare queue: " + err.Error())
    }

    err = channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
    if err != nil {
        panic("could not bind queue: " + err.Error())
    }

    item := items.Item{
        Id: "1",
    }

    body, err := proto.Marshal(&item)
    if err != nil {
        panic("could not marshal item: " + err.Error())
    }

    msg := amqp.Publishing{
        ContentType: "application/protobuf",
        AppId:       appID,
        Body:        body,
    }

    err = channel.Publish(exchangeName, routingKey, true, false, msg)
    if err != nil {
        panic("could not publish item: " + err.Error())
    }
}
