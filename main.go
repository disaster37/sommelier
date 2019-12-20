package main

import (
	"fmt"
	"log"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/options"
	_wineEstateRepo "github.com/sommelier/sommelier/v0/api/wine_estate/repository"
	"github.com/sommelier/sommelier/v0/middleware"
)

func init() {
	viper.SetConfigFile(`config.yml`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {

	// Init mongo repository
	mongoHost := viper.GetString(`mongo.host`)
	mongoPort := viper.GetString(`mongo.port`)
	mongoUser := viper.GetString(`mongo.user`)
	mongoPass := viper.GetString(`mongo.pass`)
	mongoDatabase := viper.GetString(`mongo.database`)
	mongoURL := fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)
	mongoCredential := options.Credential {
		AuthMechanism: "PLAIN",
		Username:      mongoUser,
		Password:      mongoPass,
	}
	mongoClientOpts :=  options.Client().ApplyURI(mongoURL).SetAuth(mongoCredential)
	ctx := context.TODO()
	mongoConn, err := mongo.Connect(ctx, mongoClientOpts)
	if err != nil {
		log.Fatal(err)
	}
	if err = mongoConn.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	mongoConnDB := mongoConn.Database(mongoDatabase)


	defer func() {
		err := mongoConn.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)
	_, err = _wineEstateRepo.NewMongoWineEstateRepository(mongoConnDB, viper.GetString(`mongo.collection.wine_estate`))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(e.Start(viper.GetString("server.address")))
}