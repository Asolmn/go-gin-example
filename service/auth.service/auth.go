package auth_service

import "github.com/Asolmn/go-gin-example/models"

type Auth struct {
	Username string
	Password string
}

// 检查用户认证
func (a *Auth) Check() (bool, error) {
	return models.CheckAuth(a.Username, a.Password)
}
