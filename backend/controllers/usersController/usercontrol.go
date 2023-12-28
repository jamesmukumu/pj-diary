package userscontroller

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	users "github.com/jamesmukumu/diary2024/backend/Schema/Users"
	"github.com/jamesmukumu/diary2024/backend/db/dbusers"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)
var Changedpassword users.Users
// create an instance of the class
var MyUser users.Users

//generate uniqueString

func GenerateUniquestring(uniquestring int) (string, error) {
	buffer := make([]byte, uniquestring)

	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	} else {
		return base64.URLEncoding.EncodeToString(buffer)[:uniquestring], nil

	}

}

// create token
func CreateToken(usernamefetched string) (string, error) {
	godotenv.Load()
	jwtSecret := os.Getenv("jwtSecret")
	secret := []byte(jwtSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": usernamefetched,
		"exp":      time.Now().Add(time.Minute * 10).Unix(),
	})
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return signedToken, nil

}
//sent reset Password link// send unique string on sign up
func SendResetlink() {
	godotenv.Load()
   emaiAddress := os.Getenv("myEmail")
  var emailPassowrd = os.Getenv("passwordEmail")
   mail := gomail.NewMessage()
   mail.SetHeader("From", emaiAddress)
   mail.SetHeader("To", Changedpassword.Email)
   mail.SetHeader("Subject", "Reset Link")
 

   mail.SetBody("text/plain", "Password Reset Link this is tets" )

   // setting up the dialer
   dialer := gomail.NewDialer("smtp.gmail.com", 587,emaiAddress, emailPassowrd)
   dialer.TLSConfig = &tls.Config{
	   InsecureSkipVerify: false,
	   ServerName:         "smtp.gmail.com",
   }

   // send the email
   if err := dialer.DialAndSend(mail); err != nil {
	   log.Fatal(err.Error())
   }
}















// send unique string on sign up
func Senduniquestring(contentData string) {
     godotenv.Load()
    emaiAddress := os.Getenv("myEmail")
   var emailPassowrd = os.Getenv("passwordEmail")
	mail := gomail.NewMessage()
	mail.SetHeader("From", emaiAddress)
	mail.SetHeader("To", MyUser.Email)
	mail.SetHeader("Subject", "subject")
	htmlBody := "<html><body><p style='color: red;'>" + contentData + "</p></body></html>"

	mail.SetBody("text/html", htmlBody)

	// setting up the dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587,emaiAddress, emailPassowrd)
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         "smtp.gmail.com",
	}

	// send the email
	if err := dialer.DialAndSend(mail); err != nil {
		log.Fatal(err.Error())
	}
}

func Postuser(res http.ResponseWriter, req *http.Request) {

	//decode json
	decodedJson := json.NewDecoder(req.Body).Decode(&MyUser)
	if decodedJson != nil {
		http.Error(res, "Error in decoding json", http.StatusInternalServerError)
		return
	}

	if !MyUser.CheckPasswordlengthandspecialcharchter() {
		res.Write([]byte("PAssword must be 6 charchters"))
		return
	}

	if !MyUser.Ensurenoemptyfield() {
		res.Write([]byte("No field should be empty"))
		return
	}

	if !MyUser.Validateemail() {
		res.Write([]byte("Wrong Email format"))
		return
	}

	hashedPassword, errr := bcrypt.GenerateFromPassword([]byte(MyUser.Password), 10)

	if errr != nil {
		log.Fatal(errr.Error())
		return
	}

	MyUser.Password = string(hashedPassword)

	MyUser.Uniquestring, _ = GenerateUniquestring(20)

	insertedData, err := dbusers.MongodbUsers.InsertOne(context.Background(), MyUser)
	Senduniquestring(MyUser.Uniquestring)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
			res.Write([]byte("Username or Email already exists"))
			return
		}

		res.Write([]byte(err.Error()))

		return
	}

	Tokengenerated, _ := CreateToken(MyUser.Username)
	res.Header().Set("Authorization", Tokengenerated)
	json.NewEncoder(res).Encode(map[string]interface{}{"data": insertedData})

	//set headers

}

//login

func Loginuser(res http.ResponseWriter, req *http.Request) {
	var Loggeduser users.Users

	jsonDecoded := json.NewDecoder(req.Body).Decode(&Loggeduser)
	if jsonDecoded != nil {
		log.Fatal(jsonDecoded.Error())
		return
	}

	filter := bson.M{
		"username": Loggeduser.Username,
	}

	if !Loggeduser.CheckPasswordlengthandspecialcharchter() {
		http.Error(res, "Password not long enough", http.StatusInternalServerError)
		return
	}

	var Passworddb users.Users

	matchingUsername := dbusers.MongodbUsers.FindOne(context.Background(), filter).Decode(&Passworddb)
	if matchingUsername != nil {
		http.Error(res, "Username not found", http.StatusForbidden)
		return
	}

	matchingPasswordandUsername := bcrypt.CompareHashAndPassword([]byte(Passworddb.Password), []byte(Loggeduser.Password))

	if matchingPasswordandUsername != nil {
		http.Error(res, "Invalid password", http.StatusOK)
		return
	} else {
		tokenLogin, _ := CreateToken(Loggeduser.Username)
		res.Header().Set("Authorization", tokenLogin)
		http.Error(res, "Accepted", http.StatusOK)
	}
}

//change password

func UpdatePassword(res http.ResponseWriter, req *http.Request) {
	

	emailFetched := req.URL.Query().Get("email")
	if emailFetched == "" {
		res.Write([]byte("Please provide an email for query"))
		return
	}

	Changedpassword.Email = emailFetched
	if !Changedpassword.Validateemail() {
		res.Write([]byte("Incorrect email format"))
		return
	}

	filter := bson.M{
		"email": emailFetched,
	}

	//decode the password body
	Passworddecoded := json.NewDecoder(req.Body).Decode(&Changedpassword)
	if Passworddecoded != nil {
		log.Fatal(Passworddecoded.Error())
		return
	}

	//hashPassword
	hashedPasswordinserted, _ := bcrypt.GenerateFromPassword([]byte(Changedpassword.Password), 10)
	Changedpassword.Password = string(hashedPasswordinserted)
	update := bson.M{
		"$set": bson.M{"password": Changedpassword.Password},
	}

	passwordUpdate := dbusers.MongodbUsers.FindOneAndUpdate(context.Background(), filter, update).Decode(&Changedpassword)

	if passwordUpdate != nil {
		res.Write([]byte(passwordUpdate.Error()))
		return
	}
	res.Write([]byte("Password changed"))

}

// delete Account
func Deleteaccount(res http.ResponseWriter, req *http.Request) {
	//create an instance
	var DeletedAccount users.Users

	Emailquery := req.URL.Query().Get("email")

	if Emailquery == "" {
		json.NewEncoder(res).Encode("Provide Email for query")
		return
	}

	DeletedAccount.Email = Emailquery

	if !DeletedAccount.Validateemail() {
		res.Write([]byte("Wrong email format"))
		return
	}

	filter := bson.M{
		"email": Emailquery,
	}

	ActualdeletedAaccount, err := dbusers.MongodbUsers.DeleteOne(context.Background(), filter)
	if err != nil {
		res.Write([]byte(err.Error()))
		return
	}
	if ActualdeletedAaccount.DeletedCount == 0 {
		res.Write([]byte("Sorry No account to be deleted"))
		return
	} else {
		res.Write([]byte("Account Deleted"))
	}
}

//login with unique string

func Loginwithuniquestring(res http.ResponseWriter, req *http.Request) {

	var Userdb users.Users
	decodedJson := json.NewDecoder(req.Body).Decode(&Userdb)

	if decodedJson != nil {
		panic(decodedJson.Error())

	}

	filter := bson.M{
		"uniquestring": Userdb.Uniquestring,
	}
	matchingUniqueString := dbusers.MongodbUsers.FindOne(context.Background(), filter).Decode(&Userdb)
	if matchingUniqueString != nil {
		res.Write([]byte(matchingUniqueString.Error()))
		return
	} else {
		token, _ := CreateToken(Userdb.Username)

		res.Header().Set("Authorization", token)
		res.Write([]byte("Login success"))

	}
}
