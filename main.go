package main

import (
	"gii/demo/config"
)

func main() {
	config.Router().Run("localhost:8000")
}
