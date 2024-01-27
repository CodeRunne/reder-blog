package utility

import (
	"fmt"
	"errors"
	"log"
	"os"
	"time"
	"strconv"
	"regexp"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

var (
	ErrInvalidEmail  = errors.New("Email Address Is Invalid!")
)

func init() {
	// Load env
	if err := Load("./.env"); err != nil {
		log.Fatal(err.Error())
	}
}

func Load(path string) error {
	if err := godotenv.Load(path); err != nil {
		return err
	}
	return nil
}

func CreateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"issuer": email,
		"expiry": time.Now().Add((7 * 24 * 60 * 60) * time.Hour).Unix(),
	})

	// Get secret key from env
	key := os.Getenv("SECRET_KEY")

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (map[string]interface{}, error) {
	var claims map[string]interface{}

	// Get secret key from env
	key := []byte(os.Getenv("SECRET_KEY"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return claims, err
	}

	// Get Claims
	var ok bool
	if claims, ok = token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return claims, fmt.Errorf("Invalid token")
	}

	return claims, nil
}

func ValidateEmail(email string) error {
	match, _ := regexp.MatchString(`[a-zA-Z0-9_]+@[a-z0-9]{4,5}\.(com|org)`, email)
	if match == false {
		return ErrInvalidEmail
	}
	return nil
}

func SendMail(to, subject, body string) error {

	// Get Mail Addtional Information From Environment
	from := os.Getenv("MAIL_FROM")
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	_ = os.Getenv("MAIL_MAILER")

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer(host, port, username, password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil

}
