package identity

type Email struct {
	address string

	confirmed bool
}

func CreateEmail(address string) (Email, error) {
	return Email{
		address:   address,
		confirmed: false,
	}, nil
}

func RecreateEmail(
	address string,
	confirmed bool,
) Email {
	return Email{
		address:   address,
		confirmed: confirmed,
	}
}
