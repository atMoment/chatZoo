package main

import (
	"fmt"
)

func main() {
	app := NewApp()
	app.run()
	fmt.Println("main over")
}
