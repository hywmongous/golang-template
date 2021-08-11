package identity

type Registered struct {
	id             IdentityID
	time           Time
	email          string
	hashedPassword string
}

func (registered Registered) GetId() IdentityID {
	return registered.id
}

func (registered Registered) GetTime() Time {
	return registered.time
}

func (registered Registered) GetEmail() string {
	return registered.email
}

func (registered Registered) GetHashedPassword() string {
	return registered.hashedPassword
}
