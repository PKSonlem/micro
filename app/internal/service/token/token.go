package token

type TokenService struct {
	key string
}

func NewTokenService(key string) *TokenService {
	return &TokenService{key: key}
}

func (t *TokenService) GenerateToken(userId int64) string {
	return "not implemented"
}
