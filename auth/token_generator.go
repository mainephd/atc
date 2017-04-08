package auth

import (
	"crypto/rsa"
	"crypto/sha512"
	"time"

	"encoding/base64"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

//go:generate counterfeiter . TokenGenerator

type TokenType string
type TokenValue string

const TokenTypeBearer = "Bearer"
const TokenTypeAccess = "Access"
const expClaimKey = "exp"
const teamNameClaimKey = "teamName"
const teamIDClaimKey = "teamID"
const isAdminClaimKey = "isAdmin"

type TokenGenerator interface {
	GenerateToken(expiration time.Time, teamName string, teamID int, isAdmin bool) (TokenType, TokenValue, error)
	GenerateAccessToken(teamName string, teamID int, isAdmin bool) (TokenType, TokenValue, error)
}

type tokenGenerator struct {
	privateKey *rsa.PrivateKey
}

func NewTokenGenerator(privateKey *rsa.PrivateKey) TokenGenerator {
	return &tokenGenerator{
		privateKey: privateKey,
	}
}

func (generator *tokenGenerator) GenerateToken(expiration time.Time, teamName string, teamID int, isAdmin bool) (TokenType, TokenValue, error) {
	jwtToken := jwt.NewWithClaims(SigningMethod, jwt.MapClaims{
		expClaimKey:      expiration.Unix(),
		teamNameClaimKey: teamName,
		teamIDClaimKey:   teamID,
		isAdminClaimKey:  isAdmin,
	})

	signed, err := jwtToken.SignedString(generator.privateKey)
	if err != nil {
		return "", "", err
	}

	return TokenTypeBearer, TokenValue(signed), err
}

func (generator *tokenGenerator) GenerateAccessToken(teamName string, teamID int, isAdmin bool) (TokenType, TokenValue, error) {
	hasher := sha512.New()
	hasher.Write(uuid.NewV1().Bytes())
	return TokenTypeAccess, TokenValue(base64.URLEncoding.EncodeToString(hasher.Sum(nil))), nil
}
