package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/airbloc/flux-slack-alert/slack"
	"github.com/airbloc/logger"
	"github.com/airbloc/logger/module/loggergin"
	fluxevent "github.com/fluxcd/flux/pkg/event"
	"github.com/gin-gonic/gin"
)

const (
	// EndpointPath is path of the webhook receiver endpoint.
	EndpointPath = "/v1/event"
)

type Webhook struct {
	log *logger.Logger

	sender slack.Sender
	router *gin.Engine
	server *http.Server
}

func New(port int, sender slack.Sender) *Webhook {
	router := gin.New()
	router.Use(loggergin.Middleware("webhook"))
	w := &Webhook{
		log: logger.New("webhook"),

		sender: sender,
		router: router,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		},
	}
	router.Use(w.handleRecovery)
	router.NoRoute(w.handleNoRoute)

	router.POST(EndpointPath, w.handle)
	return w
}

func (w *Webhook) handle(c *gin.Context) {
	var msg struct {
		Event fluxevent.Event
	}
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slackMsg := w.sender.Compose(msg.Event)
	if err := w.sender.Send(slackMsg); err != nil {
		dumpedMsg, _ := json.Marshal(slackMsg)
		w.log.Error("Failed to send slack message: {}", err, dumpedMsg)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (w *Webhook) handleNoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": "not found",
	})
}

func (w *Webhook) handleRecovery(c *gin.Context) {
	defer func() {
		if r := w.log.Recover(logger.Attrs{
			"method": c.Request.Method,
			"url":    c.Request.URL.Path,
			"client": c.ClientIP(),
		}); r != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()
	c.Next()
}

func (w *Webhook) Start() error {
	go func() {
		if err := w.server.ListenAndServe(); err != http.ErrServerClosed {
			w.log.Debug(reflect.TypeOf(err).String())
			w.log.Error("failed to run webhook server: {}", err)
		}
	}()
	return nil
}

func (w *Webhook) Close() error {
	return w.server.Shutdown(context.Background())
}
