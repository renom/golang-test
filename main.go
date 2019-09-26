package main

import (
	"flag"

	"encoding/json"
	"net/http"
	"strconv"

	valid "github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	conf "github.com/renom/golang-test/config"
	"github.com/renom/golang-test/errorf"
	"github.com/renom/golang-test/hash"
	"github.com/renom/golang-test/response"
	"github.com/renom/golang-test/router"
	taskModel "github.com/renom/golang-test/task"
)

var db *gorm.DB

func main() {
	var err error
	// Read db config from a file
	configPath := ""
	flag.StringVar(&configPath, "c", configPath, "A path to a .json config")
	flag.Parse()
	if configPath == "" {
		errorf.Exit("There's no config has been provided")
	}
	var config conf.Configuration
	if config, err = conf.LoadConfig(configPath); err != nil {
		errorf.Exit(err.Error())
	}

	// Connect to db
	db, err = gorm.Open("postgres", config.DbConnectionString())
	defer db.Close()
	if err != nil {
		errorf.Exit(err.Error())
	}

	// If the table doesn't exist, create it
	if !db.HasTable("tasks") {
		db.CreateTable(&taskModel.Task{})
	}

	// Routing
	r := router.NewRouter()
	r.HandleFunc("/v1/tasks", PostTask).Methods("POST")
	r.HandleFunc("/v1/tasks/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", GetTask).Methods("GET")

	// Start the server
	err = http.ListenAndServe(":"+strconv.Itoa(config.Port), r)
	if err != nil {
		errorf.Error(err.Error())
	}
}

// Create a task
func PostTask(w http.ResponseWriter, r *http.Request) {
	request := struct {
		Payload       string `json:"payload" valid:"length(1|1024),required"`
		HashRoundsCnt int    `json:"hash_rounds_cnt" valid:"range(1|32),required"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		response.RespondError(w, err.Error())
		return
	}
	if result, err := valid.ValidateStruct(request); !result {
		response.RespondError(w, err.Error())
		return
	}

	task := taskModel.Task{
		Payload:       request.Payload,
		HashRoundsCnt: request.HashRoundsCnt,
		Status:        taskModel.StatusInProgress,
	}
	if err := db.Create(&task).Error; err != nil {
		response.RespondError(w, err.Error())
		return
	}
	response.Respond(w, http.StatusCreated, task)
	go func() {
		task.Hash = hash.Calc(task.Payload, task.HashRoundsCnt)
		task.Status = taskModel.StatusFinished
		db.Save(&task)
	}()
}

// Return a single task
func GetTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var task taskModel.Task
	if err := db.First(&task, "id = ?", id).Error; err != nil {
		response.RespondError(w, err.Error())
		return
	}
	response.Respond(w, http.StatusOK, task)
}
