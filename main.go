package main

import (

	// _ "github.com/mattn/go-sqlite3"

	"fmt"
	"os/exec"
	"strings"
	// _ "lgc/router"
)

func main() {
	// engine := gin.Default()

	// engine.StaticFS("/", gin.Dir("./static", true))
	// engine.Run(":6464")

	// fmt.Printf("%v\n", cfg.Cfg())
	demo()
}

func demo() {
	cmd := exec.Command("bash", "-c", "git branch -r")

	// out := bytes.Buffer{}
	// cmd.Stdout = os.Stdout
	// cmd.Run()
	// fmt.Println(out)

	arr, _ := cmd.Output()
	// strs := strings.Split(string(arr), "\n")
	for _, v := range strings.Split(string(arr), "\n") {
		if len(v) > 0 && !strings.Contains(v, "HEAD") {
			fmt.Println(strings.TrimPrefix(strings.TrimSpace(v), "origin/"))
		}
	}
	// strings.Fields()

	// fmt.Printf("%v", strs)
}
