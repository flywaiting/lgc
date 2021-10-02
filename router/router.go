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
	r.POST("/add", addTask)

	r.Run(":6464")
}

func addTask(c *gin.Context) {
	// c.SetCookie()
	// http.StatusUnauthorized
}

func login(c *gin.Context) {

}

func register(c *gin.Context) {
	user := c.PostForm("user")
	pwd := c.PostForm("pwd")

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"pwd":  pwd,
	})

	// c.Request.ParseMultipartForm(1024)
	fmt.Println(c.Request.Form)
	// c.String(http.StatusOK, "done")
	// c.Redirect(http.StatusOK, "/")
	// c.Request.Method = http.MethodGet
	// c.Request.URL.Path = "/"
	c.Redirect(http.StatusFound, "/")
}
