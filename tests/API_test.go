package main

import (
	"net/url"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
)

const (
	host = "test-task-app:6060"
)

func TestApi_HappyPath(t *testing.T) {
	time.Sleep(1 * time.Second)
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.GET("/tokens/{user_id}").
		WithPath("user_id", "1716daab-5868-477e-9f51-0df2a0e925b7").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("RefreshToken")
}
