package main

import (
	"flag"
	"fmt"
	"github.com/wavesplatform/gowaves/cmd/bcheck/internal"
	"os"
)

func main() {
	var (
		importFile      = flag.String("import-file", "", "Path to binary blockchain file to import.")
		db              = flag.String("db", "", "Path to data base.")
		transactionType = flag.Int("transaction-type", 0, "Filter transaction by type")
	)
	flag.Parse()

	if *db == "" {
		fmt.Printf("No data base path specified\n")
		os.Exit(1)
	}
	storage := internal.Storage{Path: *db}
	err := storage.Open()
	if err != nil {
		fmt.Printf("Failed to open storage: %s\n", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := storage.Close()
		if err != nil {
			fmt.Printf("Failed to close Storage: %s\n", err.Error())
		}
	}()

	importer := internal.NewImporter(&storage, *transactionType)
	importer.Import(*importFile)
}
