package users

import (
	
	"strings"
)

type Users struct {
	Username     string `json:"username" bson:"username"`
	Password     string `json:"password" bson:"password"`
	Email        string `json:"email" bson:"email"`
	Uniquestring string `json:"uniquestring" bson:"uniquestring"`
}

//methods
//ensure email is passed
func (email *Users)EnsureEmailisprovided()bool{
	if email.Email !="" {
		return true
	}
	return false
}

 


//Ensure no field is empty

func (user *Users) Ensurenoemptyfield() bool {
	if user.Email == "" || user.Password == "" || user.Username == ""  {
		return false
	} else {
		return true
	}
}

//ensure email is in right format

func (u *Users) Validateemail() bool {
	return strings.Contains(u.Email,"@")
}


//Ensure password has at least 6charchaters and a special 

func (pass *Users)CheckPasswordlengthandspecialcharchter()bool{
if len(pass.Password ) < 6 {
return false
}else{
	return true
}

}

