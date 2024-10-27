package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// 除非連 Writer 都要覆寫
//func CustomRecovery() gin.HandlerFunc {
//	// 獲取原始的 gin.Recovery 中間件
//	return gin.CustomRecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, err interface{}) {
//		// 在這裡添加您的自定義日誌記錄邏輯
//		log.Error().Msgf("Custom recovery log: Panic recovered: %v", err)
//		// 繼續使用 gin 的內建恢復功能
//		c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
//	})
//}

func CustomRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		// 在這裡添加您的自定義日誌記錄邏輯
		log.Error().Msgf("Custom recovery log: Panic recovered: %v", err)
		// 繼續使用 gin 的內建恢復功能
		c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
	})
}
