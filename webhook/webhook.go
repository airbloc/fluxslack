package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/airbloc/flux-slack-alert/slack"
	"github.com/airbloc/logger"
	"github.com/airbloc/logger/module/loggergin"
	fluxevent "github.com/fluxcd/flux/pkg/event"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
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
	router := gin.Default()
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
	router.POST(EndpointPath, w.handle)
	router.NoRoute(w.handleNoRoute)
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
	c.Status(http.StatusCreated)
}

func (w *Webhook) handleNoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": "not found",
	})
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
