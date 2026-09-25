package usecase

import (
	"cdaq-event-worker/internal/client/typesafe"
	"cdaq-event-worker/internal/model"
	"cdaq-event-worker/internal/service"
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

const MAX_CONCURRENCY = 50

const SKIP_CLASSIFICATION = true

type OgsUseCase struct {
	logger                   *slog.Logger
	jobService               *service.JobService
	screeningWcResultService *service.ScreeningWcResultService
	typesafeClient           *typesafe.Client
}

func NewOgsUseCase(logger *slog.Logger, jobService *service.JobService, screeningWcResultService *service.ScreeningWcResultService, typesafeClient *typesafe.Client) *OgsUseCase {
	return &OgsUseCase{
		logger:                   logger,
		jobService:               jobService,
		screeningWcResultService: screeningWcResultService,
		typesafeClient:           typesafeClient,
	}
}

func (o *OgsUseCase) Handle(ctx context.Context, event *model.Event) error {
	results := make([]model.ScreeningWcResult, len(event.WorldCheck))
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(MAX_CONCURRENCY)

	for i, child := range event.WorldCheck {
		childEvent, err := model.MapToWorldCheckHits(child)
		if err != nil {
			return err
		}

		g.Go(func() error {
			res, _ := o.processOneChild(ctx, childEvent, &event.Subject)
			results[i] = *res
			return nil
		})
	}

	// TODO: Final save to screening wc + screening wc ogs

	return nil
}

func (o *OgsUseCase) processOneChild(ctx context.Context, hit *model.WorldCheckHits, subject *model.OgsEventSubject) (*model.ScreeningWcResult, error) {
	// Debug: log input data
	o.logger.Debug("processing child hit",
		"result_id", hit.ResultID,
		"reference_id", hit.ReferenceID,
		"matched_term", hit.MatchedTerm,
		"primary_name", hit.PrimaryName,
		"provider_type", hit.ProviderType,
	)

	// Build ScreeningWcResult from WorldCheckHits
	screeningWcResult := model.NewScreeningWcResultFromWorldCheckHits(hit)

	// Debug: log output data
	o.logger.Debug("converted to screening wc result",
		"result_id", screeningWcResult.ResultID,
		"reference_id", screeningWcResult.ReferenceID,
		"matched_term", screeningWcResult.MatchedTerm,
		"primary_name", screeningWcResult.PrimaryName,
		"source", screeningWcResult.Source,
	)

	if !SKIP_CLASSIFICATION {
		// Build state for classification
		typesafeState := buildTypesafeContext(hit, subject)

		// Call TypeSafe for reasoning questions
		reasoningResp, err := o.typesafeClient.Call(ctx, &typesafe.Request{
			State:     typesafeState,
			Model:     "jev-latest",
			Questions: typesafe.ReasoningQuestions,
		})
		if err != nil {
			o.logger.Error("typesafe reasoning call failed", "error", err.Error(), "result_id", screeningWcResult.ResultID)
			return nil, err
		}

		// Call TypeSafe for classification questions
		classificationResp, err := o.typesafeClient.Call(ctx, &typesafe.Request{
			State:     typesafeState,
			Model:     "jev-latest",
			Questions: typesafe.ClassificationQuestions,
		})
		if err != nil {
			o.logger.Error("typesafe classification call failed", "error", err.Error(), "result_id", screeningWcResult.ResultID)
			return nil, err
		}

		// Build AIRecommendation from responses
		aiRecommendation := &model.AIRecommendation{
			Reasoning:      mapAnswers(reasoningResp),
			Classification: mapAnswers(classificationResp),
		}
		screeningWcResult.AIRecommendation = aiRecommendation
	}

	// Save result with AI recommendation
	if err := o.screeningWcResultService.Save(ctx, screeningWcResult); err != nil {
		return nil, err
	}

	return screeningWcResult, nil
}

// buildTypesafeContext builds the context for TypeSafe API from WorldCheckHits
func buildTypesafeContext(hit *model.WorldCheckHits, subject *model.OgsEventSubject) map[string]any {
	return map[string]any{
		"subject": map[string]any{
			"name":        subject.Name,
			"dob":         subject.DOB,
			"nationality": subject.Nationality,
			"gender":      subject.Gender,
			"entity_type": subject.EntityType,
		},
		"hit": map[string]any{
			"reference_id":        hit.ReferenceID,
			"primary_name":        hit.PrimaryName,
			"matched_term":        hit.MatchedTerm,
			"provider":            hit.ProviderType,
			"category":            hit.UpdateCategory,
			"gender":              hit.KeyData.Gender,
			"nationality":         hit.KeyData.LocationDetails.Nationality,
			"location":            hit.KeyData.LocationDetails.Location,
			"aliases":             hit.Aliases,
			"comparison_fields":   hit.ComparisonData,
			"adverse_information": hit.FurtherInformation,
			"associate_count":     len(hit.ConnectionsAndRelationships),
			"associated_entities": hit.ConnectionsAndRelationships,
			"sources":             hit.Sources,
		},
	}
}

// mapAnswers converts TypeSafe Answer objects to model.AIAnswer
func mapAnswers(resp *typesafe.Response) map[string]model.AIAnswer {
	result := make(map[string]model.AIAnswer)
	for key, answer := range resp.Answers {
		result[key] = model.AIAnswer{
			Type:          answer.Type,
			Choice:        answer.Choice,
			Score:         answer.Score,
			Noul:          answer.Noul,
			Legend:        answer.Legend,
			Confidence:    answer.Confidence,
			Probabilities: answer.Probabilities,
		}
	}
	return result
}
