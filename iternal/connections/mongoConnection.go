package connections

import (
	"context"
	"fmt"
	_ "fmt"
	"github.com/SolidShake/wetherboy-tg-bot/iternal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

type Subs struct {
	ChatId         int64
	Coord          tgbotapi.Location
	LastUpdateDate string
}

type MongoConnection struct {
	client   *mongo.Client
	database string
}

func (c *MongoConnection) ConnectMongo() {
	cli, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + config.GetConfig().MongoDb.Host + ":" + config.GetConfig().MongoDb.Port))
	database := config.GetConfig().MongoDb.Database

	if err != nil {
		log.Fatal(err)
	}

	c.database = database
	c.client = cli

	err = c.client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	err = c.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Print(c.client.Database(c.database).Name())

}

func (c *MongoConnection) GetDbName() interface{} {
	return c.client.Database(c.database).Name()
}

func (c *MongoConnection) Disconnect() {
	c.client.Disconnect(context.TODO())
}

func (c *MongoConnection) AddSub(chat_id int64, location tgbotapi.Location) {
	new_sub := Subs{ChatId: chat_id, Coord: location}
	collection := c.client.Database(c.database).Collection("SUBERS")
	insertResult, err := collection.InsertOne(context.TODO(), new_sub)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted document: ", insertResult.InsertedID)
}
