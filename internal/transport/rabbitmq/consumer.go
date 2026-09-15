package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"

	"cdaq-event-worker/internal/config"
	"cdaq-event-worker/internal/handler"
)

type Consumer struct {
	cfg     config.RabbitMQConfig
	handler *handler.EventHandler
	logger  *slog.Logger
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan struct{}
}

func NewConsumer(cfg config.RabbitMQConfig, h *handler.EventHandler, logger *slog.Logger) *Consumer {
	return &Consumer{
		cfg:     cfg,
		handler: h,
		logger:  logger.With("component", "rabbitmq_consumer"),
		done:    make(chan struct{}),
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	var err error

	c.conn, err = amqp.Dial(c.cfg.URI)
	if err != nil {
		return fmt.Errorf("connecting to RabbitMQ: %w", err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("opening channel: %w", err)
	}

	if err := c.channel.Qos(c.cfg.PrefetchCount, 0, false); err != nil {
		return fmt.Errorf("setting QoS: %w", err)
	}

	if err := c.channel.ExchangeDeclare(
		c.cfg.Exchange,
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		return fmt.Errorf("declaring exchange: %w", err)
	}

	queue, err := c.channel.QueueDeclare(
		c.cfg.Queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("declaring queue: %w", err)
	}

	if err := c.channel.QueueBind(
		queue.Name,
		c.cfg.RoutingKey,
		c.cfg.Exchange,
		false, // no-wait
		nil,
	); err != nil {
		return fmt.Errorf("binding queue: %w", err)
	}

	deliveries, err := c.channel.Consume(
		queue.Name,
		"cdaq-event-worker", // consumer tag
		false,               // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("starting consumer: %w", err)
	}

	c.logger.Info("consumer started",
		"queue", queue.Name,
		"exchange", c.cfg.Exchange,
		"routing_key", c.cfg.RoutingKey,
	)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("consumer context cancelled, stopping")
			close(c.done)
			return nil
		case delivery, ok := <-deliveries:
			if !ok {
				c.logger.Warn("delivery channel closed")
				close(c.done)
				return nil
			}
			c.processDelivery(ctx, delivery)
		}
	}
}

func (c *Consumer) processDelivery(ctx context.Context, delivery amqp.Delivery) {
	if err := c.handler.Handle(ctx, delivery.Body); err != nil {
		c.logger.Error("failed to handle delivery",
			"error", err,
			"delivery_tag", delivery.DeliveryTag,
		)
		_ = delivery.Nack(false, true)
		return
	}

	if err := delivery.Ack(false); err != nil {
		c.logger.Error("failed to ack delivery",
			"error", err,
			"delivery_tag", delivery.DeliveryTag,
		)
	}
}

func (c *Consumer) Shutdown(ctx context.Context) {
	c.logger.Info("shutting down consumer")

	if c.channel != nil {
		if err := c.channel.Cancel("cdaq-event-worker", false); err != nil {
			c.logger.Error("error cancelling consumer", "error", err)
		}
	}

	select {
	case <-c.done:
		c.logger.Info("consumer stopped gracefully")
	case <-ctx.Done():
		c.logger.Warn("consumer shutdown timed out")
	}

	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
