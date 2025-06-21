package service

import (
	"time"
	"user-service/config"

	"github.com/golang-jwt/jwt/v5"
)

type JwtServiceInterface interface {
	GenerateToken(userID int64) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey string
	issuer    string
}

func (s *jwtService) GenerateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["iss"] = s.issuer
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(s.secretKey), nil
	})
}

func NewJwtService(cfg *config.Config) JwtServiceInterface {
	return &jwtService{
		secretKey: cfg.App.JwtSecretKey,
		issuer:    cfg.App.JwtIssuer,
	}
}
