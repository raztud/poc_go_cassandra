package main

import (
	"encoding/json"
	"net/http"
	//"strconv"
	"time"
	"log"

	"github.com/gorilla/mux"
	"github.com/gocql/gocql"
)

type Note struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   int64     `json:"createdon"`
}

//store
//var store = make(map[string]Note)
//var id int = 0

var session *gocql.Session



func PostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var note Note

	err := json.NewDecoder(request.Body).Decode(&note)
	if err != nil {
		panic(err)
	}

	note.CreatedOn = time.Now().Unix()
	session := getSession()

	//k := strconv.Itoa(id)
	//store[k] = note

	cql := `INSERT INTO notes (id, title, description, createdon) VALUES (?, ?, ?, ?)`
	uuid, _ := gocql.RandomUUID()
	//
	err1 := session.Query(cql, uuid, note.Title, note.Description, note.CreatedOn).Exec()
	if err1 != nil {
		log.Println(err1)
		log.Println(cql)
		log.Println(uuid, note.Title, note.Description, time.Now().Unix())
	}

	j, err := json.Marshal(note)
	if err != nil {
		panic(err)
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusCreated)
	responseWriter.Write(j)

}

//func GetHandler(responseWriter http.ResponseWriter, request *http.Request) {
//	var notes []Note
//
//	for _, v := range store {
//		notes = append(notes, v)
//	}
//	responseWriter.Header().Set("Content-Type", "application/json")
//
//	j, err := json.Marshal(notes)
//	if err != nil {
//		panic(err)
//	}
//
//	responseWriter.WriteHeader(http.StatusOK)
//	responseWriter.Write(j)
//}

func createDbSession() {
	var err error

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "razvan"
	session, err = cluster.CreateSession()
	if err != nil {
		log.Panic("Could not open DB Connection")
	}

	//defer session.Close()
}

func getSession() *gocql.Session {
	if session == nil {
		createDbSession()
	}

	return session
}

func main() {

	createDbSession()

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/notes", PostHandler).Methods("POST")
	//r.HandleFunc("/api/notes", GetHandler).Methods("GET")
	server := &http.Server {
		Addr: ":8080",
		Handler: r,
	}

	log.Println("Started...")

	server.ListenAndServe()
}
