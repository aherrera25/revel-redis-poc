package tests

import (
	sessionHandler "clock_project_v2/app/actors/redis"
	"context"

	"github.com/revel/revel/session"
	"github.com/revel/revel/testing"
)

type SessionActorTest struct {
	testing.TestSuite
}

type SessionMock struct {
	session.Session
}

func (module SessionMock) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return "", nil
}

func (t *SessionActorTest) TestsSettingsController() {
	response, status := sessionHandler.Build(SessionMock{}, context.Background(), "GET")("message")
	t.AssertEqual(status, 200)
	t.AssertEqual(response, "")
}

//TODO: err scenario
