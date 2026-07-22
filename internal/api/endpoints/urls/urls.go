package urls

import (
	"gateway/internal/auth"
	"gateway/internal/memory"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAPIURLsV1() gin.HandlerFunc {
	return func(c *gin.Context) {
		bankID := c.Param("bankId")

		memory.Mutex.Lock()

		data, ok := memory.APIURLss[bankID]

		memory.Mutex.Unlock()

		if !ok {
			c.Status(404)
			return
		}

		c.JSON(200, data)
	}
}

func SetAPIURLsV1() gin.HandlerFunc {
	return func(c *gin.Context) {
		bankID := c.Param("bankId")

		claimsBankID, ok := auth.GetCurrentBank(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"err":    "UNAUTHORIZED",
				"errmsg": "unauthorized",
			})
			return
		}

		if bankID != claimsBankID {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		var body memory.APIURLs

		if c.BindJSON(&body) != nil {
			c.Status(400)
			return
		}

		body.UpdatedAt = time.Now()

		memory.Mutex.Lock()

		memory.APIURLss[bankID] = body

		memory.Mutex.Unlock()

		c.Status(201)
	}
}
