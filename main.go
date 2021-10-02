package main

import (

	// _ "github.com/mattn/go-sqlite3"

	"fmt"
	"lgc/cfg"
	// _ "lgc/router"
)

func main() {
	// engine := gin.Default()

	// engine.StaticFS("/", gin.Dir("./static", true))
	// engine.Run(":6464")

	fmt.Printf("%v\n", cfg.Cfg())
}
