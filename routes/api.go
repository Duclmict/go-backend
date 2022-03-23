package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/Duclmict/go-backend/app/http/controller/api/users"
    jwt "github.com/appleboy/gin-jwt/v2"
)

func (r routes) addv1Api(rg *gin.RouterGroup, jwt *jwt.GinJWTMiddleware) {

    // v1
    v1 := rg.Group("/v1")
    
    v1.POST("/signup", users.Store)
    v1.POST("/login", jwt.LoginHandler)
    v1.GET("/logout", jwt.LogoutHandler)

    // simple CRUD api
    v1.GET("/users", jwt.MiddlewareFunc() ,users.Index)
    v1.GET("/users/:id", jwt.MiddlewareFunc() ,users.Get)
    v1.GET("/users/search", jwt.MiddlewareFunc(), users.Search)
    v1.PUT("/users/:id", jwt.MiddlewareFunc() ,users.Update)
    v1.DELETE("/users/:id", jwt.MiddlewareFunc() ,users.Delete)
}

