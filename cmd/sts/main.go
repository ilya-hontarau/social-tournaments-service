package main

import (
	"log"
	"net/http"
	"os"

	"github.com/illfate/social-tournaments-service/pkg/server/graphql"

	"github.com/illfate/social-tournaments-service/pkg/psql"
)

const (
	port                 = "PORT"
	userEnvVar           = "DB_USER"
	passEnvVar           = "DB_PASS"
	dbNameEnvVar         = "DB_NAME"
	userSchemeFile       = "USER_SCHEME_FILE"
	tournamentSchemeFile = "TOURNAMENT_SCHEME_FILE"
)

func main() {
	uScheme := os.Getenv(userSchemeFile)
	if uScheme == "" {
		log.Printf(`no "%s" env variable`, userSchemeFile)
		return
	}
	tScheme := os.Getenv(tournamentSchemeFile)
	if tScheme == "" {
		log.Printf(`no "%s" env variable`, tournamentSchemeFile)
		return
	}
	portNum := os.Getenv(port)
	if portNum == "" {
		log.Printf(`no "%s" env variable`, port)
		return
	}
	dbUser := os.Getenv(userEnvVar)
	if dbUser == "" {
		log.Printf(`no "%s" env variable`, userEnvVar)
		return
	}
	dbPass := os.Getenv(passEnvVar)
	dbName := os.Getenv(dbNameEnvVar)
	if dbName == "" {
		log.Printf(`no "%s" env variable`, dbNameEnvVar)
		return
	}

	db, err := psql.New(dbUser, dbPass, dbName)
	if err != nil {
		log.Print(err)
		return
	}
	defer db.Close()

	s, err := graphql.NewResolver(db, "user.graphql", "tournament.graphql")
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
