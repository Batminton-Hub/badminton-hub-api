package util

func GenUserID(email, password string) string {
	var id string
	id = Sha256(email + password)
	return id
}
