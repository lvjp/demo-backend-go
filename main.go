package main

import (
	"go.lvjp.me/demo-backend-go/cmd"

	// Importing pq to register the Postgres driver
	_ "github.com/lib/pq"
)

func main() {
	cmd.Execute()
}
