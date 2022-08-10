package main

import (
	"dts-task/controller"
	"dts-task/model"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	port := os.Getenv("PORT")

	db, err := gorm.Open(sqlite.Open("./task.db"), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	db.AutoMigrate(&model.Task{})

	taskController := controller.NewTaskController{
		Db: db,
	}
	http.HandleFunc("/", taskController.Index)
	http.HandleFunc("/create", taskController.Create)
	http.HandleFunc("/edit/", taskController.Edit)
	http.HandleFunc("/store", taskController.Store)
	http.HandleFunc("/update", taskController.Update)
	http.HandleFunc("/delete/", taskController.Delete)

	http.ListenAndServe(":"+port, nil)
}
