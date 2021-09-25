package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	r := gin.Default()

	r.StaticFS("/", gin.Dir("./static", true))
	r.POST("/register", register)

	r.Run(":6464")
}

func addTask(c *gin.Context) {

}

func register(c *gin.Context) {
	// user := c.Request.FormValue("user")
	user := c.PostForm("user")
	pwd := c.PostForm("pwd")

	fmt.Println(user, pwd)

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"pwd":  pwd,
	})
	// c.String()
	// c.Redirect()
}
