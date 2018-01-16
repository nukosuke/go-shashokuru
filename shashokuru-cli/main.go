package main

import (
	"github.com/nukosuke/go-shashokuru/shashokuru"
	"os"
)

func main() {
	homeDir := os.Getenv("HOME")
	configFile := homeDir + "/.shashokuru.toml"

	client := shashokuru.NewClient()

	//TODO:
	// - load config from configFile
	// - client.Login()

	//TODO: commands
	// - list [date]: client.Bento.GetListOnDate(date)
	// - reserve [bento_id]: client.Bento.Reserve(Bento)
}
