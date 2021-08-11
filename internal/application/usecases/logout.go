package usecases

import "github.com/hywmongous/example-service/internal/domain/identity"

type LogoutRequest struct {
	identityID identity.IdentityID
	password   string
}

type LogoutResponse struct {
}

type Logout interface {
	DoLogout(request LogoutRequest) (LogoutResponse, error)
}
