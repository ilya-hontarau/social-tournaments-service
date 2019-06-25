package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/illfate/social-tournaments-service/pkg/sts"

	"github.com/illfate/social-tournaments-service/internal/mockdb"
)

func TestAddUser(t *testing.T) {
	tt := []struct {
		name        string
		method      string
		request     string
		response    string
		status      int
		contentType string
	}{
		{
			name:        "correct test",
			method:      http.MethodPost,
			request:     `{ "name" : "ilya" }`,
			response:    `{"id":1}`,
			status:      http.StatusOK,
			contentType: "application/json",
		},
		{
			name:        "incorrect request",
			method:      http.MethodPost,
			request:     `{  : "i" }`,
			status:      http.StatusBadRequest,
			contentType: "text/plain; charset=utf-8",
		},
		{
			name:    "incorrect method",
			method:  http.MethodPatch,
			request: `{ "name" : "ilya" }`,
			status:  http.StatusMethodNotAllowed,
		},
	}
	db := new(mockdb.Connector)
	db.On("AddUser", "ilya").Return(int64(1), nil)
	s, err := NewServer(db)
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}

	server := httptest.NewServer(s)
	defer server.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, fmt.Sprintf("%s/user", server.URL), strings.NewReader(tc.request))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("couldnt get response: %s", err)
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
			if tc.status != resp.StatusCode {
				t.Fatalf("expected status %v; got %v", tc.status, resp.StatusCode)
			}
			if contentType := resp.Header.Get("Content-Type"); tc.contentType != contentType {
				t.Fatalf("expected status %v; got %v", tc.contentType, contentType)
			}
			if tc.status == http.StatusOK {
				if respBody := string(bytes.TrimSpace(b)); tc.response != respBody {
					t.Fatalf("expected %s, got %s", tc.response, respBody)
				}
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	tt := []struct {
		name        string
		id          string
		response    string
		status      int
		contentType string
	}{
		{
			name:        "correct test",
			id:          "1",
			response:    `{"id":1,"name":"ilya","balance":0}`,
			status:      http.StatusOK,
			contentType: "application/json",
		},
		{
			name:        "incorrect id",
			id:          "ahoi",
			status:      http.StatusNotFound,
			contentType: "text/plain; charset=utf-8",
		},
		{
			name:   "uncreated account",
			id:     "1000",
			status: http.StatusNotFound,
		},
	}
	db := new(mockdb.Connector)
	db.On("GetUser", int64(1)).Return(&sts.User{
		ID:      1,
		Name:    "ilya",
		Balance: 0,
	}, nil)
	db.On("GetUser", int64(1000)).Return(&sts.User{}, sts.ErrNotFound) // TODO : fix return
	s, err := NewServer(db)
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}

	server := httptest.NewServer(s)
	defer server.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/%s", server.URL, tc.id), nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("couldnt get response: %s", err)
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
			if tc.status != resp.StatusCode {
				t.Fatalf("expected status %v; got %v", tc.status, resp.StatusCode)
			}
			if contentType := resp.Header.Get("Content-Type"); tc.contentType != contentType {
				t.Fatalf("expected status %v; got %v", tc.contentType, contentType)
			}
			if tc.status == http.StatusOK {
				if respBody := string(bytes.TrimSpace(b)); tc.response != respBody {
					t.Fatalf("expected %s, got %s", tc.response, respBody)
				}
			}

		})
	}
}

func TestFund(t *testing.T) {
	tt := []struct {
		name        string
		id          string
		request     string
		status      int
		contentType string
	}{
		{
			name:        "correct test",
			id:          "1",
			request:     `{ "points" : 7 }`,
			status:      http.StatusOK,
			contentType: "application/json",
		},
		{
			name:        "empty id",
			request:     `{ "points" :​ 300 }`,
			status:      http.StatusNotFound,
			contentType: "text/plain; charset=utf-8",
		},
		{
			name:        "incorrect json",
			id:          "10",
			request:     `{ "name" : "max" }`,
			status:      http.StatusNotFound,
			contentType: "text/plain; charset=utf-8",
		},
	}
	db := new(mockdb.Connector)
	db.On("AddPoints", int64(1), int64(7)).Return(nil)
	db.On("AddPoints", int64(10), int64(0)).Return(sts.ErrNotFound)
	s, err := NewServer(db)
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}

	server := httptest.NewServer(s)
	defer server.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/%s/fund", server.URL, tc.id),
				strings.NewReader(tc.request))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("couldnt get response: %s", err)
			}
			defer resp.Body.Close()
			if tc.status != resp.StatusCode {
				t.Fatalf("expected status %v; got %v", tc.status, resp.StatusCode)
			}
			if contentType := resp.Header.Get("Content-Type"); tc.contentType != contentType {
				t.Fatalf("expected status %v; got %v", tc.contentType, contentType)
			}
		})
	}
}

func TestTake(t *testing.T) {
	tt := []struct {
		name        string
		id          string
		request     string
		status      int
		contentType string
	}{
		{
			name:        "correct test",
			id:          "1",
			request:     `{ "points" : 7 }`,
			status:      http.StatusOK,
			contentType: "application/json",
		},
		{
			name:        "empty id",
			request:     `{"points":​300}`,
			status:      http.StatusNotFound,
			contentType: "text/plain; charset=utf-8",
		},
		{
			name:        "uncreated account",
			id:          "1000",
			request:     `{ "points" : 7 }`,
			status:      http.StatusNotFound,
			contentType: "text/plain; charset=utf-8",
		},
		{
			name:        "incorrect bonus request",
			id:          "1",
			request:     `{ "points" : 7000 }`,
			status:      http.StatusInternalServerError,
			contentType: "text/plain; charset=utf-8",
		},
	}
	db := new(mockdb.Connector)
	db.On("AddPoints", int64(1), int64(-7)).Return(nil)
	db.On("AddPoints", int64(1000), int64(-7)).Return(sts.ErrNotFound)
	db.On("AddPoints", int64(1), int64(-7000)).Return(sql.ErrNoRows)
	s, err := NewServer(db)
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}

	server := httptest.NewServer(s)
	defer server.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/%s/take", server.URL, tc.id),
				strings.NewReader(tc.request))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("couldnt get response: %s", err)
			}
			defer resp.Body.Close()
			if tc.status != resp.StatusCode {
				t.Fatalf("expected status %v; got %v", tc.status, resp.StatusCode)

			}
			if contentType := resp.Header.Get("Content-Type"); tc.contentType != contentType {
				t.Fatalf("expected status %v; got %v", tc.contentType, contentType)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tt := []struct {
		name        string
		id          string
		status      int
		contentType string
	}{
		{
			name:   "correct test",
			id:     "1",
			status: http.StatusOK,
		},
		{
			name:        "incorrect id",
			id:          "ahoi",
			status:      http.StatusNotFound,
			contentType: "text/plain; charset=utf-8",
		},
		{
			name:   "uncreated user",
			id:     "100",
			status: http.StatusNotFound,
		},
	}
	db := new(mockdb.Connector)
	db.On("DeleteUser", int64(1)).Return(nil)
	db.On("DeleteUser", int64(100)).Return(sts.ErrNotFound)
	s, err := NewServer(db)
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}

	server := httptest.NewServer(s)
	defer server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/user/%s", server.URL, tc.id), nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("couldnt get response: %s", err)
			}
			defer resp.Body.Close()

			if tc.status != resp.StatusCode {
				t.Fatalf("expected status %v; got %v", tc.status, resp.StatusCode)
			}
			if contentType := resp.Header.Get("Content-Type"); tc.contentType != contentType {
				t.Fatalf("expected status %v; got %v", tc.contentType, contentType)
			}
		})
	}
}

func TestAddTournament(t *testing.T) {
	tt := []struct {
		name        string
		method      string
		request     string
		response    string
		status      int
		contentType string
	}{
		{
			name:        "correct test",
			method:      http.MethodPost,
			request:     `{"name": "poker","deposit": 1000}`,
			response:    `{"id":1}`,
			status:      http.StatusOK,
			contentType: "application/json",
		},
		{
			name:        "incorrect request",
			method:      http.MethodPost,
			request:     `{  :  }`,
			status:      http.StatusBadRequest,
			contentType: "text/plain; charset=utf-8",
		},
		{
			name:        "incorrect deposit",
			method:      http.MethodPost,
			request:     `{"name": "football","deposit": -1000}`,
			status:      http.StatusBadRequest,
			contentType: "text/plain; charset=utf-8",
		},
		{
			name:    "incorrect method",
			method:  http.MethodPatch,
			request: `{"name": "football","deposit": 1000}`,
			status:  http.StatusMethodNotAllowed,
		},
	}
	db := new(mockdb.Connector)
	db.On("AddTournament", "poker", uint64(1000)).Return(int64(1), nil)
	s, err := NewServer(db)
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}

	server := httptest.NewServer(s)
	defer server.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, fmt.Sprintf("%s/tournament", server.URL),
				strings.NewReader(tc.request))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("couldnt get response: %s", err)
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
			if tc.status != resp.StatusCode {
				t.Fatalf("expected status %v; got %v", tc.status, resp.StatusCode)
			}
			if contentType := resp.Header.Get("Content-Type"); tc.contentType != contentType {
				t.Fatalf("expected status %v; got %v", tc.contentType, contentType)
			}
			if tc.status == http.StatusOK {
				if respBody := string(bytes.TrimSpace(b)); tc.response != respBody {
					t.Fatalf("expected %s, got %s", tc.response, respBody)
				}
			}
		})
	}
}

func TestGetTournament(t *testing.T) {
	tt := []struct {
		name        string
		id          string
		response    string
		status      int
		contentType string
	}{
		{
			name:        "correct test",
			id:          "1",
			response:    `{"id":1,"name":"poker","deposit":1000,"prize":0,"winner":0,"users":[2,3]}`,
			status:      http.StatusOK,
			contentType: "application/json",
		},
		{
			name:        "incorrect id",
			id:          "ahoi",
			status:      http.StatusNotFound,
			contentType: "text/plain; charset=utf-8",
		},
		{
			name:        "uncreated account",
			id:          "1000",
			status:      http.StatusNotFound,
			contentType: "text/plain; charset=utf-8",
		},
	}
	db := new(mockdb.Connector)
	db.On("GetTournament", int64(1)).Return(&sts.Tournament{
		ID:      1,
		Name:    "poker",
		Deposit: 1000,
		Prize:   0,
		Winner:  0,
		Users:   []int64{2, 3},
	}, nil)
	db.On("GetTournament", int64(1000)).Return(&sts.Tournament{}, sts.ErrNotFound) // TODO: fix return
	s, err := NewServer(db)
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}
	server := httptest.NewServer(s)
	defer server.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/tournament/%s", server.URL, tc.id), nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("couldnt get response: %s", err)
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
			if tc.status != resp.StatusCode {
				t.Fatalf("expected status %v; got %v", tc.status, resp.StatusCode)
			}
			if contentType := resp.Header.Get("Content-Type"); tc.contentType != contentType {
				t.Fatalf("expected status %v; got %v", tc.contentType, contentType)
			}
			if tc.status == http.StatusOK {
				if respBody := string(bytes.TrimSpace(b)); tc.response != respBody {
					t.Fatalf("expected %s, got %s", tc.response, respBody)
				}
			}

		})
	}
}
