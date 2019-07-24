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

func NewResolver(db sts.Service, _ string) (*Resolver, error) {
	toutnBytes, err := ioutil.ReadFile("sts.graphql")
	if err != nil {
		return nil, errors.Wrap(err, "foobar")
	}
	userBytes, err := ioutil.ReadFile("sts-user.graphql")
	if err != nil {
		return nil, errors.Wrap(err, "foobar")
	}

	var resolver Resolver
	tSchema, err := graphql.ParseSchema(string(toutnBytes), &resolver)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't parse schema")
	}
	uSchema, err := graphql.ParseSchema(string(userBytes), &resolver)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't parse schema")
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
