package commands

type CommandHandler interface {
	VisitRegisterIdentity(registration RegisterIdentity) error
	VisitDeleteIdentity(deletion DeleteIdentity) error
	VisitIdentityLogin(login IdentityLogin) error
	VisitIdentityLogout(logout IdentityLogout) error
}
