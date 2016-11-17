package main

import (
    "log"
    "time"
    "regexp"
    "strings"
    "crypto/tls"
    "database/sql"
    mysql "github.com/go-sql-driver/mysql"
    irc "github.com/thoj/go-ircevent"
)

type comboRecord struct {
    Emote string
    Memers []string
    Timestamp string
}

var cfg Config

func main() {
    // Load our config file
    cfg = readConfig("config.toml")
    // Setup mysql config
    dsn := mysql.Config{User: cfg.DB.Username, Passwd: cfg.DB.Password, Net: cfg.DB.Net, Addr: cfg.DB.Addr, DBName: cfg.DB.Name}
    // Connect to the database
    db, err := sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    // Prepare statement for inserting our combo records
    stmtIns, err := db.Prepare("INSERT INTO chat_combos VALUES(?, ?, ?, ?, ?)")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

    emotes := fetchEmotes(cfg.EmotesURL)
    combo := comboRecord{}
    prevEmote := ""

    re := regexp.MustCompile(`\<([a-zA-Z0-9_\s]{3,20})\>\s(.*)`)

    irccon := irc.IRC(cfg.IRC.Nick, cfg.IRC.User)
    irccon.Debug = cfg.Debug
    if cfg.IRC.Secure {
        irccon.UseTLS = true
        irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
    }
    irccon.Password = cfg.IRC.Password
    irccon.AddCallback("001", func(event *irc.Event) {irccon.Join(cfg.IRC.Channel) })
    irccon.AddCallback("366", func(event *irc.Event) {})
    irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
        if event.Nick == "II" {
            msg := re.FindStringSubmatch(event.Message())

            // If msg[2] is an emote
            if index(emotes, msg[2]) != -1 {
                if prevEmote == msg[2] {
                    combo.Timestamp = time.Now().Format(time.RFC3339)
                    combo.Memers = append(combo.Memers, msg[1])
                } else {
                    if len(combo.Memers) >= cfg.MinCombo {
                        saveRecord2DB(stmtIns, combo)
                    }
                    combo.Emote = msg[2]
                    combo.Memers = []string{msg[1]}
                    prevEmote = msg[2]
                }
            } else {
                if len(combo.Memers) >= cfg.MinCombo {
                    saveRecord2DB(stmtIns, combo)
                    combo.Memers = []string{}
                }
                prevEmote =  ""
            }
        }
    });
    err = irccon.Connect(cfg.IRC.Server)
	if err != nil {
		log.Printf("Err %s", err )
		return
	}
    irccon.Loop()
}

func saveRecord2DB(stmtIns *sql.Stmt, combo comboRecord) {
    log.SetPrefix("[Database] ")
    // Insert combo record into "chat_combos" table
    _, err := stmtIns.Exec(nil, combo.Emote, len(combo.Memers), strings.Join(combo.Memers, ","), combo.Timestamp)
    if err != nil {
        log.Panic(err.Error())
        panic(err.Error()) // proper error handling instead of panic in your app
    } else if cfg.Debug {
        log.Print("Combo record has been saved into \"chat_combos\" table.")
    }
}