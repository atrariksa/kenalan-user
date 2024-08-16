package main

import (
	"flag"
	"os"

	_ "modernc.org/sqlite"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", os.Getenv("WORKDIR")+"cmd/migrations", "directory with migration files")
	args  = []string{}
)

func main() {
	flags.Parse(os.Args[1:])
	args = flags.Args()
	if args[0] == "sqlite3" {
		SQLITE3()
	}
}
