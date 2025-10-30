package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authorization(key string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Status": "Failed",
				"Message": "Token not found in header",
				"Data":    "",
				"Error":   "null token"})
			ctx.Abort()
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Status": "Failed",
				"Message": "Token not valid",
				"Data":    "",
				"Error":   err.Error()})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Status": "Failed",
				"Message": "Invalid token claims",
				"Data":    "",
				"Error":   ok})
			ctx.Abort()
			return
		}
		email, ok := claims["Email"].(string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Status": "Failed",
				"Message": "Email not found in claims",
				"Data":    "",
				"Error":   ok})
			ctx.Abort()
			return
		}

		userID, ok := claims["UserID"].(string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Status": "Failed",
				"Message": "UserID not found in token",
				"Data":    userID,
				"Error":   ok})
			ctx.Abort()
			return
		}
		ctx.Set("email", email)
		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
