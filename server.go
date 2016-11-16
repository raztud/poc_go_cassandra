package main

import (
	"encoding/json"
	"net/http"
	"time"
	"log"

	"github.com/gorilla/mux"
	"github.com/gocql/gocql"

	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Note struct {
	Id	    string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   int64     `json:"createdOn"`
}

var session *gocql.Session


func PostHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var note Note

	err := json.NewDecoder(request.Body).Decode(&note)
	if err != nil {
		panic(err)
	}

	note.CreatedOn = time.Now().Unix()
	session := getSession()

	cql := `INSERT INTO notes (id, title, description, createdon) VALUES (?, ?, ?, ?)`
	uuid, _ := gocql.RandomUUID()
	note.Id = uuid.String()

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

func GetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var note Note
	vars := mux.Vars(request)
	id := vars["id"]

	var description, title string
	var createdon int64

	if err := session.Query(`SELECT  createdon, description, title FROM notes WHERE id = ?`, id).
		Consistency(gocql.One).Scan(&createdon, &description, &title); err != nil {
		log.Println(err, id)
	}

	note.Id = id
	note.Description = description
	note.Title = title
	note.CreatedOn = createdon

	responseWriter.Header().Set("Content-Type", "application/json")

	j, err := json.Marshal(note)
	if err != nil {
		panic(err)
	}

	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write(j)
}

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

	//createDbSession()
	region := "us-east-1"
	svc := dynamodb.New(awsSession.New(&aws.Config{Region: aws.String(region)}))
	result, err := svc.ListTables(&dynamodb.ListTablesInput{})

	if err != nil {
	    log.Println(err)
	    return
	}

	for _, table := range result.TableNames {
	    log.Println(*table)
	}

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/notes", PostHandler).Methods("POST")
	r.HandleFunc("/api/notes/{id}", GetHandler).Methods("GET")
	server := &http.Server {
		Addr: ":8080",
		Handler: r,
	}

	log.Println("Started...")

	server.ListenAndServe()
}
