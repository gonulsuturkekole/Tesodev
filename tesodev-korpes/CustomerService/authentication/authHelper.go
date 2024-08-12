package authentication

import (
	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
	_ "time"
)

var SecretKey = []byte("secret")

type Claims struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	jwt.StandardClaims
}

func JwtGenerator(Id string, firstName string, lastName string, key string) string {
	//Generate Token JWT for auth
	claims := &Claims{
		ID:        Id,
		FirstName: firstName,
		LastName:  lastName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return err.Error()
	}
	return tokenString
}

// VerifyJWT checks if the token is a valid JWT
func VerifyJWT(tokenString string) error {
	_, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return echo.ErrUnauthorized
	}

	return nil
}

// HashPassword hashes the given password using bcrypt and returns the hashed password as a string.
//func HashPassword(password string) (string, error) {
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	if err != nil {
//		return "", err
//	}
//	return string(hashedPassword), nil
//}
//
//// // CheckPasswordHash checks whether the given password matches the hashed password.
//func CheckPasswordHash(password, hashedPassword string) bool {
//	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
//	return err == nil
//}

// HashPassword hashes the given password using bcrypt and the secret key,
// returning the hashed password as a string.
func HashPassword(password string) (string, error) {
	// Combine the secret key and password bytes
	combined := append(SecretKey, []byte(password)...)
	// Generate the bcrypt hash of the combined byte slice
	hashedPassword, err := bcrypt.GenerateFromPassword(combined, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash checks if the provided password matches the hashed password
// by including the secret key in the comparison process.
func CheckPasswordHash(password, hashedPassword string) bool {
	// Combine the secret key and password bytes
	combined := append(SecretKey, []byte(password)...)
	// Compare the hashed password with the combined byte slice
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), combined)
	return err == nil
}
