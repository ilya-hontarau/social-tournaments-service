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
	?  how should i respond to inccorect user input
example:
	*{"name": "Ilya", "hobby": "Programming"}
	*localhost:9000/?id=1
*/

// User represents user with name and hobby strings fields.
type User struct {
	Name  string
	Hobby string
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
	urlParsed, _ := url.Parse(urlStr)
	querry, errQuerry := url.ParseQuery(urlParsed.RawQuery)
	if errQuerry != nil {
		fmt.Fprintf(writer, "incorrect string querry\n")
		log.Printf("incorrect string querry\n")
		return
	}
	var userID int

	if i, err := strconv.Atoi(querry["id"][0]); err == nil {
		userID = i
	} else {
		log.Printf("id is not correct\n")
		return
	}
	user, ok := idUserDB[userID]
	if !ok {
		fmt.Fprintf(writer, "no user with this id %d", userID)
		log.Printf("no user with this id %d", userID)
	}
	nameHobbyDB[user.Name] = user.Hobby
	fmt.Fprintf(writer, `{"Name":"%s","Hobby":"%s"}`, user.Name, user.Hobby)
	log.Printf(`{"name":"%s","hobby":"%s"}`, user.Name, user.Hobby) //*to do new encoder
}

func getID() int {
	number++
	return number
}
