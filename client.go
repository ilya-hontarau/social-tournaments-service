package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// User represents user with name and hobby strings fields.
type User struct {
	Name  string `json:"name"`
	Hobby string `json:"hobby"`
}

const baseURL = "http://localhost:9000"

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s\n-g <name> <your hobby>\n-p <id>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	if os.Args[1] == "-p" && len(os.Args) > 3 {
		result, err := postUser(os.Args[2], os.Args[3:])
		if err != nil {
			log.Fatal("could not add user, ", err)
		}
		fmt.Println(result)
	} else if os.Args[1] == "-g" && len(os.Args) == 3 {
		result, err := getUser(os.Args[2])
		if err != nil {
			log.Fatal("could not get user, ", err)
		}
		fmt.Println(result)
	} else {
		log.Printf("incorrect request, usage: %s\n", filepath.Base(os.Args[0]))
	}
}

func postUser(name string, hobby []string) (string, error) {
	user := User{name, strings.Join(hobby, " ")}
	userData, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(userData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getUser(id string) (string, error) {
	if _, err := strconv.Atoi(id); err != nil {
		return "", errors.New("incorrect id")
	}
	resp, err := http.Get(fmt.Sprintf("%s?id=%s", baseURL, id))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return "", errors.New("id doesn't exist")
		}
		return "", err
	}
	return fmt.Sprintf("%s,your hobby is %s", user.Name, user.Hobby), nil
}
