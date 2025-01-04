package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHash(t *testing.T) {

	pwd := "this is my test password"

	hashed_password, err := HashPassword(pwd)
	if err != nil {
		t.Errorf("Error hashing password")
	}

	err = CheckPassword(pwd, hashed_password)
	if err != nil {
		t.Errorf("Hashed password does not match")
	}

}

func TestJWT(t *testing.T) {

	userID := uuid.New()
	tokenSecret := "kljaskldjkjkas9897234ldkjasd"
	expiresIn := time.Second

	tokenSigned, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("Big 'ol error here:\n %v", err)
		return
	}

	userIDmaybe, err := ValidateJWT(tokenSigned, tokenSecret)
	if err != nil {
		t.Errorf("Error in validation:\n %v", err)
		return
	}
	if userID != userIDmaybe {
		t.Errorf("Retrieved UserID does not match!")
		return
	}

	_, err = ValidateJWT(tokenSigned, "jfdksjkldkfj")
	if err == nil {
		t.Errorf("JWT validated with wrong secret key")
	}

	time.Sleep(2 * expiresIn)

	if _, err := ValidateJWT(tokenSigned, tokenSecret); err == nil {
		t.Errorf("Token did not expire correctly")
	}

}

func TestGetBearerToken(t *testing.T) {

	header := http.Header{}
	header.Set("Authorization", "Thisisthebearer Thisisthetokenstring")
	header_noauth := http.Header{}
	str, err := GetBearerToken(header)
	if err != nil {
		t.Errorf("%s", err)
	}
	if str != "Thisisthetokenstring" {
		t.Errorf("token_string not retrieved successfully")
	}
	str, err = GetBearerToken(header_noauth)
	if str != "" || err == nil {
		t.Errorf("Should not have found authorization")
	}

}
