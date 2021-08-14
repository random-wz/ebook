package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ums/internal/app/api/controller/user"
)

func GetUserListAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": user.GetUserList(),
	})
}
