package usecases

import "github.com/hywmongous/example-service/internal/domain/identity"

type DeleteIdentityRequest struct {
	IdentityID identity.IdentityID
}

type DeleteIdentityResponse struct {
}

type DeleteIdentityUseCase func(request DeleteIdentityRequest) (DeleteIdentityResponse, error)
