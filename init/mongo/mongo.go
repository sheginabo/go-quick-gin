package mongo

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Module struct {
	Client *mongo.Client
}

func NewModule() *Module {
	clientOptions := options.Client().ApplyURI(viper.GetString("MONGO_URI"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	// 測試連接是否成功
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to Ping to MongoDB")
	}

	fmt.Println("Connected to MongoDB!")

	return &Module{
		Client: client,
	}
}
