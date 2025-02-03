package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Can panic
func Init() error {
	var err error
	dsn := "host=localhost port=5432 user=Mario password=0 dbname=pokemon sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Error connecting to the database: %v", err)
		return err
	}

	tables := []string{
		/*`DROP TABLE pokemon`,*/
		/*`DROP TABLE items`,*/
		`CREATE TABLE IF NOT EXISTS pokemon (
			name TEXT PRIMARY KEY,
			count INT
		)`,
		`CREATE TABLE IF NOT EXISTS items (
			name TEXT PRIMARY KEY,
			count INT
		)`,
	}

	for _, query := range tables {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error occured while creating tables: %v", err)
			return err
		}
	}

	log.Println("Database ready!")
	return nil
}

func AddPokemon(name string) {
	query := `
		INSERT INTO pokemon (name, count)
		VALUES ($1, $2)
		ON CONFLICT (name)
		DO UPDATE SET count = pokemon.count + 1
	`
	if _, err := db.Exec(query, name, 1); err != nil {
		log.Printf("AddPokemon db query failed: %v", err)
	}
}

func GetOwnedPokemon() *map[string]int {
	result := make(map[string]int)

	query := `
		SELECT name, count FROM pokemon
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("GetOwnedPokemon failed to retrieve names from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var count int
		if err := rows.Scan(&name, &count); err != nil {
			log.Printf("GetOwnedPokemon failed to scan name from row: %v", err)
		} else {
			result[name] = count
		}
	}

	return &result
}

func AddItem(name string) {
	query := `
	INSERT INTO items (name, count)
	VALUES ($1, $2)
	ON CONFLICT (name)
	DO UPDATE SET count = items.count + 1
`
	if _, err := db.Exec(query, name, 1); err != nil {
		log.Printf("AddItem db query failed: %v", err)
	}
}

func GetOwnedItems() *map[string]int {
	result := make(map[string]int)

	query := `
		SELECT name, count FROM items
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("GetOwnedItems failed to retrieve names from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var count int
		if err := rows.Scan(&name, &count); err != nil {
			log.Printf("GetOwnedItems failed to scan name from row: %v", err)
		} else {
			result[name] = count
		}
	}

	return &result
}

func Close() error {
	err := db.Close()
	return err
}
