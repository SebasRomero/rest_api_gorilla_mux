package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

type allTask []task

var tasks = allTask{
	{
		ID:      1,
		Name:    "task1",
		Content: "task1 content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "Insert a valid task")
		return
	}

	json.Unmarshal(reqBody, &newTask) //Asigning the information to &newTask

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json") //The header shows what type of data is (JSON)
	w.WriteHeader(http.StatusCreated)                  //This sends a status code to show that everything was ok
	json.NewEncoder(w).Encode(newTask)                 // Sends to the client the new Task
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for _, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...) //Here we just add all elements before and after the ID
			fmt.Fprintf(w, "The task with ID %v has been removed succesfully", taskID)
		}
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask task //This is de new task updated

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body) //Here we received the information to be able to work with
	if err != nil {
		fmt.Fprintf(w, "Please enter a valid data")
	}
	json.Unmarshal(reqBody, &updatedTask) //Here we asign that information in &updatedTaks

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...) //Here delete the task founded
			updatedTask.ID = taskID                   //Here asign the ID founded in the new task, is the same one
			tasks = append(tasks, updatedTask)        //Here add that updated task in the list

			fmt.Fprintf(w, "The task with ID %v has been updated successfully", taskID)
		}
	}
}

//Nothing to do, this is just a Hello in the index route
func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API")
}

func main() {
	//Here we stablished the routes and the whole management
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET") //28:27
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}
