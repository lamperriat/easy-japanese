package subtitles

import (
	"backend/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type TokenizeReq struct {
	Text string `json:"text" binding:"required"`
}

type TokenizeRsp struct {
	Tokens []string `json:"tokens"`
}

// @Summary Tokenize a sentence using Kagome
// @Description 
// @Tags parser
// @Security JWTAuth
// @Accept json
// @Produce json
// @Param text body TokenizeReq true "Text to tokenize"
// @Success 200 {object} TokenizeRsp
// @Failure 400 {object} models.ErrorMsg "Invalid JSON format"
// @Router /api/subtitles/tokenize [POST]
func TokenizeSentence(t *tokenizer.Tokenizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TokenizeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, models.ErrorMsg{Error: "Invalid JSON format"})
			return
		}

		tokens := t.Tokenize(req.Text)
		var result []string
		for _, token := range tokens {
			if token.Surface != "" {
				result = append(result, token.Surface)
			}
		}

		c.JSON(200, TokenizeRsp{
			Tokens: result,
		})
		
	}
}