package usecases

type NormalizedDocumentUc interface {
	Create(id string, title string, body string) error
}
