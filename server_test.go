package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddUser(t *testing.T) {
	tt := []struct {
		name     string
		method   string
		request  string
		response string
		status   int
	}{
		{
			name:     "correct test",
			method:   http.MethodPost,
			request:  `{ "name" : "ilya" }`,
			response: `{"id":1}`,
			status:   http.StatusOK,
		},
		{
			name:    "incorrect request",
			method:  http.MethodPost,
			request: `{  : "i" }`,
			status:  http.StatusBadRequest,
		},
		{
			name:    "incorrect method",
			method:  http.MethodPatch,
			request: `{ "name" : "ilya" }`,
			status:  http.StatusNotFound,
		},
	}
	s, err := NewServer()
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}
	defer s.DB.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, "localhost:9000/user", strings.NewReader(tc.request))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			s.addUser(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
			if tc.status != res.StatusCode {
				t.Fatalf("expected status %v; got %v", tc.status, res.StatusCode)
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
		name     string
		id       string
		response string
		status   int
	}{
		{
			name:     "correct test",
			id:       "1",
			response: `{"id":1,"name":"ilya","balance":0}`,
			status:   http.StatusOK,
		},
		{
			name:   "incorrect id",
			id:     "ahoi",
			status: http.StatusBadRequest,
		},
		{
			name:   "uncreated account",
			id:     "1000",
			status: http.StatusNotFound,
		},
	}
	s, err := NewServer()
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}
	defer s.DB.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "http://localhost:9000/user/"+tc.id, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			s.getUser(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
			if tc.status != res.StatusCode {
				t.Fatalf("expected status %v; got %v", tc.status, res.StatusCode)
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
		name    string
		id      string
		request string
		status  int
	}{
		{
			name:    "correct test",
			id:      "1",
			request: `{ "points" : 7 }`,
			status:  http.StatusOK},
		{
			name:    "empty id",
			request: `{ "points" :​ 300 }`,
			status:  http.StatusBadRequest,
		},
		{
			name:    "incorrect json",
			request: `{ "name" : "max" }`,
			status:  http.StatusBadRequest,
		},
	}
	s, err := NewServer()
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}
	defer s.DB.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:9000/user/%s/fund", tc.id),
				strings.NewReader(tc.request))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			s.processBonus(rec, req)
			res := rec.Result()
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
			if tc.status != res.StatusCode {
				t.Fatalf("expected status %v; got %v", tc.status, res.StatusCode)
			}
		})
	}
}

func TestTake(t *testing.T) {
	tt := []struct {
		name    string
		id      string
		request string
		status  int
	}{
		{
			name:    "correct test",
			id:      "1",
			request: `{ "points" : 7 }`,
			status:  http.StatusOK,
		},
		{
			name:    "empty id",
			request: `{"points":​300}`,
			status:  http.StatusBadRequest,
		},
		{
			name:    "uncreated account",
			id:      "1000",
			request: `{ "points" : 7 }`,
			status:  http.StatusNotFound,
		},
		{
			name:    "incorrect bonus request",
			id:      "1",
			request: `{ "points" : 7000 }`,
			status:  http.StatusInternalServerError,
		},
	}
	s, err := NewServer()
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}
	defer s.DB.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:9000/user/%s/take", tc.id),
				strings.NewReader(tc.request))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			s.processBonus(rec, req)
			res := rec.Result()
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
			if tc.status != res.StatusCode {
				t.Fatalf("expected status %v; got %v", tc.status, res.StatusCode)

			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tt := []struct {
		name   string
		id     string
		status int
	}{
		{
			name:   "correct test",
			id:     "1",
			status: http.StatusOK,
		},
		{
			name:   "incorrect id",
			id:     "ahoi",
			status: http.StatusBadRequest,
		},
		{
			name:   "uncreated user",
			id:     "100",
			status: http.StatusNotFound,
		},
	}
	s, err := NewServer()
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}
	defer s.DB.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "http://localhost:9000/user/"+tc.id, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			s.deleteUser(rec, req)
			res := rec.Result()
			if tc.status != res.StatusCode {
				t.Errorf("expected status %v; got %v", tc.status, res.StatusCode)
				return
			}

		})
	}
}

func TestHandler(t *testing.T) {
	tt := []struct {
		name    string
		method  string
		url     string
		request string
		status  int
	}{
		{
			name:   "empty url",
			method: "POST",
			status: http.StatusNotFound,
		},
		{
			name:   "incorrect method - process bonus",
			url:    "/user/1/fund",
			method: http.MethodDelete,
			status: http.StatusNotFound,
		},
		{
			name:    "incorrect request - process bonus",
			url:     "/user/1/fund",
			method:  http.MethodPost,
			request: "empty req",
			status:  http.StatusBadRequest,
		},
		{
			name:    "incorrect url - process bonus",
			url:     "/user/666/fud",
			method:  http.MethodPost,
			request: `{ "points" : 7 }`,
			status:  http.StatusNotFound,
		},
		{
			name:   "get user incorrect id",
			url:    "/user/ahoi",
			method: http.MethodGet,
			status: http.StatusBadRequest,
		},
		{
			name:   "delete user incorrect id",
			url:    "/user/ahoi",
			method: http.MethodDelete,
			status: http.StatusBadRequest,
		},
	}
	s, err := NewServer()
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}
	defer s.DB.Close()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, "http://localhost:9000"+tc.url,
				strings.NewReader(tc.request))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			s.switchHandler(rec, req)
			res := rec.Result()
			if tc.status != res.StatusCode {
				t.Errorf("expected status %v; got %v", tc.status, res.StatusCode)
				return
			}
		})
	}
}
