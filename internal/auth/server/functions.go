package server

import (
	"fmt"
	"time"

	"github.com/Vysogota99/school/internal/auth/jwt"
	"github.com/Vysogota99/school/internal/auth/redis"
	"github.com/Vysogota99/school/pkg/authService"
)

// CreateUser - создает пользователя в бд или возвращает ошибку, если такой пользователь
//				уже существует
func (s *GRPCServer) CreateUser(user *authService.User) error {
	_, err := s.Store.User().GetUser(user.TelephoneNumber)
	if err == nil {
		return fmt.Errorf("Пользователь с таким номером телефона уже существует")
	}

	if err := s.Store.User().CreateUser(user); err != nil {
		return fmt.Errorf("Ошибка при создании пользователя. %w", err)
	}

	return nil
}

// CreateAuth - function that will be used to save the JWTs metadata in redis
func (s *GRPCServer) CreateAuth(telephoneNumber string, td *jwt.TokenDetails) error {
	redis, err := redis.NewClient(s.Conf.RedisDSN)
	redisClient := redis.Client
	defer redisClient.Close()

	if err != nil {
		return err
	}

	accessTExp := time.Unix(td.AccTExpires, 0)
	refresTExp := time.Unix(td.RefTExpires, 0)
	now := time.Now()

	if errAccess := redisClient.Set(td.AccessUUID, telephoneNumber, accessTExp.Sub(now)).Err(); errAccess != nil {
		return errAccess
	}

	if errRefresh := redisClient.Set(td.RefreshUUID, telephoneNumber, refresTExp.Sub(now)).Err(); errRefresh != nil {
		return errRefresh
	}

	return nil
}

// DeleteAuth - delete a JWT metadata from redis
func (s *GRPCServer) DeleteAuth(givenUUID string) (int64, error) {
	redis, err := redis.NewClient(s.Conf.RedisDSN)
	redisClient := redis.Client
	defer redisClient.Close()

	deleted, err := redisClient.Del(givenUUID).Result()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}

// FetchAuth -  accepts the AccessDetails from the ExtractTokenMetadata function,
//				then looks it up in redis. If the record is not found, it may mean
//				the token has expired, hence an error is thrown.
func (s *GRPCServer) FetchAuth(accessUUID string) error {
	redis, err := redis.NewClient(s.Conf.RedisDSN)
	redisClient := redis.Client
	defer redisClient.Close()

	value, err := redisClient.Get(accessUUID).Result()
	if err != nil || value == "" {
		return fmt.Errorf("Пользователь еще не авторизован")
	}

	return nil
}
