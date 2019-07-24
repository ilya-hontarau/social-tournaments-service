package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/illfate/social-tournaments-service/pkg/server/graphql"

	"github.com/illfate/social-tournaments-service/pkg/psql"
)

const (
	port         = "PORT"
	userEnvVar   = "DB_USER"
	passEnvVar   = "DB_PASS"
	dbNameEnvVar = "DB_NAME"
)

func main() {
	portNum := os.Getenv(port)
	if portNum == "" {
		log.Printf(`no "%s" env variable`, port)
		return
	}
	dbUser := os.Getenv(userEnvVar)
	if dbUser == "" {
		log.Print(fmt.Errorf(`no "%s" env variable`, userEnvVar))
		return
	}
	dbPass := os.Getenv(passEnvVar)
	dbName := os.Getenv(dbNameEnvVar)
	if dbName == "" {
		log.Print(fmt.Errorf(`no "%s" env variable`, dbNameEnvVar))
		return
	}

	db, err := psql.New(dbUser, dbPass, dbName)
	if err != nil {
		log.Print(err)
		return
	}
	defer db.Close()

	// todo app flag or env
	s, err := graphql.NewResolver(db, "sts.graphql")
	if err != nil {
		log.Printf("couldn't start graphql: %s", err)
		return
	}

	err = http.ListenAndServe(":"+portNum, s)
	if err != nil {
		log.Print(err)
		return
	}
}
