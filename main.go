package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pranavmangal/termq/common"
	"github.com/pranavmangal/termq/config"
	"github.com/pranavmangal/termq/providers/cerebras"
	"github.com/pranavmangal/termq/providers/gemini"
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

	if !common.ModelCacheExists() {
		common.CreateModelCache()
	} else {
		go func() {
			common.UpdateModelCache()
		}()
	}

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

	// Prefer fastest provider first
	if cfg.Cerebras.Exists() {
		queryResp, err = cerebras.RunQuery(query, cfg)
	} else if cfg.Groq.Exists() {
		queryResp, err = groq.RunQuery(query, cfg)
	} else if cfg.Gemini.Exists() {
		queryResp, err = gemini.RunQuery(query, cfg)
	} else {
		s.Stop()
		fmt.Println("No provider is configured. Please configure one and try again.")
		return
	}

	if err != nil {
		s.Stop()
		log.Fatalf("Error while running query: %v\n", err)
	}

	r, _ := glamour.NewTermRenderer(glamour.WithAutoStyle())
	mdOutput, err := r.Render(queryResp)
	s.Stop()
	if err != nil {
		fmt.Println(queryResp)
	} else {
		fmt.Println(mdOutput)
	}
}
