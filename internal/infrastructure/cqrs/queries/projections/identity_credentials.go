package projections

import identity "github.com/hywmongous/example-service/internal/domain/identity/values"

type IdentityCredentials struct {
	IdentityId identity.IdentityID
	Email      identity.Email
}

func (credentials IdentityCredentials) RecreateIdentityCredentials(
	identityId identity.IdentityID,
	email identity.Email,
) IdentityCredentials {
	return IdentityCredentials{
		IdentityId: identityId,
		Email:      email,
	}
}
