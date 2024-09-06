package security

type TokenHandler interface {
	GenerateToken(email string, username string) (string, error)
}
