package repositories

import "lambda-metrics-nir/service/application/domain"

type DocumentMetricsRepository interface {
	Save(domain.NormalizedDocument) error
}
