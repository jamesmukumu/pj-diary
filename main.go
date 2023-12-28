package main

import (
	router "github.com/jamesmukumu/diary2024/backend/Router"
	dbdiary "github.com/jamesmukumu/diary2024/backend/db/dbDiary"
	"github.com/jamesmukumu/diary2024/backend/db/dbusers"
)

func main() {
	dbdiary.Connectdbdiary()
     dbusers.Connectdbuses()
	router.Serversetup()

}