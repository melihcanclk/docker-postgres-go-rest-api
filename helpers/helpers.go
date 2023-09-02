package helpers

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/melihcanclk/docker-postgres-go-rest-api/models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/unicode/norm"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func isIncludesNonAscii(input string) error {
	for _, r := range input {
		if r > unicode.MaxASCII {
			return fmt.Errorf("input contains non-ascii characters")
		}
	}
	return nil
}

func IsIncludesNonAscii(input *string) error {
	normalized := norm.NFKD.String(*input)
	return isIncludesNonAscii(normalized)
}

func ValidateToken(token string, publicKey string) (*models.TokenDetails, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return &models.TokenDetails{
		TokenUuid: fmt.Sprint(claims["token_uuid"]),
		UserID:    fmt.Sprint(claims["sub"]),
	}, nil
}
func GenerateJWTToken(id string, ttl *time.Duration, privateKey string) (*models.TokenDetails, error) {
	now := time.Now().UTC()
	td := &models.TokenDetails{
		ExpiresIn: new(int64),
		Token:     new(string),
	}
	*td.ExpiresIn = now.Add(*ttl).Unix()
	fmt.Println(*ttl)
	td.TokenUuid = uuid.New().String()
	td.UserID = id

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode token private key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("create: parse token private key: %w", err)
	}

	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = td.UserID
	atClaims["token_uuid"] = td.TokenUuid
	atClaims["exp"] = td.ExpiresIn
	atClaims["iat"] = now.Unix()
	atClaims["nbf"] = now.Unix()

	*td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("create: sign token: %w", err)
	}

	return td, nil
}
