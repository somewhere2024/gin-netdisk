package services

import (
	"errors"
	"gin-netdisk/internal/dao/mysql"
	"gin-netdisk/internal/models"
	"gin-netdisk/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtkey = []byte("thisisasecret")

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	}
}

func CreateToken(payload jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString(jwtkey)
}

func GetCurrentUser(token string) (jwt.MapClaims, error) {
	data := jwt.MapClaims{}
	jwtstr, err := jwt.ParseWithClaims(token, &data, Secret()) //注意指针
	if err != nil {
		utils.Logger.Error("token解析失败")
		return nil, err
	}
	if !jwtstr.Valid {
		utils.Logger.Error("token验证失败")
		return nil, err
	}

	return data, nil
}

func PasswordHash(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

func CheckPassword(password, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	// return err == nil
	if err != nil {
		return false
	}
	return true
}

func CreateUser(user *models.User) error {
	if user.Username == "" || user.Password_hash == "" {
		utils.Logger.Info("创建用户失败")
		return errors.New("创建用户失败")
	}
	rel := mysql.DB.Create(user)
	if rel.Error != nil {
		utils.Logger.Info("创建用户失败")
		return rel.Error
	}
	return nil
}

func Authenticate(username, password string) (*models.User, error) {
	user := &models.User{}
	rel := mysql.DB.Where("username = ?", username).First(user)
	if rel.Error != nil {
		utils.Logger.Info("用户名不存在")
		return nil, rel.Error
	}
	if !CheckPassword(password, user.Password_hash) {
		utils.Logger.Info("密码错误")
		return nil, errors.New("密码错误")
	}
	return user, nil
}
