# ComboTracker

irc bot responsible for recording chat combos in [#destinyecho](https://qchat.rizon.net/?channels=#destinyecho) channel and for storing them in a database. 

## Installation

* ### Get dependencies
```console
$ go get github.com/BurntSushi/toml
$ go get github.com/thoj/go-ircevent
$ go get github.com/go-sql-driver/mysql
```

* ### Compile packages and dependencies
```console
$ go build
```

* ### Update configuration file ```config.toml```

## Usage

* ### Run the executable

## TODO
- [ ] Improve comments
- [ ] Improve logging
- [ ] Handle errors where possible
- [ ] Save records to database in intervals
- [ ] If database is unavailable store records in a slice / file to dump later when database is back online
- [ ] Provide unit file for systemd
- [ ] Ability to update emotes without restarting the app
- [ ] Implement commands for controlling the bot

## Built With

* [TOML](https://github.com/BurntSushi/toml) - Parser and encoder for Go with reflection.
* [go-ircevent](https://github.com/thoj/go-ircevent) - Event based IRC client library in Go.
* [Go-MySQL-Driver](https://github.com/go-sql-driver/mysql) - A MySQL-Driver for Go's database/sql package.

## [License](LICENSE)