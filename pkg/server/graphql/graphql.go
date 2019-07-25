package graphql

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/illfate/social-tournaments-service/pkg/sts"
	"github.com/pkg/errors"
)

type Resolver struct {
	s sts.Service
	http.Handler
}

func NewResolver(db sts.Service, userSchemeFile, tournamentSchemeFile string) (*Resolver, error) {
	tournamentBytes, err := ioutil.ReadFile(tournamentSchemeFile)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read tournament graphql schema")
	}
	userBytes, err := ioutil.ReadFile(userSchemeFile)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read user graphql schema")
	}

	var resolver Resolver
	tSchema, err := graphql.ParseSchema(string(tournamentBytes), &resolver)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't parse tournament schema")
	}
	uSchema, err := graphql.ParseSchema(string(userBytes), &resolver)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't parse user schema")
	}

	t := func(w http.ResponseWriter, req *http.Request) {
		tHandler := relay.Handler{
			Schema: tSchema,
		}
		tHandler.ServeHTTP(w, req)
	}
	u := func(w http.ResponseWriter, req *http.Request) {
		uHandler := relay.Handler{
			Schema: uSchema,
		}
		uHandler.ServeHTTP(w, req)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/tournament", t)
	mux.HandleFunc("/user", u)
	resolver.Handler = mux
	resolver.s = db
	return &resolver, nil

}

func decodeID(id graphql.ID) (int64, error) {
	intID, err := strconv.ParseInt(string(id), 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "couldn't parse id [%s]", id)
	}
	if intID < 1 {
		return 0, errors.Errorf("invalid ID: %d", intID)
	}
	return intID, nil
}

func encodeID(id int64) graphql.ID {
	return graphql.ID(strconv.FormatInt(id, 10))
}
