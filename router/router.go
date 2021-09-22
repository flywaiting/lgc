package router

import "github.com/gin-gonic/gin"

func init() {
	r := gin.Default()

	r.StaticFS("/", gin.Dir("./static", true))
	r.Run(":6464")
}

func addTask(c *gin.Context) {

}
