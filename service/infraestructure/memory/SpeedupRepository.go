package memory

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"lambda-metrics-nir/service/application/domain"
	"lambda-metrics-nir/service/application/exception"
	"lambda-metrics-nir/service/application/repositories"
	"lambda-metrics-nir/service/infraestructure/speedup"
	"log"
	"time"
)

type SpeedupRepository struct {
}

func NewSpeedupRepository() repositories.IndexMemoryRepository {
	return &SpeedupRepository{}
}

func (r *SpeedupRepository) Save(term string, document domain.NormalizedDocument) error {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("172.31.2.165:9000", grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
		return exception.ThrowValidationError("Not is possible connect to RCP Server.")
	}
	defer conn.Close()

	client := speedup.NewDataServiceClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	value, err := json.Marshal(document)
	response, err := client.SetData(ctx, &speedup.RequestDataKeyValue{
		Key:   term,
		Value: string(value),
	})

	if ctx.Err() == context.Canceled {
		log.Println(err.Error())
		return exception.ThrowValidationError("RPC Client cancelled, abandoning.")
	}

	if err != nil {
		return err
	}

	if response.GetException() != "" {
		log.Println(err.Error())
		return exception.ThrowValidationError(response.GetException())
	}

	return nil

}
