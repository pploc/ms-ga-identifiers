package messaging

import (
	"context"
	"encoding/json"

	"github.com/gym-api/ms-ga-identifier/pkg/config"
	"github.com/segmentio/kafka-go"
)

type EventType string

const (
	EventIdentityRegistered EventType = "identity.registered"
	EventIdentityLoggedIn   EventType = "identity.logged_in"
	EventIdentityLoggedOut  EventType = "identity.logged_out"
	EventPasswordChanged   EventType = "identity.password_changed"
)

type IdentityEvent struct {
	Type      EventType `json:"type"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp string    `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type KafkaProducer struct {
	writer *kafka.Writer
	topic  string
}

func NewKafkaProducer(cfg *config.KafkaConfig) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer: writer,
		topic:  cfg.Topic,
	}
}

func (p *KafkaProducer) PublishEvent(ctx context.Context, event IdentityEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.UserID),
		Value: data,
	})
}

func (p *KafkaProducer) PublishIdentityRegistered(ctx context.Context, userID, email string) error {
	event := IdentityEvent{
		Type:   EventIdentityRegistered,
		UserID: userID,
		Email:  email,
	}
	return p.PublishEvent(ctx, event)
}

func (p *KafkaProducer) PublishIdentityLoggedIn(ctx context.Context, userID, email string, metadata map[string]interface{}) error {
	event := IdentityEvent{
		Type:     EventIdentityLoggedIn,
		UserID:   userID,
		Email:    email,
		Metadata: metadata,
	}
	return p.PublishEvent(ctx, event)
}

func (p *KafkaProducer) PublishIdentityLoggedOut(ctx context.Context, userID, email string) error {
	event := IdentityEvent{
		Type:   EventIdentityLoggedOut,
		UserID: userID,
		Email:  email,
	}
	return p.PublishEvent(ctx, event)
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
