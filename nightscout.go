package main

import (
	"encoding/json"
	"errors"
	"github.com/gabe565/nightscout-systray/internal/nightscout"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

func init() {
	flag.StringP("url", "u", "", "Nightscout base URL")
	if err := viper.BindPFlag("url", flag.Lookup("url")); err != nil {
		panic(err)
	}
}

func fetchFromNightscout() error {
	url := viper.GetString("url")
	if url == "" {
		return errors.New("url is required")
	}

	// Fetch JSON
	resp, err := http.Get(url + "/api/v2/properties/bgnow,buckets,delta,direction")
	if err != nil {
		errorChan <- err
		return err
	}

	// Decode JSON
	var properties nightscout.Properties
	if err := json.NewDecoder(resp.Body).Decode(&properties); err != nil {
		errorChan <- err
		return err
	}

	updateChan <- properties
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
