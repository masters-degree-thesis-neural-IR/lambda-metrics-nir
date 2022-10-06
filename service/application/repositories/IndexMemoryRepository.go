package repositories

import "lambda-metrics-nir/service/application/domain"

type IndexMemoryRepository interface {
	Save(term string, document domain.NormalizedDocument) error
}
