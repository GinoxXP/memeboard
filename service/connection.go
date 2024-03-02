package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/ini.v1"
)

var connection *pgxpool.Pool

type Config struct {
	Username string `ini:"username"`
	Password string `ini:"password"`
	Address  string `ini:"address"`
	Port     int    `ini:"port"`
	DbName   string `ini:"dbname"`
}

func Init() {
	inidata, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		panic(err)
	}

	var config Config
	err = inidata.MapTo(&config)

	bdUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		config.Username,
		config.Password,
		config.Address,
		config.Port,
		config.DbName)
	if err != nil {
		fmt.Printf("Fail to map config: %v", err)
		panic(err)
	}

	conn, err := pgxpool.New(context.Background(), bdUrl)
	if err != nil {
		fmt.Printf("Fail to connect to database: %v", err)
		panic(err)
	}

	connection = conn
}
