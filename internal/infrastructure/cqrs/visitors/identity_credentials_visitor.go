package visitors

import (
	"github.com/hywmongous/example-service/internal/domain/identity/values"
	"github.com/hywmongous/example-service/internal/infrastructure/cqrs/commands"
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

func (visitor IdentityCredentialsVisitor) Visit(commands []commands.Command) error {
	for _, command := range commands {
		if err := command.Apply(visitor); err != nil {
			return err
		}
	}
	return nil
}

func (visitor IdentityCredentialsVisitor) VisitRegisterIdentity(command commands.RegisterIdentity) error {
	visitor.credentials.IdentityId = command.IdentityID
	visitor.credentials.Email = values.RecreateEmail(command.EmailAddress, false)
	return nil
}

func (visitor IdentityCredentialsVisitor) VisitDeleteIdentity(deletion commands.DeleteIdentity) error {
	return nil
}

func (visitor IdentityCredentialsVisitor) VisitIdentityLogin(login commands.IdentityLogin) error {
	return nil
}

func (visitor IdentityCredentialsVisitor) VisitIdentityLogout(logout commands.IdentityLogout) error {
	return nil
}
