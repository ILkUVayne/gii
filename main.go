package main

import (
	"gii/demo/config"
)

type User struct {
	Name string
	Age  int
}

func main() {
	config.Router().Run("localhost:8000")
}
