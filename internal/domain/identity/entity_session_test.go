package identity

import (
	"testing"
)

func TestCreateSession(t *testing.T) {
	var session Session
	var err error

	session, err = CreateSession()

	if err != nil {
		t.Error("CreateSession failed with err:", err)
	}
	if session.id == "" {
		t.Error("CreateSession created a session with an empty id")
	}
	if session.revoked {
		t.Error("CreateSession created a session which was immediately revoked")
	}
}
