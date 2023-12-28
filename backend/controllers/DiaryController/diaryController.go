package diarycontroller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jamesmukumu/diary2024/backend/Middlewares/helpers"
	"github.com/jamesmukumu/diary2024/backend/Schema/diary"
	dbdiary "github.com/jamesmukumu/diary2024/backend/db/dbDiary"

	"go.mongodb.org/mongo-driver/bson"
)

func AddEvent(res http.ResponseWriter, req *http.Request){
var mydiary diary.Diary

decodeJson := json.NewDecoder(req.Body).Decode(&mydiary)
if decodeJson != nil {
	http.Error(res,"Error decoding json",http.StatusInternalServerError)
	return
}

if !mydiary.Validatedateformat() {
	http.Error(res,"Wrong date format.",http.StatusOK)
	return
}  

 

timeFormtted,err := time.Parse("2006-01-02",mydiary.Day)
if err != nil {
	log.Fatal(err.Error())
	return
}
mydiary.Daytimeintimeformat = timeFormtted


if !mydiary.Checkemptyfields() {
	http.Error(res,"No empty field allowed",http.StatusBadRequest)
	return
}



if !mydiary.Ensuresignatueislongenough() {
	http.Error(res,"String must be long enough",http.StatusBadRequest)
	return
}

  



insertedEvent,errr := dbdiary.Mongodbdiary.InsertOne(context.Background(),mydiary)
if errr != nil {
if strings.Contains(errr.Error(),"E11000 duplicate key error collection") {
	res.Write([]byte("date already in use"))
	return
}

 


	res.Write([]byte(errr.Error()))
	return
}
  
json.NewEncoder(res).Encode(map[string]interface{}{"data":insertedEvent})


}




//get Events based on signature

func GetEvents(res http.ResponseWriter,req *http.Request){
  var Myevents []diary.Diary


//query using the signature
signatureQuery := req.URL.Query().Get("sign")
if len(signatureQuery) != 20  || signatureQuery == ""{
	http.Error(res,"Query signature should be 20 charachters long and not blank",http.StatusAccepted)
	return
}


 


filter := bson.M{
	"sign":signatureQuery,
}




fetchedResults,err := dbdiary.Mongodbdiary.Find(context.Background(),filter)
if err != nil {
	res.Write([]byte(err.Error()))
	return
}


finalResults := fetchedResults.All(context.Background(),&Myevents)
if finalResults != nil {
log.Fatal(finalResults.Error())

}else{
	json.NewEncoder(res).Encode(map[string]interface{}{"data":Myevents})
}





}





func Getasingleevent(res http.ResponseWriter,req *http.Request){
	var Singlevent diary.Diary


//query on both signature and day
signatureQuery := req.URL.Query().Get("sign")
dayQuery := req.URL.Query().Get("day")


//check that signature query is 20 long,dayQuery is format yyyy/mm/dd and none is blank
if len(signatureQuery) != 20 {
	http.Error(res,"Signature should be 20 charachters long",http.StatusAccepted)
	return
}else if  dayQuery == ""  || signatureQuery == ""{
	http.Error(res,"Please provide signature and day query",http.StatusAccepted)
	return
}



filter := bson.M{
	"sign":signatureQuery,
	"day":dayQuery,
}


matchingDocument := dbdiary.Mongodbdiary.FindOne(context.Background(),filter)
decodedResult := matchingDocument.Decode(&Singlevent)
if decodedResult != nil {
	
	res.Write([]byte(decodedResult.Error()))
	return
}

marshelledDecodedresult,_ := json.Marshal(&Singlevent)
res.Write([]byte(marshelledDecodedresult))
}
 




//update the message based on signature and day
func Updatemessagebasedonsing(res http.ResponseWriter,req *http.Request){
var Messageupdate diary.Diary

//decode the body of request
decodedRequest := json.NewDecoder(req.Body).Decode(&Messageupdate)
if decodedRequest != nil {
	log.Fatal(decodedRequest.Error())
}

dayQuery := req.URL.Query().Get("day")
signatureQuery := req.URL.Query().Get("sign")


if len(signatureQuery) != 20 {
	http.Error(res,"Signature should be 20 charachters long",http.StatusAccepted)
	return
}
 if  dayQuery == ""  || signatureQuery == ""{
	http.Error(res,"Please provide signature and day query",http.StatusAccepted)
	return
}



//filter on day and signature


filter := bson.M{
	"sign":signatureQuery,
	"day":dayQuery,
}

update:= bson.M{
	"$set":bson.M{"message":Messageupdate.Message},
}



actualUpdate := dbdiary.Mongodbdiary.FindOneAndUpdate(context.Background(),filter,update).Decode(&Messageupdate)

if actualUpdate != nil {
	res.Write([]byte(actualUpdate.Error()))
	return
}

info := map[string]string{"message":"Your message was updated sucessfully"}
infojson,_ := json.Marshal(info) 
res.Write(infojson)
}





//delete an Event based

func Deleteevent(res http.ResponseWriter,req *http.Request){
// var Eventtodelete diary.Diary

//delete on signature && date
signatureQuery := req.URL.Query().Get("sign")
dayQuery := req.URL.Query().Get("day")

if len(signatureQuery) != 20 {
	http.Error(res,"Signature should be 20 charachters long",http.StatusAccepted)
	return
}
 if  dayQuery == ""  || signatureQuery == ""{
	http.Error(res,"Please provide signature and day query",http.StatusAccepted)
	return
}



filter := bson.M{
"sign":signatureQuery,
"day":dayQuery,
}

Deletedevent,err := dbdiary.Mongodbdiary.DeleteOne(context.Background(),filter)
if err != nil {
res.Write([]byte(err.Error()))
return
}

if Deletedevent.DeletedCount == 0 {
	message := map[string]string{"message":"No record to be deleted found"}
	jsonMessage,_ := json.Marshal(message)
	res.Write(jsonMessage)
	return
}else{
	json.NewEncoder(res).Encode(Deletedevent)
}

}




//clear all events based on signature 

func Deleteallevents(res http.ResponseWriter,req *http.Request){
	//create an instance of the Confirmation
	var Confirmationtodelete helpers.Response
//decode the Confirmation wrt to confirmationDelete
decodedConfirmation := json.NewDecoder(req.Body).Decode(&Confirmationtodelete)
if decodedConfirmation != nil {
	res.Write([]byte(decodedConfirmation.Error()))
	return
}
if !Confirmationtodelete.ValidateDeletion() {
	http.Error(res,"Please provide yes or no before deletion",http.StatusAccepted)
	return
}

	//delete on signature && date
	signatureQuery := req.URL.Query().Get("sign")
	
	
	if len(signatureQuery) != 20 {
		http.Error(res,"Signature should be 20 charachters long or provide your signature accordingly",http.StatusAccepted)
		return
	}
	
	
	filter := bson.M{
	"sign":signatureQuery,
	
	}
	
	Deletedevent,err := dbdiary.Mongodbdiary.DeleteMany(context.Background(),filter)
	if err != nil {
	res.Write([]byte(err.Error()))
	return
	}
	
	if Deletedevent.DeletedCount == 0 {
		message := map[string]string{"message":"No record to be deleted found"}
		jsonMessage,_ := json.Marshal(message)
		res.Write(jsonMessage)
		return
	}else{
		json.NewEncoder(res).Encode(Deletedevent)
	}
	
	}
	
	