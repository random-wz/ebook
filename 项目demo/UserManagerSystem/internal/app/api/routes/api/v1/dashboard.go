package v1

import "github.com/gin-gonic/gin"

func DashboardAPI(c *gin.Context) {
	_, err := c.Writer.Write([]byte("Hello World"))
	if err != nil {}
}
