package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware handles panics and logs them
func RecoveryMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		if ne, ok := recovered.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
					strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
					logger.Warn("Connection error", "error", se.Error())
					c.Abort()
					return
				}
			}
		}

		httpRequest, _ := httputil.DumpRequest(c.Request, false)
		logger.Error("Panic recovered",
			"error", recovered,
			"request", string(httpRequest),
			"stack", string(debug.Stack()),
		)

		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
