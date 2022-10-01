package main

import (
	"GoVoteApi/cmd"
	"GoVoteApi/config"

	_ "github.com/joho/godotenv/autoload"
)

var cfg = &config.Config{}

func init() {
	if err := config.ParseEnv(cfg); err != nil {
		panic(err)
	}
}

func main() {
	if err := cmd.Run(cfg); err != nil {
		panic(err)
	}
}
