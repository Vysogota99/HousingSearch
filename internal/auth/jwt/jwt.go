package jwt

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

// CreateToken ...
func CreateToken(telephoneNumer string, id int64, expAccessTime, expRefreshTime int64) (*TokenDetails, error) {
	tokenDet := &TokenDetails{}
	tokenDet.AccTExpires = time.Now().Add(time.Minute * time.Duration(expAccessTime)).Unix()
	tokenDet.AccessUUID = uuid.NewV4().String()

	tokenDet.RefTExpires = time.Now().Add(time.Hour * 24 * time.Duration(expRefreshTime)).Unix()
	tokenDet.RefreshUUID = uuid.NewV4().String()

	var err error
	// Creating access token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = tokenDet.AccessUUID
	atClaims["telephone_number"] = telephoneNumer
	atClaims["user_id"] = id
	atClaims["exp"] = tokenDet.AccTExpires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	tokenDet.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// Creating refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = tokenDet.RefreshUUID
	rtClaims["telephone_number"] = telephoneNumer
	rtClaims["user_id"] = id
	rtClaims["exp"] = tokenDet.RefTExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	tokenDet.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return tokenDet, nil
}

// VerifyToken - verify the token
func VerifyToken(headerString string) (*jwt.Token, error) {
	tokenString := ExtractToken(headerString)
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
func ExtractToken(bearerToken string) string {
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// TokenValid - check the validity of this token
func TokenValid(tokenString string) error {
	token, err := VerifyToken(tokenString)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

// ExtractTokenMetadata - extract the token metadata that will lookup in our redis store
func ExtractTokenMetadata(tokenString string) (*AccessDetails, error) {
	token, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, fmt.Errorf("No access_uuid in jwt token")
		}
		telephoneNumber, ok := claims["telephone_number"].(string)
		if !ok {
			return nil, fmt.Errorf("No telephone_number in jwt token")
		}
		userID, ok := claims["user_id"]
		if !ok {
			return nil, fmt.Errorf("No user_id in jwt token")
		}
		userIDInt64 := int64(userID.(float64))
		return &AccessDetails{
			AccessUUID:      accessUUID,
			TelephoneNumber: telephoneNumber,
			UserID:          userIDInt64,
		}, nil
	}
	return nil, err
}
