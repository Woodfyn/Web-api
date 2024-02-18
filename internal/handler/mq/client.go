package mq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Woodfyn/Web-api/internal/config"
	"github.com/Woodfyn/auditLog/pkg/core"
	"github.com/streadway/amqp"
)

type Server struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func New() *Server {
	return &Server{
		conn: new(amqp.Connection),
		ch:   new(amqp.Channel),
	}
}

func (s *Server) Init(cfg *config.Config) error {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.MQ.Username, cfg.MQ.Password, cfg.MQ.Host, cfg.MQ.Port)

	conn, err := amqp.Dial(addr)
	if err != nil {
		return err
	}
	s.conn = conn

	ch, err := s.conn.Channel()
	if err != nil {
		return err
	}
	s.ch = ch

	return nil
}

func (s *Server) Close() error {
	if err := s.ch.Close(); err != nil {
		return err
	}

	if err := s.conn.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Publisher(item core.LogItem) error {
	t, err := json.Marshal(core.LogItem{
		Entity:    item.Entity,
		Action:    item.Action,
		EntityID:  item.EntityID,
		Timestamp: item.Timestamp,
	})
	if err != nil {
		return err
	}

	q, err := s.ch.QueueDeclare(
		"log", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	if err := s.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        t,
		}); err != nil {
		log.Fatal("failed to declare a queue")
	}

	return nil
}
