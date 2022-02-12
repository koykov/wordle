package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/koykov/wordle"
)

var (
	fdb string
	db  wordle.DB

	fpat, fneg string
)

func init() {
	var err error

	rf := func(v *string, names []string, value, usage string) {
		for i := range names {
			flag.StringVar(v, names[i], value, usage)
		}
	}
	rf(&fdb, []string{"database", "db"}, "", "Path to nouns5 database.")
	rf(&fpat, []string{"pattern", "pat"}, "", "Pattern to match words.")
	rf(&fneg, []string{"negative", "neg"}, "", "List of chars to exclude.")
	flag.Parse()

	if len(fdb) == 0 {
		log.Fatalln("param -database is required")
	}
	if len(fpat) < wordle.WordSize {
		log.Fatalln("param -pattern is empty or too short")
	}
	if _, err = os.Stat(fdb); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("database file '%s' doesn't exists\n", fdb)
	}
	if err = db.Load(fdb); err != nil {
		log.Fatalf("error '%s' caught on database load '%s'\n", err.Error(), fdb)
	}
}

func main() {
	dst := make([]string, 0, 50)
	var err error
	dst, err = db.Unwordle(dst, fpat, fneg)
	if err != nil {
		log.Fatalf("error '%s' caught on check\n", err.Error())
	}
	if len(dst) == 0 {
		log.Fatalln("empty set found")
	}
	fmt.Println("List of possible words:")
	for i := 0; i < len(dst); i++ {
		fmt.Println(dst[i])
	}
}
