package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

var number int

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
		fmt.Fprintf(writer, "cannot decode json: %s", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, found := nameHobbyDB[user.Name]; found {
		log.Printf("this user %s with hobby %s already exists\n", user.Name, user.Hobby)
		writer.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprintf(writer, "this user already exists\n")
		return
	}
	nameHobbyDB[user.Name] = user.Hobby
	userID := getID()
	idUserDB[userID] = user
	fmt.Fprintf(writer, "hello, %s, your id is %d\n", user.Name, userID)
	log.Println(user, userID)

}

func getUser(writer http.ResponseWriter, req *http.Request) {
	urlStr := req.URL.String()
	urlParsed, err := url.Parse(urlStr)
	if err != nil {
		writer.WriteHeader(http.StatusPreconditionFailed)
		log.Printf("incorrect string query\n")
		return
	}
	query, err := url.ParseQuery(urlParsed.RawQuery)
	if err != nil {
		fmt.Fprintf(writer, "incorrect string query\n")
		writer.WriteHeader(http.StatusPreconditionFailed)
		log.Printf("incorrect string query\n")
		return
	}
	var userID int
	idNum, found := query["id"]
	if !found {
		writer.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprintf(writer, "incorrect key\n")
		log.Printf("incorrect key\n")
		return
	}

	if i, err := strconv.Atoi(idNum[0]); err == nil {
		userID = i
	} else {
		log.Printf("id is not correct\n")
		writer.WriteHeader(http.StatusPreconditionFailed)
		return
	}
	user, ok := idUserDB[userID]
	if !ok {
		writer.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprintf(writer, "no user with this id %d", userID)
		log.Printf("no user with this id %d", userID)
		return
	}
	nameHobbyDB[user.Name] = user.Hobby
	json.NewEncoder(writer).Encode(user)
}

func getID() int {
	number++
	return number
}
