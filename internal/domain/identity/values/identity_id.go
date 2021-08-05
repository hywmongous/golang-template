package values

type IdentityID uuidValue

func GenerateIdentityID() IdentityID {
	return IdentityID(generateUuidValue())
}

func CreateIdentityID(value string) (IdentityID, error) {
	uuid, err := createUuidValue(value)
	if err != nil {
		return IdentityID(""), err
	}
	return IdentityID(uuid), nil
}

func RecreateIdentityID(value string) IdentityID {
	return IdentityID(value)
}

func (id IdentityID) ToString() string {
	return uuidValue(id).toString()
}

func (id IdentityID) Equals(other IdentityID) bool {
	return uuidValue(id).equals(uuidValue(other))
}
