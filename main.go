package main

import (
	"github.com/gin-gonic/gin"
	// _ "github.com/mattn/go-sqlite3"
)

func main() {
	engine := gin.Default()

	engine.StaticFS("/", gin.Dir("./static", true))
	engine.Run(":6464")
}
