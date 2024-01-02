package main

import (
	"MainServer/Strutture"
	database "MainServer/database_main"
	"MainServer/kafka"
	"MainServer/routes"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func handleRegisterRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito+2", http.StatusMethodNotAllowed)
		return
	}

	// Leggi il corpo della richiesta
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Errore nella lettura del corpo della richiesta", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Deserializza il corpo JSON in un oggetto
	var register Strutture.Utenti
	err = json.Unmarshal(body, &register)
	if err != nil {
		http.Error(w, "Errore nella deserializzazione JSON+1", http.StatusBadRequest)
		return
	}

	hasher := sha256.New()

	// Scrivere i dati nel hasher
	hasher.Write([]byte(register.Password))

	// Calcolare l'hash
	hash := hasher.Sum(nil)

	// Convertire l'hash in una rappresentazione esadecimale
	hashString := hex.EncodeToString(hash)

	register.Password = hashString

	jsonData, err := json.Marshal(register)
	if err != nil {
		fmt.Println("Errore nella serializzazione JSON:", err)
		return
	}
	// Creare la richiesta HTTP
	req, err := http.NewRequest("POST", "http://authserver:8081/api/register", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Errore nella creazione della richiesta:", err)
		return
	}

	// Eseguire la richiesta HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errore nell'esecuzione della richiesta:", err)
		return
	}
	defer resp.Body.Close()

	// Leggere e stampare la risposta
	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Errore nella lettura della risposta:", err)
		return
	}

	//fmt.Println(string(body) + "body;")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func handleAuthRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito+1", http.StatusMethodNotAllowed)
		return
	}

	// Leggi il corpo della richiesta
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Errore nella lettura del corpo della richiesta", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	fmt.Println(body)
	fmt.Println(string(body))
	// Deserializza il corpo JSON in un oggetto Person
	var Auth Strutture.Authentication

	err = json.Unmarshal(body, &Auth)
	if err != nil {
		http.Error(w, "Errore nella deserializzazione JSON", http.StatusBadRequest)
		return
	}
	fmt.Println(Auth)
	// Credenziali per l'autenticazione di base
	username := Auth.Email
	password := Auth.Password

	resp1, err1 := json.Marshal(Auth.Id_tg)
	if err1 != nil {
		http.Error(w, "Errore nella serializzazione JSON", http.StatusBadRequest)
		return
	}
	// Creare la stringa di autorizzazione di base
	authString := fmt.Sprintf("%s:%s", username, password)
	base64AuthString := base64.StdEncoding.EncodeToString([]byte(authString))
	authHeader := fmt.Sprintf("Basic %s", base64AuthString)

	// Creare la richiesta HTTP
	req, err := http.NewRequest("POST", "http://authserver:8081/api/protetta", bytes.NewBuffer([]byte(resp1)))
	if err != nil {
		fmt.Println("Errore nella creazione della richiesta:", err)
		return
	}

	// Aggiungere l'intestazione di autorizzazione
	req.Header.Add("Authorization", authHeader)

	// Eseguire la richiesta HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errore nell'esecuzione della richiesta+1:", err)
		return
	}
	defer resp.Body.Close()

	// Leggere e stampare la risposta
	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Errore nella lettura della risposta:", err)
		return
	}

	fmt.Println(string(body))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func handleRegisterRouteRequest(w http.ResponseWriter, r *http.Request) {
	var route Strutture.Routes
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito+3", http.StatusMethodNotAllowed)
		return
	}

	// Leggi il corpo della richiesta
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Errore nella lettura del corpo della richiesta", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Deserializza il corpo JSON in un oggetto

	err = json.Unmarshal(body, &route)
	if err != nil {
		http.Error(w, "Errore nella deserializzazione JSON+1", http.StatusBadRequest)
		return
	}
	fmt.Println(route)

	fmt.Println(route)
	resp2, err1 := database.CreateRoute(route)
	if err1 != nil || resp2 == nil {
		http.Error(w, "Errore nella richiesta", http.StatusInternalServerError)
		return
	}
	//fmt.Println(string(body) + "body;")

	w.Write([]byte(*resp2))
}

func handleDeleteRouteRequest(w http.ResponseWriter, r *http.Request) {
	var route Strutture.Routes
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito+3", http.StatusMethodNotAllowed)
		return
	}

	// Leggi il corpo della richiesta
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Errore nella lettura del corpo della richiesta", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Deserializza il corpo JSON in un oggetto

	fmt.Println(string(body))

	// Deserializza il corpo JSON in un oggetto

	err = json.Unmarshal(body, &route)
	if err != nil {
		http.Error(w, "Errore nella deserializzazione JSON+1", http.StatusBadRequest)
		return
	}
	fmt.Println(route)

	resp2, err1 := database.DeleteRoute(route)
	if err1 != nil {
		http.Error(w, "Errore nella richiesta", http.StatusInternalServerError)
		return
	}
	//fmt.Println(string(body) + "body;")

	w.Write([]byte(*resp2))
}

func handleEnableRouteRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito+3", http.StatusMethodNotAllowed)
		return
	}

	// Leggi il corpo della richiesta
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Errore nella lettura del corpo della richiesta", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Deserializza il corpo JSON in un oggetto
	var route Strutture.Routes
	fmt.Println(string(body))

	// Deserializza il corpo JSON in un oggetto

	err = json.Unmarshal(body, &route)
	if err != nil {
		http.Error(w, "Errore nella deserializzazione JSON+1", http.StatusBadRequest)
		return
	}
	fmt.Println(route)

	resp2, err1 := database.EnableRoutes(route)
	if err1 != nil {
		http.Error(w, "Errore nella richiesta", http.StatusInternalServerError)
		return
	}
	//fmt.Println(string(body) + "body;")
	// Marshal dei dati in JSON
	jsonData, err := json.Marshal(resp2)
	if err != nil {
		http.Error(w, "Errore nel marshalling dei dati in JSON", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(jsonData))
}

func handleDisableRouteRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito+3", http.StatusMethodNotAllowed)
		return
	}

	// Leggi il corpo della richiesta
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Errore nella lettura del corpo della richiesta", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Deserializza il corpo JSON in un oggetto
	var route Strutture.Routes
	fmt.Println(string(body))

	// Deserializza il corpo JSON in un oggetto

	err = json.Unmarshal(body, &route)
	if err != nil {
		http.Error(w, "Errore nella deserializzazione JSON+1", http.StatusBadRequest)
		return
	}
	fmt.Println(route)

	resp2, err1 := database.DisableRoutes(route)
	if err1 != nil {
		http.Error(w, "Errore nella richiesta", http.StatusInternalServerError)
		return
	}
	//fmt.Println(string(body) + "body;")

	w.Write([]byte(*resp2))
}

func CallMapsApi(province []string) {
	var a, b string
	for i := range province {

		for j := range province {

			if province[i] != province[j] {
				a = string(province[i])
				b = string(province[j])
				resp := routes.Routing(a, b)
				for _, s := range resp.Routes {
					fmt.Println("Durata Attuale: ", (s.Duration.Seconds)/60, " minuti", "\nDurata Tipica: ", (s.StaticDuration.Seconds)/60, " minuti")
					// Specifica il nome del topic da creare
					topic := fmt.Sprintf("%s_%s", a, b)

					kafka.KafkaProducer(topic)

				}
			}
		}

	}
}
func main() {
	province := []string{"Agrigento", "Caltanissetta", "Catania", "Enna", "Messina", "Palermo", "Ragusa", "Siracusa", "Trapani"}
	intervallo := time.Hour

	fmt.Println("Avvio tra 2 minuti...")
	time.Sleep(120 * time.Second)
	fmt.Println("Avviato")

	kafka.KafkaStartup(province)
	database.StartDBRoute(province)

	// Goroutine per eseguire la tua funzione ogni ora
	go func() {
		CallMapsApi(province)
		for {
			// Attendi per un'ora
			time.Sleep(intervallo)

			// Esegui la tua funzione ogni ora
			CallMapsApi(province)
		}
	}()

	go func() {
		// Configura il router per autenticazione di base
		http.HandleFunc("/api/v1/authentication", handleAuthRequest)
		http.HandleFunc("/api/v1/register", handleRegisterRequest)

		// Configura il router per autenticazione di base
		http.HandleFunc("/api/v1/deletesRoute", handleDeleteRouteRequest)
		http.HandleFunc("/api/v1/registerRoute", handleRegisterRouteRequest)
		http.HandleFunc("/api/v1/enableRoute", handleEnableRouteRequest)
		http.HandleFunc("/api/v1/disableRoute", handleDisableRouteRequest)
		// Avvia il server sulla porta 25535
		fmt.Println("Server in ascolto su http://localhost:25536")
		http.ListenAndServe(":25536", nil)
	}()
	select {}
}
