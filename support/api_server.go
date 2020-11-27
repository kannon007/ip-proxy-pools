package support

import "github.com/gin-gonic/gin"

func ApiServerStart() {

	router := gin.Default()

	router.GET("/api/list", ipList)
	router.Run(":8000")
}

func ipList(c *gin.Context) {
	c.JSON(200, GetIPs())
}
