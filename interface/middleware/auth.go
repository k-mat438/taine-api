package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func ClerkSessionAuth() gin.HandlerFunc {
	issuer := os.Getenv("CLERK_ISSUER")
	jwksURL := os.Getenv("CLERK_JWKS_URL")

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval:   time.Hour,
		RefreshTimeout:    10 * time.Second,
		RefreshUnknownKID: true,
	})
	if err != nil {
		panic("failed to load JWKS: " + err.Error())
	}

	return func(c *gin.Context) {
		ah := c.GetHeader("Authorization")
		if !strings.HasPrefix(ah, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer"})
			return
		}
		raw := strings.TrimPrefix(ah, "Bearer ")

		token, err := jwt.Parse(raw, jwks.Keyfunc)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token", "details": err.Error()})
			return
		}
		if !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "token is not valid"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid claims"})
			return
		}

		// issuer チェック
		if !claims.VerifyIssuer(issuer, true) {
			c.AbortWithStatusJSON(401, gin.H{"error": "bad issuer"})
			return
		}

		// 有効期限チェック
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			c.AbortWithStatusJSON(401, gin.H{"error": "expired"})
			return
		}

		if sub, _ := claims["sub"].(string); sub != "" {
			c.Set("sub_id", sub)
		}

		// JWT templateで追加したカスタムクレームを取得
		if orgID, _ := claims["org_id"].(string); orgID != "" {
			c.Set("org_external_id", orgID)
		}
		if orgRole, _ := claims["org_role"].(string); orgRole != "" {
			c.Set("org_role", orgRole)
		}

		c.Next()
	}
}
