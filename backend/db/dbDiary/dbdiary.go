package dbdiary

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongodbdiary *mongo.Collection
var Client *mongo.Client 

func Connectdbdiary() {
	godotenv.Load()
	myconnectionString :=os.Getenv("mongoConnection")
	connectionString := options.Client().ApplyURI(myconnectionString)

	Connection, err := mongo.Connect(context.Background(), connectionString)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	Client = Connection

	Mongodbdiary = Client.Database("Diary").Collection("mydiary")

	//set up somefield unique

	dateOption := options.Index().SetUnique(true)
	dateIndex := mongo.IndexModel{
		Keys:    bson.M{"day":1},
		Options: dateOption,
	}

	_, errr := Mongodbdiary.Indexes().CreateOne(context.Background(), dateIndex)
	if errr != nil {
		log.Fatal(errr.Error())
		return
	}

	
	

	fmt.Println("Connected to DB Diary successfully")

}