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
	s, err := NewServer()
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}
	defer s.DB.Close()

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
			if contentType := strings.Join(resp.Header["Content-Type"], ""); tc.contentType != contentType {
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
	s, err := NewServer()
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}
	defer s.DB.Close()

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
			if contentType := strings.Join(resp.Header["Content-Type"], ""); tc.contentType != contentType {
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
			request:     `{ "name" : "max" }`,
			status:      http.StatusNotFound,
			contentType: "text/plain; charset=utf-8",
		},
	}
	s, err := NewServer()
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}
	defer s.DB.Close()

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
			if contentType := strings.Join(resp.Header["Content-Type"], ""); tc.contentType != contentType {
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
			name:    "uncreated account",
			id:      "1000",
			request: `{ "points" : 7 }`,
			status:  http.StatusNotFound,
		},
		{
			name:        "incorrect bonus request",
			id:          "1",
			request:     `{ "points" : 7000 }`,
			status:      http.StatusInternalServerError,
			contentType: "text/plain; charset=utf-8",
		},
	}
	s, err := NewServer()
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}
	defer s.DB.Close()

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
			if contentType := strings.Join(resp.Header["Content-Type"], ""); tc.contentType != contentType {
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
	s, err := NewServer()
	if err != nil {
		t.Fatalf("couldn't create db connection: %v", err)
	}
	defer s.DB.Close()

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
			if contentType := strings.Join(resp.Header["Content-Type"], ""); tc.contentType != contentType {
				t.Fatalf("expected status %v; got %v", tc.contentType, contentType)
			}
		})
	}
}
