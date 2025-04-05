package main

import (
	"go.lvjp.me/demo-backend-go/cmd"

	_ "github.com/lib/pq"
)

func main() {
	cmd.Execute()
}
