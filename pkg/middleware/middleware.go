package middleware

import (
	"fp_pinjaman_online/model/dto/json"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	applicationName  = os.Getenv("APP_NAME")
	jwtSigningMethod = jwt.SigningMethodHS256
	jwtSignatureKey  = []byte(os.Getenv("JWT_SIGNATURE_KEY"))
)

func GenerateTokenJwt(userId, email, roleName, status string, expiredAt int64) (string, error) {
	loginExpDuration := time.Duration(expiredAt) * time.Hour
	myExpiresAt := time.Now().Add(loginExpDuration).Unix()

	claims := json.JwtClaim{
		StandardClaims: jwt.StandardClaims{
			Issuer:    applicationName,
			ExpiresAt: myExpiresAt,
		},
		UserId: userId,
		Email:  email,
		Roles: roleName,
		Status: status,
	}

	token := jwt.NewWithClaims(
		jwtSigningMethod,
		claims,
	)

	signedToken, err := token.SignedString(jwtSignatureKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			json.NewResponseUnauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)
		claims := &json.JwtClaim{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSignatureKey, nil
		})
		if err != nil {
			json.NewResponseUnauthorized(c, "Invalid token")
			c.Abort()
			return
		}
		if !token.Valid {
			json.NewResponseUnauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		// set user's id in the context
		c.Set("roleName", claims.Roles)
		c.Set("userId", claims.UserId)
		c.Set("status", claims.Status)
		c.Next()
	}
}

func JWTAuthWithRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			json.NewResponseUnauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)
		claims := &json.JwtClaim{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSignatureKey, nil
		})

		if err != nil {
			json.NewResponseUnauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		if !token.Valid {
			json.NewResponseForbidden(c, "Forbidden")
			c.Abort()
			return
		}

		// validation role
		validRole := false
		if len(roles) > 0 {
			for _, role := range roles {
				if role == claims.Roles {
					validRole = true
					break
				}
			}
		}
		if !validRole {
			json.NewResponseForbidden(c, "Forbidden")
			c.Abort()
			return
		}
		c.Set("roleName", claims.Roles)
		c.Set("userId", claims.UserId)
		c.Set("status", claims.Status)
		c.Next()
	}
}

func VerifiedOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		status, _ := ctx.Get("status")
		if status.(string) != "verified" {
			json.NewAbortForbidden(ctx, "user not verified")
			return
		}
		ctx.Next()
	}
}
