package middleware

import (
	"fmt"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
    "github.com/gin-gonic/gin"
	"github.com/Duclmict/go-backend/app/model"
	"github.com/Duclmict/go-backend/app/service/log_service"
)

type login struct {
	Email string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func JwtInit() (res *jwt.GinJWTMiddleware) {

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "core_server",
		Key:         []byte("secret key 123"),
		Timeout:     time.Hour * 2,
		MaxRefresh:  time.Hour * 2,
		IdentityKey: model.IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
		  if convert,ok := data.(model.Users); ok {
			return jwt.MapClaims{
				model.IdentityKey: convert.Email,
			}
		  }
		  return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
		  claims := jwt.ExtractClaims(c)
		  return &model.Users{
			Email: claims[model.IdentityKey].(string),
		  }
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
		  var loginVals login
		  if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		  }

		  userID := loginVals.Email
		  password := loginVals.Password

		  users, err_v := model.VerifyLogin(userID, password)

		  if (err_v != nil) {
			return nil, jwt.ErrFailedAuthentication
		  }
		  	
		  return users, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
		  if _, ok := data.(*model.Users); ok {
			return true
		  }
	
		  return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
		  c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		  })
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
	
		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",
	
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log_service.Error("MIDDLEWARE:[JWT] " + fmt.Sprint(err))
		return nil
	}

	return authMiddleware
}