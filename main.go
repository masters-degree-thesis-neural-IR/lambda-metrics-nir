package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"lambda-metrics-nir/service/application/service"
	"lambda-metrics-nir/service/infraestructure/dto"
	"lambda-metrics-nir/service/infraestructure/dydb"
	zapplog "lambda-metrics-nir/service/infraestructure/log"
)

var TableName string
var AwsRegion string

func handler(ctx context.Context, event events.SQSEvent) error {

	logger := zapplog.NewZapLogger()

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion)},
	)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	repository := dydb.NewMetricsRepository(awsSession, TableName)
	service := service.NewMetricsService(logger, repository)

	logger.Info("Lambda Accepted Request")

	for _, message := range event.Records {

		var mp map[string]string
		json.Unmarshal([]byte(message.Body), &mp)

		doc := &dto.Document{}
		err := json.Unmarshal([]byte(mp["Message"]), doc)

		logger.Info("Documento recebido")
		logger.Info(doc)

		if err != nil {
			logger.Fatal(err.Error())
			return err
		}

		err = service.Create(doc.Id, doc.Title, doc.Body)

		if err != nil {
			logger.Fatal(err.Error())
			return err
		}

	}

	return nil
}

func main() {
	//TopicArn = "arn:aws:sns:us-east-1:149501088887:mestrado-document-created" //os.Getenv("BAR")
	AwsRegion = "us-east-1"
	TableName = "NIR_Metrics"
	lambda.Start(handler)
}
