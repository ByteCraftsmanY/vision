package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/scylladb/gocqlx/v2"
	"github.com/twmb/franz-go/pkg/kgo"
	"net/http"
	"vision/db"
	"vision/dtos"
	"vision/kafka"
)

type HealthHandler struct {
	client  *kgo.Client
	session *gocqlx.Session
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		client:  kafka.GetClient(),
		session: db.GetSession(),
	}
}

func (c *HealthHandler) Status(ctx *gin.Context) {
	kafkaErr := c.client.Ping(ctx)
	dbStatus := c.session.Closed()
	health := dtos.HealthDTO{
		Status:      http.StatusOK,
		Message:     "Success",
		DBStatus:    !dbStatus,
		KafkaStatus: kafkaErr == nil,
	}
	ctx.AbortWithStatusJSON(http.StatusOK, health)
}
