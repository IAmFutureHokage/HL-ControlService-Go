package main

import (
	"github.com/IAmFutureHokage/HL-ControlService-Go/database"
)

func main() {
	_, err := database.OpenDB()
	if err != nil {
		panic(err)
	}
}
