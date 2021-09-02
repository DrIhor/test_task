package jwt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	er "github.com/DrIhor/test_task/internal/errors"
	"github.com/golang-jwt/jwt"
)

const googleCertURL string = "https://www.googleapis.com/oauth2/v1/certs"

func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get(googleCertURL)
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", er.JWTKeyNotFound
	}
	return key, nil
}

func ValidateGoogleExpireJWT(tokenString string) (bool, error) {
	claimsStruct := GoogleModel{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				return nil, err
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err
			}
			return key, nil
		},
	)
	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(*GoogleModel)
	if !ok {
		return false, er.InvalidTocken
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return false, er.TockenExpired
	}

	fmt.Println("FirstName: ", claims.FirstName)
	fmt.Println("LastName: ", claims.LastName)

	return true, nil
}
