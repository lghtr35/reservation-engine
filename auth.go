package main

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lghtr35/reservation-engine/models"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func keyFunc(secretKey string, logger *zerolog.Logger) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Error().Err(http.ErrAbortHandler).Msg("jwtAuthMiddleware: had errors when get key while trying to Parse token")
			return nil, http.ErrAbortHandler
		}
		return []byte(secretKey), nil
	}
}

func jwtAuthMiddleware(configuration *models.Configuration, db *gorm.DB, logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if configuration == nil {
			logger.Error().Msg("jwtAuthMiddleware: had an error when parsing jwt token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Token could not be parsed"})
			c.Abort()
			return
		}
		tokenString := c.GetHeader("Authorization")

		token, err := jwt.Parse(tokenString, keyFunc(configuration.Secret, logger))
		if err != nil {
			logger.Error().Err(err).Msg("jwtAuthMiddleware: had an error when parsing jwt token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Token could not be parsed"})
			c.Abort()
			return
		}
		if !token.Valid {
			logger.Debug().Msg(fmt.Sprintf("jwtAuthMiddleware: token is not valid: %v", tokenString))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Token is not valid"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var customer models.Customer
			res := db.First(&customer, claims["customerId"])
			if res.Error != nil {
				logger.Debug().Msg(fmt.Sprintf("jwtAuthMiddleware: claims are not valid: %v", tokenString))
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No customer with given id"})
				c.Abort()
				return
			}
			c.Set("claims", claims)
		} else {
			logger.Debug().Msg(fmt.Sprintf("jwtAuthMiddleware: claims are not valid: %v", tokenString))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Claims are not valid"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func apiKeyAuthMiddleware(db *gorm.DB, logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiSecret := c.GetHeader("x-api-secret")
		apiToken := c.GetHeader("x-api-token")

		var secret models.Secret
		res := db.Where("value = ?", apiSecret).First(&secret)
		if res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				logger.Debug().Msg(fmt.Sprintf("apiKeyAuthMiddleware: secret is not valid: %v", apiSecret))
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Cant match given secret"})
				c.Abort()
				return
			}
			logger.Err(res.Error).Msg(fmt.Sprintf("apiKeyAuthMiddleware: an error occured: %s", res.Error.Error()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - An error occured"})
			c.Abort()
			return
		}

		var token models.ApiToken
		res = db.Where("customerId = ? AND token = ?", secret.CustomerID, apiToken).First(&token)
		if res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				logger.Debug().Msg(fmt.Sprintf("apiKeyAuthMiddleware: token is not valid: %v", apiToken))
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Cant match given api token"})
				c.Abort()
				return
			}
			logger.Err(res.Error).Msg(fmt.Sprintf("apiKeyAuthMiddleware: an error occured: %s", res.Error.Error()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - An error occured"})
			c.Abort()
			return
		}

		c.Next()
	}
}
