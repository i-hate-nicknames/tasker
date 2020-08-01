package main

import (
	"fmt"

	"github.com/i-hate-nicknames/tasker/db"
)

func main() {
	fmt.Println("tasker")
	// web.StartServer()
	db.Run()
}
