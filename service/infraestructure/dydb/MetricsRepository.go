package dydb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"lambda-metrics-nir/service/application/domain"
	"lambda-metrics-nir/service/application/repositories"
)

type MetricsRepository struct {
	AwsSession       *session.Session
	TableName        string
	MemoryRepository repositories.IndexMemoryRepository
}

func NewMetricsRepository(awsSession *session.Session, tableName string, memoryRepository repositories.IndexMemoryRepository) repositories.DocumentMetricsRepository {
	return MetricsRepository{
		AwsSession:       awsSession,
		TableName:        tableName,
		MemoryRepository: memoryRepository,
	}
}

func (i MetricsRepository) Save(document domain.NormalizedDocument) error {

	item, err := dynamodbattribute.MarshalMap(document)
	i.MemoryRepository.Save("metrics"+document.Id, document)

	if err != nil {
		return err
	}

	svc := dynamodb.New(i.AwsSession)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(i.TableName),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		return err
	}

	return nil
}
