package identity

type Registered struct {
	Id             IdentityID
	Time           Time
	Email          string
	HashedPassword string
}
