package encrypt

import (
	"github.com/xh-polaris/psych-user/biz/infrastructure/consts"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

var defaultPwd string
var once sync.Once

func GetDefaultPwd() string {
	once.Do(func() {
		var pwd string
		var err error
		if pwd, err = BcryptEncrypt(consts.DefaultPassword); err != nil {
			panic(err)
		}
		defaultPwd = pwd
	})
	return defaultPwd
}

// BcryptEncrypt Bcrypt加密函数
func BcryptEncrypt(password string) (string, error) {
	// 使用 bcrypt 生成哈希值
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// BcryptCheck Bcrypt校验函数
func BcryptCheck(password string, hash string) bool {
	// 检查密码是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
