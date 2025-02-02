package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Pokemon struct {
	Name      string   `json:"name"`
	Height    int      `json:"height"`
	Weight    int      `json:"weight"`
	Xp        int      `json:"base_experience"`
	Abilities []string `json:"abilities"`
}

func Init() error {
	var err error
	dsn := "host=localhost port=5432 user=Mario password=0 dbname=pokemon sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS pokemon (
			id UUID PRIMARY KEY,
			name TEXT,
			xp INT,
			weight INT,
			height INT,
			effects TEXT[]
		)`}

	for _, query := range tables {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Error occured while creating tables: %v", err)
		}
	}

	log.Println("Database ready!")
	return nil
}

func AddPokemon(name string) {
	/*
		var pokemon Pokemon
		pokemon.Name = name

		resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name)
		if err != nil {
			log.Fatalf("Error creating http request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response: %v", err)
		}

		json.Unmarshal(body, &pokemon)
	*/
	query := `
		INSERT INTO pokemon ()
	`
	db.Exec(query)
}

func GetOwnedPokemon() []Pokemon {
	return nil
}

func Close() error {
	err := db.Close()
	return err
}
