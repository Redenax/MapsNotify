package database

import (
	"MainServer/Strutture"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbname1 = "routes"
)

func dsn(dbName string) string {
	username := os.Getenv("Username_DB")
	// Configura l'username del database mysql
	if username == "" {
		username = "root"
	} else {
		fmt.Printf("The value of the environment variable is: %s\n", username)
	}
	password := os.Getenv("Password_DB")
	// Configura password del database mysql
	if password == "" {
		password = "tuapassword"
	} else {
		fmt.Printf("The value of the environment variable is: %s\n", password)
	}
	hostname := os.Getenv("Hostname_DB")
	// Configura hostname del database mysql
	if hostname == "" {
		hostname = "mysql_main:3306"
	} else {
		fmt.Printf("The value of the environment variable is: %s\n", hostname)
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func StartDBRoute(province []string) {

	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname1)
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

	//----------------------------------------------------------------------
	db, err = sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}
	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err = db.ExecContext(ctx, "USE "+dbname1)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS province (id INT NOT NULL AUTO_INCREMENT,name VARCHAR(255),PRIMARY KEY (id))")
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

	for i := range province {

		ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()

		res, err = db.ExecContext(ctx, "INSERT INTO province (name) SELECT * FROM (SELECT '"+province[i]+"') AS tmp WHERE NOT EXISTS (SELECT name FROM province WHERE name = '"+province[i]+"') LIMIT 1;")
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

	}

	db.Close()
	//----------------------------------------------------------------------
	db, err = sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}
	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err = db.ExecContext(ctx, "USE "+dbname1)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS routes (name VARCHAR(255),partenza int,destinazione int,email varchar(255),notify tinyint, PRIMARY KEY (partenza, destinazione, email))")
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

func CreateRoute(dati Strutture.Routes) (*string, error) {
	db, err := sql.Open("mysql", dsn(dbname1))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}
	defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname1)
	var id_partenza, id_destinazione int
	err = db.QueryRow("SELECT id FROM province WHERE name='" + dati.Partenza + "'").Scan(&id_partenza)
	err1 := db.QueryRow("SELECT id FROM province WHERE name='" + dati.Destinazione + "'").Scan(&id_destinazione)
	if id_partenza == 0 {
		log.Printf("Partenza non trovata")
		return nil, err
	}
	if id_destinazione == 0 {
		log.Printf("Destinazione non trovata")
		return nil, err1
	}
	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	fmt.Println(dati, "quo")
	dati.Nome = dati.Partenza + "_" + dati.Destinazione
	var query string

	//modificare token e variabile da inserire
	//if dati.Tgid != "nullo" {
	query = fmt.Sprintf("INSERT  INTO routes (name, partenza, destinazione,email,notify) VALUES ('" + dati.Nome + "','" + strconv.Itoa(id_partenza) + "','" + strconv.Itoa(id_destinazione) + "','" + dati.Email + "','0')")

	/*} else {

		log.Printf("Errore tg id non settato\n")
		return nil, err
	}*/

	res, err := db.ExecContext(ctx, query)
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
	a := "ok"
	return &a, nil
}

func DeleteRoute(dati Strutture.Routes) (*string, error) {
	db, err := sql.Open("mysql", dsn(dbname1))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}
	defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname1)
	var id_partenza, id_destinazione int
	err = db.QueryRow("SELECT id FROM province WHERE name='" + dati.Partenza + "'").Scan(&id_partenza)
	err1 := db.QueryRow("SELECT id FROM province WHERE name='" + dati.Destinazione + "'").Scan(&id_destinazione)
	if id_partenza == 0 {
		log.Printf("Partenza non trovata")
		return nil, err
	}
	if id_destinazione == 0 {
		log.Printf("Destinazione non trovata")
		return nil, err1
	}
	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	query := fmt.Sprintf("DELETE FROM routes WHERE partenza='" + strconv.Itoa(id_partenza) + "' AND destinazione='" + strconv.Itoa(id_destinazione) + "' AND email='" + dati.Email + "'")
	//modificare token e variabile da inserire
	fmt.Println(query)
	res, err := db.ExecContext(ctx, query)
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
	a := "ok"
	return &a, nil
}
func EnableRoutes(route Strutture.Routes) ([]Strutture.Routes, error) {

	db, err := sql.Open("mysql", dsn(dbname1))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}
	defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname1)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	query := fmt.Sprintf("UPDATE routes SET notify='%d' WHERE email='"+route.Email+"'", 1)
	//modificare token e variabile da inserire
	fmt.Println(query)

	fmt.Println()
	res, err := db.ExecContext(ctx, query)
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

	db, err = sql.Open("mysql", dsn(dbname1))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}
	defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname1)
	// Esegui la query
	rows, err := db.Query("SELECT name FROM routes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Crea un vettore per immagazzinare i risultati
	var routes []Strutture.Routes

	// Itera sui risultati della query
	for rows.Next() {
		var route Strutture.Routes
		err := rows.Scan(&route.Nome)
		if err != nil {
			log.Fatal(err)
		}
		routes = append(routes, route)
	}

	// Gestisci eventuali errori dopo rows.Next()
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Stampa i risultati
	fmt.Println(routes)

	return routes, nil
}
func DisableRoutes(route Strutture.Routes) (*string, error) {

	db, err := sql.Open("mysql", dsn(dbname1))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}
	defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname1)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	query := fmt.Sprintf("UPDATE routes SET notify='%d' WHERE email='"+route.Email+"'", 0)
	//modificare token e variabile da inserire
	fmt.Println(query)

	fmt.Println()
	res, err := db.ExecContext(ctx, query)
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

	ok := "ok"
	return &ok, nil
}
