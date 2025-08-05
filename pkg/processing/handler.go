package processing

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/censys/scan-takehome/pkg/database"
	"github.com/censys/scan-takehome/pkg/scanning"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

type Handler interface {
	SetSubscription(ctx context.Context, topicID string) error
	Receive(ctx context.Context) error
}

type handler struct {
	db     database.Handler
	client *pubsub.Client
	sub    *pubsub.Subscription
}

func NewHandler(db database.Handler, projectID string) *handler {
	client, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatal("Failed to create PubSub client:", err)
	}
	return &handler{
		db:     db,
		client: client,
	}
}

func (h *handler) processMessage(ctx context.Context, msg *pubsub.Message) {
	var scan scanning.Scan
	err := json.Unmarshal(msg.Data, &scan)
	if err != nil {
		log.Error("error unmarshaling message data:", err)
		return
	}
	err = h.db.ProcessMessage(ctx, scan)
	if err != nil {
		log.Error("error processing message:", err)
		return
	}
	msg.Ack() // Only ack if processing was successful
}

func (h *handler) SetSubscription(ctx context.Context, topicID *string) error {
	if topicID == nil {
		return fmt.Errorf("topicID cannot be nil")
	}
	topicString := *topicID
	subID := "sub_" + topicString + "_" + uuid.New().String()
	sub, err := h.client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		Topic: h.client.Topic(topicString),
	})
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	log.Info("Processor subscribed to topic:", topicString, "with subscription ID:", subID)
	h.sub = sub
	return nil
}

func (h *handler) Receive(ctx context.Context) error {
	if h.sub == nil {
		return fmt.Errorf("subscription is not set")
	}
	return h.sub.Receive(ctx, h.processMessage)
}
