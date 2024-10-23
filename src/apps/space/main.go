package main

import (
	"fmt"
)

func main() {
	app := NewApp()
	app.Run()
	fmt.Println("server main over")
}
