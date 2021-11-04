package main

import (
	"github.com/ParkChangSun/nomadcoin/cli"
	"github.com/ParkChangSun/nomadcoin/db"
)

func main() {
	defer db.Close()
	cli.Start()

}
