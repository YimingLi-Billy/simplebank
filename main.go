package main

import (
	"database/sql"
	"log"

	"github.com/YimingLi-Billy/simplebank/api"
	db "github.com/YimingLi-Billy/simplebank/db/sqlc"
	"github.com/YimingLi-Billy/simplebank/util"
	_ "github.com/lib/pq" // a GO postgres driver required to talk to the database
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
