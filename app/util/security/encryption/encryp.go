package encryption

import (
	"golang.org/x/crypto/bcrypt"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 加密库
 * @Date: 2025-02-19 20:46
 */

func EncryptPassword(password string) (salt string, encryptedPassword string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	salt = string(hashedPassword[:29])              // 提取出盐
	encryptedPassword = string(hashedPassword[29:]) // 加密后的密码

	return salt, encryptedPassword, nil
}

func VerifyPassword(encryptedPassword, password, salt string) bool {
	fullEncryptedPassword := salt + encryptedPassword
	err := bcrypt.CompareHashAndPassword([]byte(fullEncryptedPassword), []byte(password))
	return err == nil
}
