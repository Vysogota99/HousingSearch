package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

// CreateToken ...
func CreateToken(steamID string) (*TokenDetails, error) {
	tokenDet := &TokenDetails{}
	tokenDet.AccTExpires = time.Now().Add(time.Minute * 15).Unix()
	tokenDet.AccessUUID = uuid.NewV4().String()

	tokenDet.RefTExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokenDet.RefreshUUID = uuid.NewV4().String()

	var err error
	// Creating access token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = tokenDet.AccessUUID
	atClaims["steam_id"] = steamID
	atClaims["exp"] = tokenDet.AccTExpires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	tokenDet.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// Creating refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = tokenDet.RefreshUUID
	rtClaims["steam_id"] = steamID
	rtClaims["exp"] = tokenDet.RefTExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	tokenDet.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return tokenDet, nil
}

// VerifyToken - verify the token
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// ExtractToken - extract the token from the request header
func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// TokenValid - check the validity of this token
func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

// ExtractTokenMetadata - extract the token metadata that will lookup in our redis store
func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, ok := claims["steam_id"].(string)
		if !ok {
			return nil, err
		}
		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}
	return nil, err
}
