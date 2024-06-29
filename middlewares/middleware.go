package middlewares

import (
	"final-task-pbi-fullstackdev/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		//mengambil token value
		token, err := ctx.Cookie("jwt-token")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		claims := &app.JWTClaim{}
		tokenParse, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return app.JWT_KEY, nil
		})
		//cek apakah token malformed/expired/invalid
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token malformed/expired/invalid",
			})
			ctx.Abort()
			return
		}
		//cek apakah token valid atau tidak
		if !tokenParse.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		//jika token ada dan tidak menampilkan error, maka proses akan dialnjutkan ke controller handler
		ctx.Next()
	}
}
