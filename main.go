package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pranavmangal/termq/config"
	"github.com/pranavmangal/termq/providers/cerebras"
	"github.com/pranavmangal/termq/providers/groq"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/glamour"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Please provide a query to run.")
		return
	}

	s := spinner.New(spinner.CharSets[14], 20*time.Millisecond)
	s.Start()

	if !config.Exists() {
		configFilePath := config.Create()
		fmt.Println("No configuration file was found. It has been created for you at:")
		fmt.Println("  " + configFilePath)
		fmt.Println("Please fill atleast one provider and try again.")
		return
	}

	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("The configuration is invalid: %v\nPlease fix the errors and try again.", err)
	}

	query := args[1]
	var queryResp string

	if cfg.Cerebras.IsValid() {
		res, err := cerebras.RunQuery(query, cfg)
		if err != nil {
			fmt.Println("Failed to run query: ", err)
			return
		}
		queryResp = res

	} else if cfg.Groq.IsValid() {
		res, err := groq.RunQuery(query, cfg)
		if err != nil {
			fmt.Println("Failed to run query: ", err)
			return
		}
		queryResp = res

	} else {
		fmt.Println("No provider is configured. Please configure one and try again.")
		return
	}

	r, _ := glamour.NewTermRenderer(glamour.WithAutoStyle())
	mdOutput, err := r.Render(queryResp)
	if err != nil {
		s.Stop()
		fmt.Println(queryResp)
		os.Exit(0)
	}

	s.Stop()
	fmt.Println(mdOutput)
}
