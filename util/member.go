package util

func GenUserID(email, password string) string {
	return Sha256(email + password)
}
