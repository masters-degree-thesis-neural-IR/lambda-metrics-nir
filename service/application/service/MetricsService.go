package service

import (
	"lambda-metrics-nir/service/application/domain"
	"lambda-metrics-nir/service/application/logger"
	"lambda-metrics-nir/service/application/nlp"
	"lambda-metrics-nir/service/application/repositories"
	"lambda-metrics-nir/service/application/usecases"
)

type MetricsService struct {
	Logger                    logger.Logger
	DocumentMetricsRepository repositories.DocumentMetricsRepository
}

func NewMetricsService(logger logger.Logger, documentMetricsRepository repositories.DocumentMetricsRepository) usecases.NormalizedDocumentUc {

	return MetricsService{
		Logger:                    logger,
		DocumentMetricsRepository: documentMetricsRepository,
	}
}

func (m MetricsService) Create(id string, title string, body string) error {

	tokens := nlp.Tokenizer(body, true)
	normalizedTokens, err := nlp.RemoveStopWords(tokens, "en")

	if err != nil {
		m.Logger.Error(err.Error())
	}

	normalizedDocument := domain.NormalizedDocument{
		Id:     id,
		Length: len(normalizedTokens),
		Tf:     nlp.TermFrequency(normalizedTokens),
	}

	m.DocumentMetricsRepository.Save(normalizedDocument)

	return nil
}
