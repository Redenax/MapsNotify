package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/kafka-go"
)

func GetKafkaAddress() string {
	brokerAddress := os.Getenv("KafkaAddress")
	// Configura l'indirizzo del broker Kafka
	if brokerAddress == "" {
		brokerAddress = "kafka:9093"
	} else {
		fmt.Printf("The value of the environment variable is: %s\n", brokerAddress)
	}

	fmt.Printf("The value of the environment variable is: %s\n", brokerAddress)

	return brokerAddress
}

func KafkaStartup(province []string) {
	brokerAddress := GetKafkaAddress()

	var a, b, c string
	var inserimento bool
	time.Sleep(15 * time.Second)
	// Configura l'amministratore Kafka
	admin, err := kafka.DialContext(context.Background(), "tcp", brokerAddress)
	if err != nil {
		log.Fatalf("Errore connessione a Kafka: %v", err)
	}
	defer admin.Close()

	// Ottieni la lista dei topic
	topics, err := admin.ReadPartitions()
	if err != nil {
		log.Fatalf("Errore ottenimento lista topic: %v", err)
	}

	fmt.Println("Lista dei topic:")
	for i := range province {

		for j := range province {

			if province[i] != province[j] {
				a = string(province[i])
				b = string(province[j])

				// Specifica il nome del topic da creare
				topic := fmt.Sprintf("%s_%s", a, b)

				// Stampa la lista dei topic

				for _, topiz := range topics {

					c = fmt.Sprintf(topiz.Topic)
					if topic == c {
						fmt.Println("Topic:", c, "Non inserito, gia presente \n")
						inserimento = true
						break
					} else {

						inserimento = false
					}
				}

				if !inserimento {
					// Crea il topic
					err = admin.CreateTopics(kafka.TopicConfig{
						Topic:             topic,
						NumPartitions:     1, // Numero di partizioni del topic
						ReplicationFactor: 1, // Fattore di replicazione del topic
					})
					if err != nil {
						log.Fatalf("Errore creazione topic Kafka: %v", err)
					}
					fmt.Printf("\n Il topic %s è stato creato con successo.\n", topic)
				}
				inserimento = true
			}
		}
	}

}

func KafkaProducer(topic string) {
	brokerAddress := GetKafkaAddress()
	// Set up a connection to Kafka broker
	conn, err := kafka.DialContext(context.Background(), "tcp", brokerAddress)
	if err != nil {
		log.Fatalf("Error connecting to Kafka broker: %v\n", err)
	}
	defer conn.Close()

	// Create a new Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	// Produce a message to the topic
	message := kafka.Message{
		Value: []byte("Ciao, Kafka! " + topic),
	}

	err = writer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Fatalf("Error producing message: %v\n", err)
	}

	fmt.Println("Message sent successfully! in Topic: " + topic)

	// Close the writer
	err = writer.Close()
	if err != nil {
		log.Fatalf("Error closing Kafka writer: %v\n", err)
	}

}
