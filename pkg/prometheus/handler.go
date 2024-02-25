package prometheus

import (
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
)

func GatewayHandler() gin.HandlerFunc {
	handler := promhttp.Handler()
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func RpcHandler(addr string) {
	port := strings.Split(addr, ":")[1]
	http.Handle("/metrics", promhttp.Handler())
	log.LogrusObj.Panic(http.ListenAndServe(":"+port, nil))
}
