package main

import (
    "os"
    "log"
    "net/http"
    "encoding/json"
	"github.com/BurntSushi/toml"
)

// Config is a struct for storing our toml config values
type Config struct {
	Debug    bool
    EmotesURL string
	MinCombo int
	DB databaseInfo `toml:"database"`
    IRC ircInfo
}

type databaseInfo struct {
	Username     string
	Password string
	Net      string
	Addr     string
	Name   string
}

type ircInfo struct {
	Secure  bool
	Server  string
	Channel string
	Nick    string
    Password string
	User    string
}

func readConfig(configFile string) Config {
    log.SetPrefix("[Config] ")
    _, err := os.Stat(configFile)
	if err != nil {
		log.Fatalf("Config file \"%s\" not found!", configFile)
	}

    var config Config
    if _, err = toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}
    return config
}

func fetchEmotes(url string) []string {
    log.SetPrefix("[FetchEmotes] ")
    resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
    // Decode json array
    emotes := []string{}
    err = dec.Decode(&emotes)
    if err != nil {
        log.Fatal(err)
    }
    return emotes
}

func index(vs []string, t string) int {
    for i, v := range vs {
        if v == t {
            return i
        }
    }
    return -1
}