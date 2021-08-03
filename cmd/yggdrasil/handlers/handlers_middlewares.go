package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/juliangruber/go-intersect"
	authentication2 "github.com/subzerobo/ratatoskr/internal/services/authentication"
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"net/http"
	"regexp"
)

var bearerRegexp = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

// CORSMiddleware provides CORS ability for the router
func (h *YggdrasilHandler) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH, HEAD")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

func (h *YggdrasilHandler) JWTMiddleware(secret string, allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractBearerToken(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token, err := parseJWTClaims(tokenString, secret)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if allowedRoles != nil {
			appClaims := token.Claims.(*authentication2.AppClaims)
			interSection := intersect.Hash(appClaims.Roles, allowedRoles)
			if len(interSection.([]interface{})) == 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
		c.Set(TokenKey, token)
		c.Next()
	}
}

// extractBearerToken get the jwt token from header
func extractBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("this endpoint requires a Bearer token")
	}
	
	matches := bearerRegexp.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", errors.New("this endpoint requires a Bearer token")
	}
	return matches[1], nil
}

// parseJWTClaims tries to decode the token using provided secret
func parseJWTClaims(bearer string, secret string) (*jwt.Token, error) {
	p := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}
	token, err := p.ParseWithClaims(bearer, &authentication2.AppClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
	
}
