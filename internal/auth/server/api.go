package server

import (
	"context"
	"fmt"

	"github.com/Vysogota99/HousingSearch/internal/auth/jwt"
	"github.com/Vysogota99/HousingSearch/internal/auth/store/postgresstore"
	"github.com/Vysogota99/HousingSearch/pkg/authService"
)

// SignupUser - создает пользователя в базе данных, генерирует для него jwt токен,
//				после чего возвращает его в ответ
func (s *GRPCServer) SignupUser(ctx context.Context, req *authService.SignUPUserRequest) (*authService.SignUPUserResponse, error) {
	if err := s.CreateUser(req.User); err != nil {
		return nil, err
	}

	tokens, err := jwt.CreateToken(req.User.TelephoneNumber, req.User.ID)
	if err != nil {
		return nil, err
	}

	if saveErr := s.CreateAuth(req.User.TelephoneNumber, tokens); saveErr != nil {
		return nil, err
	}
	req.User.Password = ""

	return &authService.SignUPUserResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         req.User,
	}, nil
}

// LogOutUser - выход пользователя из системы путем удаления jwt токена
func (s *GRPCServer) LogOutUser(ctx context.Context, req *authService.LogOutRequest) (*authService.LogOutResponse, error) {
	au, err := jwt.ExtractTokenMetadata(req.Jwt)
	if err != nil {
		return nil, fmt.Errorf("Невалидный jwt. %w", err)
	}

	deleted, delErr := s.DeleteAuth(au.AccessUUID)
	if delErr != nil || deleted == 0 {
		return nil, fmt.Errorf("Пользователь еще не авторизован")
	}

	return &authService.LogOutResponse{
		Message: "Deleted",
	}, nil
}

// CheckAuthUser - проверяет авторизован ли пользователь
func (s *GRPCServer) CheckAuthUser(ctx context.Context, req *authService.CheckAuthRequest) (*authService.CheckAuthResponse, error) {
	if err := jwt.TokenValid(req.Jwt); err != nil {
		return nil, err
	}

	jwtData, err := jwt.ExtractTokenMetadata(req.Jwt)
	if err != nil {
		return nil, err
	}

	if err := s.FetchAuth(jwtData.AccessUUID); err != nil {
		return nil, err
	}

	return &authService.CheckAuthResponse{
		TelephoneNumber: jwtData.TelephoneNumber,
		UserID:          jwtData.UserID,
	}, nil
}

// LoginUser - авторизация пользователя
func (s *GRPCServer) LoginUser(ctx context.Context, req *authService.LoginUserRequest) (*authService.SignUPUserResponse, error) {
	user, err := s.Store.User().GetUser(req.User.TelephoneNumber)
	if err != nil {
		return nil, fmt.Errorf("Пользователь не найден")
	}

	if areEqual := postgresstore.ComparePasswords(user.Password, req.User.Password); !areEqual {
		return nil, fmt.Errorf("Неправильное имя пользователя или пароль")
	}

	user.Password = ""
	tokens, err := jwt.CreateToken(user.TelephoneNumber, user.ID)
	if err != nil {
		return nil, err
	}

	if saveErr := s.CreateAuth(user.TelephoneNumber, tokens); saveErr != nil {
		return nil, err
	}
	return &authService.SignUPUserResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user,
	}, nil
}
