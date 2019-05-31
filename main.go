package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

/*
example:
	*{"name": "Ilya", "hobby": "Programming"}
	*localhost:9000/?id=1
*/

// User represents user with name and hobby strings fields.
type User struct {
	Name  string `json:"name"`
	Hobby string `json:"hobby"`
}

var idUserDB = make(map[int]User)

var nameHobbyDB = make(map[string]string)

var nextID int

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe("localhost:9000", nil)
	if err != nil {
		log.Fatal("failed to start server", err)
	}
}

func handler(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getUser(writer, req)
	case http.MethodPost:
		addUser(writer, req)
	}
}

func addUser(writer http.ResponseWriter, req *http.Request) {
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "cannot decode json: %s", err)
		return
	}
	if _, found := nameHobbyDB[user.Name]; found {
		log.Printf("this user %s with hobby %s already exists", user.Name, user.Hobby)
		writer.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprintf(writer, "this user already exists")
		return
	}
	nameHobbyDB[user.Name] = user.Hobby
	userID := getNextID()
	idUserDB[userID] = user
	fmt.Fprintf(writer, "hello, %s, your id is %d", user.Name, userID)
	log.Println(user, userID)

}

func getUser(writer http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	idNum, found := query["id"]
	if !found {
		writer.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprintf(writer, "incorrect key")
		log.Printf("incorrect key")
		return
	}
	userID, err := strconv.Atoi(idNum[0])
	if err != nil {
		writer.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprintf(writer, "id is not correct")
		log.Printf("id is not correct")
		return
	}
	user, ok := idUserDB[userID]
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "no user with this id %d", userID)
		log.Printf("no user with this id %d", userID)
		return
	}
	_ = json.NewEncoder(writer).Encode(user)
}

func getNextID() int {
	nextID++
	return nextID
}
