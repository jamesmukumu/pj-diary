package helpers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	users "github.com/jamesmukumu/diary2024/backend/Schema/Users"
	"github.com/jamesmukumu/diary2024/backend/db/dbusers"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

//creating a class of the response
type Response struct{
	Confirmdeletion string `json:"confirmdeletion"`
}

//method validate if confirm deletion is either a yes or a no
func (res *Response) ValidateDeletion () bool{
if res.Confirmdeletion != "yes" && res.Confirmdeletion != "no" {
	return false
}else if res.Confirmdeletion == "" {
	return false
}
return true

}//

//create Token that a user will use to acceess the reset api
func Createtokenonemailvalidation(email string)(string,interface{}){
godotenv.Load()
jwtSecretforreset := os.Getenv("jwtSecretforreset")
byteForm := []byte(jwtSecretforreset)


//create actual token
token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
	"username":email,
	"exp":time.Now().Add(time.Minute * 5).Unix(),
})
//sign the token to secret
signedToken,err := token.SignedString(byteForm)
if err != nil {
	log.Fatal(err.Error())
	return "",err.Error()
}

	return signedToken,nil



}




//middleware to confirm email before accesing reset link

func ConfirmemailExistenceandgetToken(next http.HandlerFunc)http.HandlerFunc{
return func(res http.ResponseWriter,req *http.Request){
	//create an instance 
	var EmailConfirmation users.Users
	jsonDecodedforemail := json.NewDecoder(req.Body).Decode(&EmailConfirmation)
     
    if jsonDecodedforemail != nil {
		log.Fatal(jsonDecodedforemail.Error())
		return
	}
    

	//create a method to ensure Email is passed
	if !EmailConfirmation.EnsureEmailisprovided() {
		message := map[string]string{"message":"Email should be provided"}
		jsonMeesage,_ := json.Marshal(message)
		res.Write([]byte(jsonMeesage))
		return
	}
	filter := bson.M{
		"email":EmailConfirmation.Email,
	}

	matchingEmail := dbusers.MongodbUsers.FindOne(context.Background(),filter).Decode(&EmailConfirmation)
    
	if matchingEmail != nil {
		http.Error(res,matchingEmail.Error(),http.StatusInternalServerError)
		return
	}else{
		createdToken,_ := Createtokenonemailvalidation(EmailConfirmation.Email)
		res.Header().Set("Authorization",createdToken)
		next.ServeHTTP(res,req)
	} 





}




}





