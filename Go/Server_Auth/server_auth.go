package main

import (
	"Server_Auth/Strutture"
	database "Server_Auth/database_auth"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Configura le credenziali di autenticazione di base

func middlewareAutenticazioneBase(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ottiene le credenziali di autenticazione dall'intestazione Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			unauthorized(w)
			return
		}

		// Verifica le credenziali di autenticazione di base
		auth := strings.SplitN(authHeader, " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			unauthorized(w)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			fmt.Println("errore decode")
			unauthorized(w)
			return
		}

		//----------------
		var person Strutture.Utenti

		//----------------
		pair := strings.SplitN(string(payload), ":", 2)
		person.Email = pair[0]

		hasher := sha256.New()

		// Scrivere i dati nel hasher
		hasher.Write([]byte(pair[1]))

		// Calcolare l'hash
		hash := hasher.Sum(nil)

		// Convertire l'hash in una rappresentazione esadecimale
		hashString := hex.EncodeToString(hash)

		pair[1] = hashString
		var temp string
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Errore nella lettura del corpo della richiesta", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		fmt.Println(body)
		fmt.Println(string(body))
		// Deserializza il corpo JSON in un oggetto Person

		err = json.Unmarshal(body, &temp)
		if err != nil {
			http.Error(w, "Errore nella deserializzazione JSON", http.StatusBadRequest)
			return
		}
		person.Password = pair[1]
		person.Id_tg = temp
		fmt.Println(person, "qui")

		res, err := database.ReadeUser(person)
		//yfmt.Println(res)
		if res == nil {
			unauthorized(w)
			return
		}
		if err != nil {
			http.Error(w, "Errore Lettura ", http.StatusBadRequest)
			return
		}
		if len(pair) != 2 || pair[0] != res.Email || pair[1] != res.Password {

			fmt.Println(pair[0] + " tua " + res.Email + "  secondacoppia " + pair[1] + "tua " + res.Password)
			unauthorized(w)
			return
		}

		// Chiama la prossima funzione nella catena di middleware
		next.ServeHTTP(w, r)
	})
}

func middlewareRegister(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
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
		//fmt.Println(body)
		err = json.Unmarshal(body, &register)
		if err != nil {
			http.Error(w, "Errore nella deserializzazione JSON", http.StatusBadRequest)
			return
		}

		//----------------
		//fmt.Println(register)
		res, err := database.CreateUser(register)
		if res == nil {
			erroreRegister(w)
		}
		if err != nil {
			http.Error(w, "Errore Creazione ", http.StatusBadRequest)
			return
		}

		// Chiama la prossima funzione nella catena di middleware
		next.ServeHTTP(w, r)
	})
}

func Getidtg(w http.ResponseWriter, r *http.Request) {

	var person Strutture.Utenti

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Errore nella lettura del corpo della richiesta", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Println(body)

	// Deserializza il corpo JSON in un oggetto Person

	person.Email = string(body)
	fmt.Println(person, "qui")

	res, err := database.ReadeUserTg(person)
	//yfmt.Println(res)
	if res == nil {
		unauthorized(w)
		return
	}
	if err != nil {
		http.Error(w, "Errore Creazione ", http.StatusBadRequest)
		return
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	w.Write(jsonData)
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("401 Unauthorized\n"))
}

func gestoreProtetto(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authorized"))
}

func erroreRegister(w http.ResponseWriter) {

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 Bad Request - Errore Registrazione\n"))
}

func okRegister(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registrazione riuscita! Benvenuto nell'API protetta.\n"))
}

func main() {
	fmt.Println("Avvio tra 2 minuti...")
	time.Sleep(120 * time.Second)
	fmt.Println("Avviato")
	database.StartDBUtenti()
	// Configura il router e aggiungi il middleware di autenticazione di base
	http.Handle("/api/protetta", middlewareAutenticazioneBase(http.HandlerFunc(gestoreProtetto)))
	http.Handle("/api/register", middlewareRegister(http.HandlerFunc(okRegister)))
	http.HandleFunc("/api/idtg", Getidtg)
	// Avvia il server sulla porta 8080
	fmt.Println("Server in ascolto su http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
