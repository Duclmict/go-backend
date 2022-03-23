package routes

import (
    jwt "github.com/appleboy/gin-jwt/v2"
    "github.com/gin-gonic/gin"
    "github.com/Duclmict/go-backend/app/http/middleware"
)

type routes struct {
    router *gin.Engine
    jwt *jwt.GinJWTMiddleware
}

func InitRouter(){

    r := routes{
        router: gin.Default(),
        jwt: middleware.JwtInit(),
    }

    // add cors setting
    r.router.Use(middleware.Cors())

    // // Jwt Middleware setting
    api := r.router.Group("/api/")

    r.addv1Api(api, r.jwt)
	r.router.Run("127.0.0.1:7777")

	return 
}