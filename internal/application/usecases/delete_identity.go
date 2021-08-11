package usecases

import "github.com/hywmongous/example-service/internal/domain/identity"

type DeletionRequest struct {
	IdentityID identity.IdentityID
}

type DeletionResponse struct {
}

type Delete interface {
	DoDelete(request DeletionRequest) (DeletionResponse, error)
}
