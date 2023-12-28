package authentication

import (
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

//verify token

func Verifytoken(tokenString string)error{
//grab jwt secret
godotenv.Load()
jwtSecret := os.Getenv("jwtSecret")
jwtSecretbyte := []byte(jwtSecret)

//verifying the token 

token,err := jwt.Parse(tokenString,func(token *jwt.Token)(interface{},error){

	return []byte(jwtSecretbyte),nil
})
if err != nil {
	log.Fatal(err.Error())
}
if !token.Valid {
	log.Fatal("Invalid token")
}
return nil

}




//middleware to validate token
func Middlewarefortokenvalidation(next http.HandlerFunc)http.HandlerFunc{
return func(res http.ResponseWriter,req *http.Request)  {
	tokenFromheaders := req.Header.Get("Authorization")
	if tokenFromheaders == "" {
		http.Error(res,"Unauthorized no token",http.StatusUnauthorized)
		return
	}

	validatedToken := Verifytoken(tokenFromheaders)

	if validatedToken != nil {
		res.Write([]byte(validatedToken.Error()))
		return
	}else{
		next.ServeHTTP(res,req)
	}


}




}