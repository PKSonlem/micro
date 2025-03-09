package token

type TokenService struct {
	key []byte
}

func NewTokenService(key []byte) *TokenService {
	return &TokenService{key: key}
}

func (t *TokenService) GenerateToken(userId int64) string {
	return "not implemented"
}
