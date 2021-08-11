package visitors

import (
	identity "github.com/hywmongous/example-service/internal/domain/identity"
	"github.com/hywmongous/example-service/internal/infrastructure/cqrs/queries/projections"
)

type IdentityCredentialsVisitor struct {
	credentials projections.IdentityCredentials
}

func CreateIdentityCredentialsVisitor() (IdentityCredentialsVisitor, error) {
	return IdentityCredentialsVisitor{
		credentials: projections.IdentityCredentials{},
	}, nil
}

func RecreateIdentityCredentialsVisitor(
	credentials projections.IdentityCredentials,
) (IdentityCredentialsVisitor, error) {
	return IdentityCredentialsVisitor{
		credentials: credentials,
	}, nil
}

func (visitor IdentityCredentialsVisitor) VisitRegisterIdentity(command identity.Registered) error {
	visitor.credentials.IdentityId = command.GetId()
	visitor.credentials.Email = identity.RecreateEmail(command.GetEmail(), false)
	return nil
}
