package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type PackEvent struct {
	ID     uint64  `db:"id"`
	Type   string  `db:"type"`
	Entity []uint8 `db:"payload"`
}

type Pack struct {
	ID      uint64    `db:"id"`
	Name    string    `db:"name"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}

type ParsedPackEvent struct {
	ID     uint64 `db:"id"`
	Type   string `db:"type"`
	Entity *Pack  `db:"payload"`
}

type Consumer struct{}

// at start
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// at stop
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// consume
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var msg PackEvent
		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			return fmt.Errorf("parse PackEvent: %w", err)
		}
		result, err := parsePackEvent(&msg)
		if err != nil {
			return err
		}

		fmt.Println("MSG: ", *result)

		session.MarkMessage(message, "consumed")
	}

	return nil
}

// func for parsing PackEvent into ParsedPackEvent type
func parsePackEvent(packEvent *PackEvent) (*ParsedPackEvent, error) {
	var parsedEvent ParsedPackEvent
	var pack Pack
	err := json.Unmarshal(packEvent.Entity, &pack)
	if err != nil {
		return nil, fmt.Errorf("pack parse: %w", err)
	}
	parsedEvent.ID = packEvent.ID
	parsedEvent.Type = packEvent.Type
	parsedEvent.Entity = &pack
	return &parsedEvent, nil
}
