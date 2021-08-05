package projections

import "github.com/hywmongous/example-service/internal/domain/identity/values"

type IdentityCredentials struct {
	IdentityId values.IdentityID
	Email      values.Email
}

func (credentials IdentityCredentials) RecreateIdentityCredentials(
	identityId values.IdentityID,
	email values.Email,
) IdentityCredentials {
	return IdentityCredentials{
		IdentityId: identityId,
		Email:      email,
	}
}
