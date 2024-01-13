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

type HealthController struct {
	client  *kgo.Client
	session *gocqlx.Session
}

func NewHealthController() *HealthController {
	client := kafka.GetClient()
	session := db.GetSession()
	return &HealthController{
		client:  client,
		session: session,
	}
}

func (c *HealthController) Status(ctx *gin.Context) {
	//kafkaErr := c.client.Ping(ctx)
	dbStatus := c.session.Closed()
	health := dtos.HealthDTO{
		Status:   http.StatusOK,
		DBStatus: !dbStatus,
		//KafkaStatus: kafkaErr != nil,
	}
	ctx.AbortWithStatusJSON(http.StatusOK, health)
}
