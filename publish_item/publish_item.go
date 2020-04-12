package main

import (
    "github.com/edebernis/sizematch-protobuf/build/go/items"
    "github.com/golang/protobuf/proto"
    "github.com/streadway/amqp"
)

var exchangeName = "sizematch-items"

/* Parser configuration */
var routingKey = "items.parse.ikea"
var queueName = "sizematch-item-parser-ikea"

/* Normalizer configuration
var routingKey = "items.normalize"
var queueName = "sizematch-item-normalizer"
*/
/* Saver configuration
var routingKey = "items.save"
var queueName = "sizematch-item-saver"
*/

var unparsedItem = items.Item{
    Source: "ikea",
    Lang:   items.Lang_FR,
    Urls: []string{
        "https://www.ikea.com/fr/fr/p/aptitlig-planche-a-decouper-bambou-80233430",
    },
}

var parsedItem = items.Item{
    Id:     "123",
    Source: "ikea",
    Lang:   items.Lang_EN,
    Urls: []string{
        "https://www.ikea.com/gb/en/p/leifarne-swivel-chair-dark-yellow-balsberget-white-s29301700/",
    },
    Name:        "LEIFARNE Swivel chair - dark yellow, Balsberget white",
    Description: "LEIFARNE Swivel chair - dark yellow, Balsberget white. You sit comfortably thanks to the restful flexibility of the scooped seat and shaped back. The self-adjusting plastic feet adds stability to the chair. A special surface treatment on the seat prevents you from sliding.",
    Categories: []string{
        "Dining chairs",
    },
    ImageUrls: []string{
        "https://www.ikea.com/gb/en/images/products/leifarne-swivel-chair__0742962_PE742882_S5.JPG",
        "https://www.ikea.com/gb/en/images/products/leifarne-swivel-chair__0802274_PH165465_S5.JPG",
        "https://www.ikea.com/gb/en/images/products/leifarne-swivel-chair__0799469_PH165466_S5.JPG",
        "https://www.ikea.com/gb/en/images/products/leifarne-swivel-chair__0742966_PE742886_S5.JPG",
        "https://www.ikea.com/gb/en/images/products/leifarne-swivel-chair__0742965_PE742883_S5.JPG",
        "https://www.ikea.com/gb/en/images/products/leifarne-swivel-chair__0787512_PE763253_S5.JPG",
    },
    Dimensions: map[string]string{
        "Tested for: ":       "100 kg",
        "Width: ":            "69 cm",
        "Depth: ":            "69 cm",
        "Max. height: ":      "87 cm",
        "Seat width: ":       "45 cm",
        "Seat depth: ":       "36 cm",
        "Min. seat height: ": "41 cm",
        "Max. seat height: ": "51 cm",
    },
    Price:         31.0,
    PriceCurrency: "GBP",
}

var normalizedItem = items.NormalizedItem{
    Id:     "123",
    Source: "ikea",
    Lang:   items.Lang_EN,
    Urls: []string{
        "https://www.ikea.com/gb/en/p/leifarne-swivel-chair-dark-yellow-balsberget-white-s29301700/",
    },
    Name:        "LEIFARNE Swivel chair - dark yellow, Balsberget white",
    Description: "LEIFARNE Swivel chair - dark yellow, Balsberget white. You sit comfortably thanks to the restful flexibility of the scooped seat and shaped back. The self-adjusting plastic feet adds stability to the chair. A special surface treatment on the seat prevents you from sliding.",
    Categories: []string{
        "Dining chairs",
    },
    ImageUrls: []string{
        "https://www.ikea.com/gb/en/images/products/leifarne-swivel-chair__0742962_PE742882_S5.JPG",
        "https://www.ikea.com/gb/en/images/products/leifarne-swivel-chair__0802274_PH165465_S5.JPG",
        "https://www.ikea.com/gb/en/images/products/leifarne-swivel-chair__0799469_PH165466_S5.JPG",
        "https://www.ikea.com/gb/en/images/products/leifarne-swivel-chair__0742966_PE742886_S5.JPG",
        "https://www.ikea.com/gb/en/images/products/leifarne-swivel-chair__0742965_PE742883_S5.JPG",
        "https://www.ikea.com/gb/en/images/products/leifarne-swivel-chair__0787512_PE763253_S5.JPG",
    },
    Dimensions: []*items.Dimension{
        {
            Name:  items.Dimension_HEIGHT,
            Value: 87.0,
            Unit:  items.Dimension_CM,
        },
        {
            Name:  items.Dimension_WIDTH,
            Value: 69.0,
            Unit:  items.Dimension_CM,
        },
        {
            Name:  items.Dimension_DEPTH,
            Value: 69.0,
            Unit:  items.Dimension_CM,
        },
    },
    Price: &items.Price{
        Amount:   39,
        Currency: items.Price_GBP,
    },
}

func main() {
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

    body, err := proto.Marshal(&unparsedItem)
    if err != nil {
        panic("could not marshal item: " + err.Error())
    }

    msg := amqp.Publishing{
        ContentType: "application/protobuf",
        AppId:       "sizematch-tool-publish-item",
        Body:        body,
    }

    err = channel.Publish(exchangeName, routingKey, true, false, msg)
    if err != nil {
        panic("could not publish item: " + err.Error())
    }
}
