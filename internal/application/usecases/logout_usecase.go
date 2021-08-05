package usecases

import "github.com/hywmongous/example-service/internal/domain/identity/values"

type LogoutRequest struct {
	identityID values.IdentityID
	password   string
}

type LogoutResponse struct {
}

type Logout interface {
	DoLogout(request LogoutRequest) (LogoutResponse, error)
}
