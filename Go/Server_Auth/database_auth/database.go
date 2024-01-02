package database

import (
	"Server_Auth/Strutture"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbname = "utenti"
)

func dsn(dbName string) string {
	username := os.Getenv("Username_DB")
	// Configura l'indirizzo del broker Kafka
	if username == "" {
		username = "root"
	} else {
		fmt.Printf("The value of the environment variable is: %s\n", username)
	}
	password := os.Getenv("Password_DB")
	// Configura l'indirizzo del broker Kafka
	if password == "" {
		password = "tuapassword"
	} else {
		fmt.Printf("The value of the environment variable is: %s\n", password)
	}
	hostname := os.Getenv("Hostname_DB")
	// Configura l'indirizzo del broker Kafka
	if hostname == "" {
		hostname = "mysqlaut:3307"
	} else {
		fmt.Printf("The value of the environment variable is: %s\n", hostname)
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func StartDBUtenti() {

	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}
	db.SetConnMaxLifetime(time.Minute * 5) // Set maximum connection lifetime
	db.SetMaxOpenConns(10)                 // Set maximum open connections
	db.SetMaxIdleConns(5)                  // Set maximum idle connections
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}
	log.Printf("rows affected: %d\n", no)
	db.Close()

	//----------------------------------------------------------------------------------------
	db, err = sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}
	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err = db.ExecContext(ctx, "USE "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS users (name VARCHAR(255),Cognome VARCHAR(255),Email VARCHAR(255),Password VARCHAR(255),id_tg VARCHAR(255) DEFAULT '',active TINYINT,PRIMARY KEY (Email))")
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	no, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}
	log.Printf("rows affected: %d\n", no)
	db.Close()

}

func CreateUser(dati Strutture.Utenti) (*Strutture.Utenti, error) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
	dati.Id_tg = "nullo"
	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	fmt.Println(dati)
	//modificare token e variabile da inserire
	res, err := db.ExecContext(ctx, "INSERT INTO users (name, Cognome, Email,Password,id_tg,active) VALUES ('"+dati.Nome+"','"+dati.Cognome+"','"+dati.Email+"','"+dati.Password+"','"+dati.Id_tg+"','1')")
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return nil, err
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return nil, err
	}
	log.Printf("rows affected: %d\n", no)
	return &dati, nil
}

func ReadeUser(dati Strutture.Utenti) (*Strutture.Utenti, error) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
	var utente Strutture.Utenti
	var temp = dati.Id_tg
	err = db.QueryRow("SELECT * FROM users WHERE Email='"+dati.Email+"' and Password='"+dati.Password+"'").Scan(&utente.Nome, &utente.Cognome, &utente.Email, &utente.Password, &utente.Id_tg, &utente.Active)
	if temp != "nullo" && utente.Id_tg == "nullo" {
		fmt.Println(dati, "qua")
		UpdateUserTg(dati, temp)
	}

	if err != nil {
		log.Printf("Error %s when creating DB\n", err)

		return nil, err
	}

	return &utente, nil

}

func ReadeUserTg(dati Strutture.Utenti) (*Strutture.Utenti, error) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
	var utente Strutture.Utenti

	err = db.QueryRow("SELECT id_tg FROM users WHERE Email='" + dati.Email + "'").Scan(&utente.Id_tg)

	if err != nil {
		log.Printf("Error %s when creating DB\n", err)

		return nil, err
	}

	return &utente, nil

}

func DeleteUser(dati Strutture.Utenti) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return
	}
	log.Printf("Connected to DB %s successfully\n", dbname)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	//Aggiungere ricerca prima per verificare se esiste
	res, err := db.ExecContext(ctx, "DELETE FROM users WHERE Email='"+dati.Email+"'")
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}
	log.Printf("rows affected: %d\n", no)

}

func UpdateUserTg(dati Strutture.Utenti, temp string) {

	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return
	}
	log.Printf("Connected to DB %s successfully\n", dbname)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	fmt.Println("UPDATE users SET id_tg='" + temp + "' WHERE Email='" + dati.Email + "' AND id_tg='nullo'")
	res, err := db.ExecContext(ctx, "UPDATE users SET id_tg='"+temp+"' WHERE Email='"+dati.Email+"' AND id_tg='nullo'")
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}
	log.Printf("rows affected: %d\n", no)

}

func UpdateUser(dati Strutture.Utenti) {

	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return
	}
	log.Printf("Connected to DB %s successfully\n", dbname)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.ExecContext(ctx, "UPDATE users SET Nome='"+dati.Nome+"',Cognome='"+dati.Cognome+"',Email='"+dati.Email+"', WHERE Email='"+dati.Email+"'")
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}
	log.Printf("rows affected: %d\n", no)

}
