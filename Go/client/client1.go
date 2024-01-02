// socket-client project main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"
)

type Utenti struct {
	Nome    string `json:"nome"`
	Cognome string `json:"cognome"`
	Email   string `json:"email"`
}

type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

const (
	brokerAddress = "localhost:9092"
	topic         = "Palermo_Catania"
	groupID       = "3114"
)

func main() {

	// Configura un segnale per intercettare l'interruzione da tastiera
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Crea un nuovo reader (consumer)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerAddress},
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.LastOffset,
	})

	// Chiudi il reader quando il programma termina
	defer reader.Close()

	// Loop per leggere i messaggi
	for {

		// Configura un timeout di 10 secondi
		timeout := 86400 * time.Second

		// Configura un contesto con timeout
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Errore durante la lettura del messaggio: %v\n", err)
			break
		}

		fmt.Printf("Messaggio ricevuto: %s\n", string(msg.Value))

		// Esci dal loop se il programma riceve un segnale di interruzione
		select {
		case <-signals:
			fmt.Println("Programma terminato.")
			return
		default:
		}
	}
	/*

		// URL dell'API della NGO
		apiURL := "http://localhost:23356/api/v1/authentication" // Assicurati di utilizzare l'URL corretto del tuo server

		// Creazione dell'oggetto Person
		person := Utenti{
			Nome:    "John Doe",
			Cognome: "30",
			Email:   "a@a.it",
		}

		// Serializza l'oggetto in formato JSON
		payload, err := json.Marshal(person)
		if err != nil {
			fmt.Printf("Errore durante la serializzazione JSON: %v\n", err)
			return
		}

		// Esempio di richiesta POST con dati JSON
		response, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Printf("Errore durante la richiesta: %v\n", err)
			return
		}
		defer response.Body.Close()

		// Leggi il corpo della risposta
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Errore durante la lettura del corpo della risposta: %v\n", err)
			return
		}

		// Deserializza il corpo JSON della risposta
		var jsonResponse map[string]string
		err = json.Unmarshal(body, &jsonResponse)
		if err != nil {
			fmt.Printf("Errore nella deserializzazione JSON della risposta: %v\n", err)
			return
		}

		// Stampa la risposta
		fmt.Printf("Risposta dalla ONG: %+v\n", jsonResponse)

		//establish connection
		/*connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
		if err != nil {
			panic(err)
		}
		///send some data
		_, err = connection.Write([]byte("Hello Server! Greetings."))
		buffer := make([]byte, 1024)
		mLen, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		fmt.Println("Received: ", string(buffer[:mLen]))
		defer connection.Close()*/
}
