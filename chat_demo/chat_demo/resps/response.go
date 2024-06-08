package resps

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Resp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var (
	internal = Resp{
		Status:  50000,
		Message: "internal error ",
	}

	success = Resp{
		Status:  10000,
		Message: "success",
	}

	param = Resp{
		Status:  20000,
		Message: "param error ",
	}
)

func OK(c *gin.Context) {
	c.JSON(http.StatusOK, success)
}

func InternalErr(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, internal)
}

func ParamErr(c *gin.Context) {
	c.JSON(http.StatusBadRequest, param)
}

func OKWithData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"status":  10000,
		"message": "success",
		"data":    data,
	})
}
