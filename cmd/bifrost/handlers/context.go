package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/subzerobo/ratatoskr/internal/services/authentication"
	"github.com/subzerobo/ratatoskr/pkg/utils"
	"strconv"
)

const (
	TokenKey     = "jwt_token"
	LocalDateKey = "local_date"
)

func getToken(c *gin.Context) *jwt.Token {
	obj, exists := c.Get(TokenKey)
	if !exists {
		return nil
	}
	return obj.(*jwt.Token)
}

func getClaims(c *gin.Context) *authentication.AppClaims {
	token := getToken(c)
	if token == nil {
		return nil
	}
	return token.Claims.(*authentication.AppClaims)
}

func getLocalDate(c *gin.Context) string {
	obj, exists := c.Get(LocalDateKey)
	if !exists {
		return ""
	}
	return obj.(string)
}

func getPagination(c *gin.Context) (utils.Paging, error) {
	pageNum := 1
	limitNum := 20
	var err error
	if pageStr := c.Query("page"); pageStr != "" {
		pageNum, err = strconv.Atoi(pageStr)
		if err != nil {
			return utils.Paging{}, err
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		limitNum, err = strconv.Atoi(limitStr)
		if err != nil {
			return utils.Paging{}, err
		}
	}
	return utils.Paging{
		Page: pageNum,
		Size: limitNum,
	}, nil
}
