package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextAuthClaimsKey = "auth.claims"
	ContextAuthBankIDKey = "auth.bank_id"
)

func GetCurrentBank(c *gin.Context) (string, bool) {
	bankIDValue, bankIDExists := c.Get(ContextAuthBankIDKey)
	if !bankIDExists {
		return "", false
	}

	bankID, bankIDOk := bankIDValue.(string)
	if !bankIDOk || bankID == "" {
		return "", false
	}

	return bankID, true
}

func RequireAuth(jwt *JWT, h gin.HandlerFunc) gin.HandlerFunc {
	mw := CheckAuth(jwt)
	return func(c *gin.Context) {
		mw(c)
		if c.IsAborted() {
			return
		}
		h(c)
	}
}

func CheckAuth(jwtManager *JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		if jwtManager == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"err":    "UNKNOWN",
				"errmsg": "auth is not configured",
			})
			return
		}

		authorizationHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err":    "UNAUTHORIZED",
				"errmsg": "missing or invalid authorization header",
			})
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, "Bearer "))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err":    "UNAUTHORIZED",
				"errmsg": "missing bearer token",
			})
			return
		}

		claims, err := jwtManager.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err":    "UNAUTHORIZED",
				"errmsg": "invalid token",
			})
			return
		}

		authClaims, claimsOk := ClaimsFromToken(claims)
		if !claimsOk {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err":    "UNAUTHORIZED",
				"errmsg": "invalid token claims",
			})
			return
		}

		bankID := authClaims.BankID

		c.Set(ContextAuthClaimsKey, claims)
		c.Set(ContextAuthBankIDKey, bankID)
		c.Next()
	}
}
