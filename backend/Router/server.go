package router

import (
	"net/http"
	"os"

	"fmt"

	"github.com/gorilla/mux"
	authentication "github.com/jamesmukumu/diary2024/backend/Middlewares/Authentication"
	"github.com/jamesmukumu/diary2024/backend/Middlewares/helpers"
	diarycontroller "github.com/jamesmukumu/diary2024/backend/controllers/DiaryController"
	userscontroller "github.com/jamesmukumu/diary2024/backend/controllers/usersController"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func Serversetup() {

   godotenv.Load()

   PORT := os.Getenv("port")



	Router := mux.NewRouter()
     Handler := cors.Default().Handler(Router)


//Routes user

Router.HandleFunc("/post/user",userscontroller.Postuser).Methods("POST")
Router.HandleFunc("/login/user",userscontroller.Loginuser).Methods("POST")
Router.HandleFunc("/change/password",helpers.ConfirmemailExistenceandgetToken(userscontroller.UpdatePassword)).Methods("PUT")
Router.HandleFunc("/delete/account",userscontroller.Deleteaccount).Methods("DELETE")
Router.HandleFunc("/login/uniquestring",userscontroller.Loginwithuniquestring).Methods("POST")






//Routes for diary
Router.HandleFunc("/add/event",authentication.Middlewarefortokenvalidation(diarycontroller.AddEvent)).Methods("POST")
Router.HandleFunc("/fetch/all/events",authentication.Middlewarefortokenvalidation(diarycontroller.GetEvents)).Methods("GET")
Router.HandleFunc("/fetch/single/event",authentication.Middlewarefortokenvalidation(diarycontroller.Getasingleevent)).Methods("GET")
Router.HandleFunc("/update/message",authentication.Middlewarefortokenvalidation(diarycontroller.Updatemessagebasedonsing)).Methods("PUT")
Router.HandleFunc("/delete/event",authentication.Middlewarefortokenvalidation(diarycontroller.Deleteevent)).Methods("DELETE")
Router.HandleFunc("/delete/all/events",authentication.Middlewarefortokenvalidation(diarycontroller.Deleteallevents)).Methods("DELETE")
 

 

     
    fmt.Printf("Server listening for request at %s",PORT)
	 http.ListenAndServe(":"+PORT,Handler)



}