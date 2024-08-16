package main

import (
	"log"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func SQLITE3() {

	if len(args) < 3 {
		flags.Usage()
		return
	}

	dbstring, command := args[1], args[2]

	db, err := goose.OpenDBWithDriver("sqlite", dbstring)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
