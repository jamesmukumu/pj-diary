package dbusers

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

var MongodbUsers *mongo.Collection
var Client *mongo.Client

     

func Connectdbuses(){
	godotenv.Load()
	myconnectionString := os.Getenv("mongoConnection")
connectionString := options.Client().ApplyURI(myconnectionString)

Connection,err := mongo.Connect(context.TODO(),connectionString)

if err !=nil {
	log.Fatal(err.Error())
	return
}



Client = Connection

MongodbUsers = Client.Database("Gloryyear").Collection("People")


//set up somefield unique

usernameOption := options.Index().SetUnique(true)
usernameIndex := mongo.IndexModel{
	Keys: bson.M{"username":1},
	Options: usernameOption,
}

_,errr := MongodbUsers.Indexes().CreateOne(context.Background(),usernameIndex)
if errr != nil {
	log.Fatal(errr.Error())
	return
}





emailOption := options.Index().SetUnique(true)
emailIndex := mongo.IndexModel{
	Keys: bson.M{"email":1},
	Options: emailOption,
}

_,fault := MongodbUsers.Indexes().CreateOne(context.Background(),emailIndex)
if fault != nil {
	log.Fatal(fault.Error())
	return
}
 
 


fmt.Println("Connected to DB successfully")


}