package usecases

import "github.com/hywmongous/example-service/internal/domain/identity/values"

type DeletionRequest struct {
	IdentityID values.IdentityID
}

type DeletionResponse struct {
}

type Delete interface {
	DoDelete(request DeletionRequest) (DeletionResponse, error)
}
