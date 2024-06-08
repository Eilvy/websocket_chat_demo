package routers

import (
	"github.com/gin-gonic/gin"
	"go_code/chat_demo/chat_demo/Wschat"
)

func InitRouters() {
	r := gin.Default()

	r.GET("/chat/:id", Wschat.WsChat)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
