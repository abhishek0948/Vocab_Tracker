package routes

import (
	"github.com/abhishek0948/Vocab_Tracker/controllers"
	"github.com/abhishek0948/Vocab_Tracker/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}
}

func VocabRoutes(router *gin.Engine) {
	vocab := router.Group("/vocab")
	vocab.Use(middleware.AuthMiddleware())
	{
		vocab.GET("", controllers.GetVocabulary)
		vocab.POST("", controllers.CreateVocabulary)
		vocab.PUT("/:id", controllers.UpdateVocabulary)
		vocab.DELETE("/:id", controllers.DeleteVocabulary)
	}
}