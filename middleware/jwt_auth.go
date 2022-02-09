package middleware

import (
	"log"
	"net/http"

	"github.com/adetiamarhadi/golang-rest-api/helper"
	"github.com/adetiamarhadi/golang-rest-api/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := helper.BuildErrorResponse("failed to process request", "no token found", nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		token, err := jwtService.ValidateToken(authHeader)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claim[user_id]: ", claims["user_id"])
			log.Println("Clain[issuer]: ", claims["issuer"])
		} else {
			log.Println(err)
			response := helper.BuildErrorResponse("token is not valid", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
	}
}
