package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type rabbitArtifacts struct {
	userEventExchangeName    string
	userImageUpdateQueueName string
}

func setupRabbitMqTopicsAndQueues(channel *amqp.Channel, userEventExchangeName string, userImageEventQueueName string, userImageEventUpdateRoutingKey string) rabbitArtifacts {
	exchangeErr := channel.ExchangeDeclare(userEventExchangeName, "topic", true, false, false, false, nil)
	failOnError(exchangeErr, "failed to declare exchange")

	_, queueDeclarationErr := channel.QueueDeclare(
		userImageEventQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(queueDeclarationErr, "Failed to declare queue")

	bindErr := channel.QueueBind(userImageEventQueueName, userImageEventUpdateRoutingKey, userEventExchangeName, false, nil)
	failOnError(bindErr, "Failed to bind queries queue to topic exchange")

	log.Printf("created topics and queues %s, %s", userImageEventQueueName, userEventExchangeName)

	return rabbitArtifacts{userEventExchangeName: userEventExchangeName, userImageUpdateQueueName: userImageEventQueueName}
}

func handleIncomingImageUpdateMessages(inBound <-chan amqp.Delivery, outBound chan<- ImageUpdate) {
	for msg := range inBound {

		var imageUpdate ImageUpdate
		jsonErr := json.Unmarshal(msg.Body, &imageUpdate)

		if jsonErr != nil {
			log.Println("failed to consume image update message")
			msg.Nack(false, false)
		} else {
			log.Println("successfully consumed image update message")
			outBound <- imageUpdate
			msg.Ack(false)
		}
	}
}

func handleOutgoingImageUpdateMessages(inBound <-chan ImageUpdate) {
	for imageUpdate := range inBound {
		log.Printf("implement me. would send imageupdates to rabbitmq %v\n", imageUpdate)
	}
}

func connectRabbit(conf rabbitConf) *amqp.Connection {
	for {
		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", conf.username, conf.password, conf.hostname, conf.port))
		if err == nil && conn != nil {
			log.Println("connected to rabbitmq")
			return conn
		}
		log.Println(fmt.Sprintf("failed to connect to rabbitmq will retry in %d. current cause: %s", conf.timeout, err))
		time.Sleep(conf.timeout)
	}
}