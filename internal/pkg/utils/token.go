package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string)(string, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}


func CheckPassword(password, hash string)bool{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}


func GenerateJWT(userID uint)(string, error){
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == ""{
		secretKey = "my_super_secret_key"
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}


func ValidateJWT(tokenString string) (uint, error){
	secretKey := os.Getenv("JWT_SECRET")

	if secretKey == ""{
		secretKey = "my_super_secret_key"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, errors.New("not'g'ri imzolash usuli")
		}
		return  []byte(secretKey), nil
	})

	if err != nil || !token.Valid{
		return 0, errors.New("muddati o'tgan token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("token ma'lumotlarini o'qib bolmadi")
	}

	userID := uint(claims["user_id"].(float64))
	return userID, nil
}