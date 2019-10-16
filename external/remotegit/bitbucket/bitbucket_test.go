/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package bitbucket

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"golang.org/x/oauth2"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("", http.DefaultClient)
	if err != nil {
		t.Errorf("Unable to create new client: %v", err)
		return
	}

	if clientType := reflect.TypeOf(client); clientType != reflect.TypeOf(&Bitbucket{}) {
		t.Errorf("Returned client (%v) has a different type from expected", clientType)
	}
}

func TestNewOauth2Client(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		_, _ = rw.Write([]byte(`{"access_token":"my-access-token","token_type":"token-type","refresh_token":"my-refresh-token"}`))
	}))
	defer server.Close()

	key := "key"
	secret := "secret"
	username := "username"
	password := "password"
	client, err := NewOauth2Client(key, secret, username, password, http.DefaultClient, &oauth2.Endpoint{TokenURL: server.URL})
	if err != nil {
		t.Errorf("Unable to create new client: %v", err)
		return
	}

	if client.token.AccessToken != "my-access-token" {
		t.Errorf("Want my-access-token, got %v", client.token.AccessToken)
	}
	if client.token.TokenType != "token-type" {
		t.Errorf("Want token-type, got %v", client.token.TokenType)
	}
	if client.token.RefreshToken != "my-refresh-token" {
		t.Errorf("Want my-refresh-token, got %v", client.token.RefreshToken)
	}

	if clientType := reflect.TypeOf(client); clientType != reflect.TypeOf(&Bitbucket{}) {
		t.Errorf("Returned client (%v) has a different type from expected", clientType)
	}
}
