package service_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"

	"cdaq-event-worker/internal/model"
	"cdaq-event-worker/internal/service"
)

type mockEventRepository struct {
	saveFunc func(ctx context.Context, event *model.Event) error
	saved    []*model.Event
}

func (m *mockEventRepository) Save(ctx context.Context, event *model.Event) error {
	m.saved = append(m.saved, event)
	if m.saveFunc != nil {
		return m.saveFunc(ctx, event)
	}
	return nil
}

func newTestService(repo *mockEventRepository) *service.EventService {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return service.NewEventService(repo, logger)
}

func TestProcess_ValidEvent(t *testing.T) {
	repo := &mockEventRepository{}
	svc := newTestService(repo)

	event := &model.Event{
		ID:   uuid.New().String(),
		Type: "user.created",
	}

	err := svc.Process(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repo.saved) != 1 {
		t.Fatalf("expected 1 saved event, got %d", len(repo.saved))
	}

	if repo.saved[0].CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}
}

func TestProcess_EmptyID(t *testing.T) {
	repo := &mockEventRepository{}
	svc := newTestService(repo)

	event := &model.Event{
		Type: "user.created",
	}

	err := svc.Process(context.Background(), event)
	if err == nil {
		t.Fatal("expected error for empty ID")
	}

	if len(repo.saved) != 0 {
		t.Fatal("expected no events saved")
	}
}

func TestProcess_EmptyType(t *testing.T) {
	repo := &mockEventRepository{}
	svc := newTestService(repo)

	event := &model.Event{
		ID: uuid.New().String(),
	}

	err := svc.Process(context.Background(), event)
	if err == nil {
		t.Fatal("expected error for empty type")
	}

	if len(repo.saved) != 0 {
		t.Fatal("expected no events saved")
	}
}

func TestProcess_RepoError(t *testing.T) {
	repoErr := errors.New("connection refused")
	repo := &mockEventRepository{
		saveFunc: func(_ context.Context, _ *model.Event) error {
			return repoErr
		},
	}
	svc := newTestService(repo)

	event := &model.Event{
		ID:   uuid.New().String(),
		Type: "user.created",
	}

	err := svc.Process(context.Background(), event)
	if err == nil {
		t.Fatal("expected error from repo")
	}

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repo error, got %v", err)
	}
}
