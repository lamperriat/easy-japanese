package proxy

import (
	"backend/pkg/models"
	"io"
	"net/http"

	// "os"

	// "github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type ProxyReq struct {
	URL string `json:"url" binding:"required"`
	Headers map[string]string `json:"headers,omitempty"`
	Body string `json:"body,omitempty"`
	Selector string `json:"selector,omitempty"`
}


// @Summary Proxy to bypass CORS issues
// @Description 
// @Tags proxy
// @Security JWTAuth
// @Accept json
// @Produce json
// @Failure 404 {object} models.ErrorMsg "User grammar not found"
// @Failure 400 {object} models.ErrorMsg "Invalid JSON format"
// @Failure 500 {object} models.ErrorMsg "Other error"
// @Router /api/proxy [POST]
func CORSProxy() func(c *gin.Context) {
	// selector: deprecated because we have to execute js on server in that case
	return func(c *gin.Context) {
		var req ProxyReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, models.ErrorMsg{Error: "Invalid JSON format"})
			return
		}

		rsp, err := http.Get(req.URL)
		if err != nil || rsp.StatusCode != http.StatusOK {
			c.JSON(500, models.ErrorMsg{Error: "Failed to fetch the URL"})
			return
		}
		defer rsp.Body.Close()

		for key, values := range rsp.Header {
			for _, value := range values {
				c.Writer.Header().Add(key, value)
			}
		}

		io.Copy(c.Writer, rsp.Body)
	}	
}