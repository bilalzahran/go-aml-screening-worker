package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/google/uuid"

	amqp "github.com/rabbitmq/amqp091-go"
)

// TestCase mirrors the shape of each entry in test_screening_data.json.
type TestCase struct {
	Comment        string          `json:"_comment"`
	ExpectedVerdict string         `json:"expected_verdict"`
	TestID         string          `json:"test_id"`
	Data           json.RawMessage `json:"data"`
}

// CaseData is the minimal extraction needed to build the OGS message.
type CaseData struct {
	Summary struct {
		CaseRecord struct {
			CaseID       string `json:"case_id"`
			Name         string `json:"name"`
			DateOfBirth  string `json:"date_of_birth"`
			Citizenship  string `json:"citizenship"`
			Gender       string `json:"gender"`
			EntityType   string `json:"entity_type"`
		} `json:"case_record"`
		WorldCheck []map[string]any `json:"world_check"`
	} `json:"summary"`
}

// Subject contains identifying information about the screened subject.
type Subject struct {
	Name       string `json:"name"`
	DOB        string `json:"dob"`
	Nationality string `json:"nationality"`
	Gender     string `json:"gender"`
	EntityType string `json:"entity_type"`
}

// OGSMessage is the RabbitMQ message conforming to ogs-schema.json.
type OGSMessage struct {
	ID               string            `json:"id"`
	EventType        string            `json:"event_type"`
	TenantID         *string           `json:"tenant_id"`
	CaseID           string            `json:"case_id"`
	Subject          Subject           `json:"subject"`
	WorldCheck       []map[string]any  `json:"world_check"`
	AIRecommendation []any             `json:"ai_recommendation"`
}

func main() {
	filePath := flag.String("file", "test_screening_data.json", "path to test data JSON file")
	amqpURL := flag.String("url", "amqp://guest:guest@localhost:5672/", "RabbitMQ connection URL")
	exchange := flag.String("exchange", "aml.screening", "exchange name")
	queue := flag.String("queue", "aml.screening.ogs", "queue name")
	routingKey := flag.String("routing-key", "screening.ogs", "routing key")
	tenantID := flag.String("tenant-id", "test-tenant", "tenant ID for published messages")
	flag.Parse()

	raw, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("failed to read test data file: %v", err)
	}

	var testCases []TestCase
	if err := json.Unmarshal(raw, &testCases); err != nil {
		log.Fatalf("failed to parse test data: %v", err)
	}

	log.Printf("loaded %d test cases from %s", len(testCases), *filePath)

	conn, err := amqp.Dial(*amqpURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(*exchange, "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
	}

	q, err := ch.QueueDeclare(*queue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	if err := ch.QueueBind(q.Name, *routingKey, *exchange, false, nil); err != nil {
		log.Fatalf("failed to bind queue: %v", err)
	}

	log.Printf("exchange=%s  queue=%s  routing_key=%s", *exchange, q.Name, *routingKey)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	published := 0
	for i, tc := range testCases {
		var caseData CaseData
		if err := json.Unmarshal(tc.Data, &caseData); err != nil {
			log.Printf("[%d] SKIP %s — failed to parse case data: %v", i+1, tc.TestID, err)
			continue
		}

		eventID := uuid.New().String()

		msg := OGSMessage{
			ID:        eventID,
			EventType: "wc.ogs",
			TenantID:  tenantID,
			CaseID:    caseData.Summary.CaseRecord.CaseID,
			Subject: Subject{
				Name:        caseData.Summary.CaseRecord.Name,
				DOB:         caseData.Summary.CaseRecord.DateOfBirth,
				Nationality: caseData.Summary.CaseRecord.Citizenship,
				Gender:      caseData.Summary.CaseRecord.Gender,
				EntityType:  caseData.Summary.CaseRecord.EntityType,
			},
			WorldCheck:       caseData.Summary.WorldCheck,
			AIRecommendation: []any{},
		}

		body, err := json.Marshal(msg)
		if err != nil {
			log.Printf("[%d] SKIP %s — failed to marshal OGS message: %v", i+1, tc.TestID, err)
			continue
		}

		err = ch.PublishWithContext(ctx,
			*exchange,
			*routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				MessageId:    eventID,
				Timestamp:    time.Now(),
				Headers: amqp.Table{
					"test_id":          tc.TestID,
					"expected_verdict": tc.ExpectedVerdict,
				},
				Body: body,
			},
		)
		if err != nil {
			log.Printf("[%d] FAIL %s — publish error: %v", i+1, tc.TestID, err)
			continue
		}

		hitCount := len(caseData.Summary.WorldCheck)
		log.Printf("[%d] OK   %s  case=%s  hits=%d  expected=%s",
			i+1, tc.TestID, caseData.Summary.CaseRecord.CaseID, hitCount, tc.ExpectedVerdict)
		published++
	}

	log.Printf("done — published %d/%d messages", published, len(testCases))
}
