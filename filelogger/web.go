package filelogger

import "github.com/gin-gonic/gin"

func RegisterAPI(e *gin.Engine) {
	e.GET("/log/:biz", func(c *gin.Context) {
		logger := getLogger(c.Param("biz"))
		if logger == nil {
			c.String(400, "logger not found")
			return
		}
		c.String(200, logger.GetFromFile())
	})

	e.GET("/log/clean/:biz", func(c *gin.Context) {
		biz := c.Param("biz")
		if biz == "" {
			for _, logger := range loggerMap {
				logger.Clean()
			}
		} else {
			logger := getLogger(biz)
			if logger == nil {
				c.String(400, "logger not found")
				return
			}
			err := logger.Clean()
			if err != nil {
				c.JSON(500, err)
				return
			}
			c.String(200, "ok")
		}
	})
}
