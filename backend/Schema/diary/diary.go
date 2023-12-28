package diary

import (
	"log"
	"time"
)

type Diary struct {
	Day string `json:"day" bson:"day"`
	Message string `json:"message" bson:"message"`
	Signature string `json:"sign" bson:"sign"`
	Daytimeintimeformat time.Time `json:"timeform" bson:"timeform"`
}


func (diary *Diary)Checkmessagenotempty()bool{
	if  diary.Message == ""  {
	
			log.Println("Field should not be empty")
			return false
	}else{
		return true
	}
	
	}


//create a method to check non empty fields


func (diary *Diary)Checkemptyfields()bool{
if  diary.Message == "" || diary.Signature == "" {

		log.Println("Field should not be empty")
		return false
}else{
	return true
}

}


//method to emsure Signature is 20 charchters long

func(signature *Diary)Ensuresignatueislongenough()bool{
if len(signature.Signature ) !=  20  {
	
log.Print("Signature must be 20 charachters long")
return false
}else{
	return true
}



} 




//check day format ensure it has 2024-mm-dd as a string

func (day*Diary)Validatedateformat()bool{


	dateLayout := "2006-01-01" //YYYY/MM/dd
    _,err := time.Parse(dateLayout,day.Day)
	if err != nil {
		log.Println("Incorrect date format")
		return false
	}else{
		return true
	}



}