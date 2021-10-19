package authentication_test

import (
	"testing"

	"github.com/hywmongous/example-service/internal/domain/authentication"
)

func TestCreateSession(t *testing.T) {
	t.Parallel()

	_, err := authentication.CreateSession()

	if err != nil {
		t.Error("CreateSession failed with err:", err)
	}
}
