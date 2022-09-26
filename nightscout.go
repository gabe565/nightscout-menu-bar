package main

import (
	"encoding/json"
	"errors"
	"github.com/gabe565/nightscout-systray/internal/nightscout"
	"log"
	"net/http"
	"os"
	"time"
)

var url = os.Getenv("NIGHTSCOUT_URL")

func fetchFromNightscout() error {
	if url == "" {
		return errors.New("url is required")
	}

	// Fetch JSON
	resp, err := http.Get(url + "/api/v2/properties/bgnow,buckets,delta,direction")
	if err != nil {
		return err
	}

	// Decode JSON
	var properties nightscout.Properties
	if err := json.NewDecoder(resp.Body).Decode(&properties); err != nil {
		return err
	}

	updateTitle <- properties.String()
	updateHistory <- properties.Buckets
	updateLastReading <- properties.Bgnow.Time()
	return nil
}

func tick() {
	if err := fetchFromNightscout(); err != nil {
		log.Println(err)
	}

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := fetchFromNightscout(); err != nil {
					log.Println(err)
				}
			}
		}
	}()
}
