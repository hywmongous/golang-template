package visitors

import (
	"github.com/hywmongous/example-service/internal/domain/identity/aggregate"
	"github.com/hywmongous/example-service/internal/domain/identity/values"
	"github.com/hywmongous/example-service/internal/infrastructure/cqrs/commands"
	merr "github.com/hywmongous/example-service/pkg/errors"
)

type AggregateVisitor struct {
	root aggregate.Identity
}

func CreateAggregateVisitor() (AggregateVisitor, error) {
	return AggregateVisitor{
		root: aggregate.Identity{},
	}, nil
}

func RecreateAggregateVisitor(root aggregate.Identity) (AggregateVisitor, error) {
	return AggregateVisitor{
		root: root,
	}, nil
}

func (visitor AggregateVisitor) Visit(commands []commands.Command) error {
	for _, command := range commands {
		if err := command.Apply(visitor); err != nil {
			return err
		}
	}
	return nil
}

func (visitor AggregateVisitor) VisitRegisterIdentity(command commands.RegisterIdentity) error {
	password, err := values.CreatePassword(command.Password)
	if err != nil {
		return merr.CreateFailedStructInvocation("AggregateVisitor", "VisitRegisterIdentity", err)
	}
	visitor.root = aggregate.RecreateIdentity(
		command.IdentityID,
		values.RecreateEmail(command.EmailAddress, false),
		password,
		[]aggregate.Session{},
		[]aggregate.Scope{},
	)
	return nil
}

func (visitor AggregateVisitor) VisitDeleteIdentity(deletion commands.DeleteIdentity) error {
	return nil
}

func (visitor AggregateVisitor) VisitIdentityLogin(login commands.IdentityLogin) error {
	if _, err := visitor.root.Login(login.Password); err != nil {
		return merr.CreateFailedStructInvocation("AggregateVisitor", "VisitIdentityLogin", err)
	}
	return nil
}

func (visitor AggregateVisitor) VisitIdentityLogout(logout commands.IdentityLogout) error {
	if err := visitor.root.Logout(logout.SessionID); err != nil {
		return merr.CreateFailedStructInvocation("AggregateVisitor", "VisitIdentityLogout", err)
	}
	return nil
}
