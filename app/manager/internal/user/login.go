package user

import (
	"github.com/woshilaixuex/csd_chat_backend/app/util/model/manager"
	"github.com/woshilaixuex/csd_chat_backend/app/util/security/encryption"
	"github.com/woshilaixuex/csd_chat_backend/app/util/security/xtoken"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xerr"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 登录业务
 * @Date: 2025-02-18 21:44
 */

type LoginEntity struct {
	Email    string
	Password string
}

func Login(entity LoginEntity) (token string, err error) {
	found, user, err := manager.GetUserEmail(entity.Email)
	if err != nil {
		return
	}
	if !found {
		return token, xerr.UserNotExist
	}
	if verified := encryption.VerifyPassword(user.HashPassword, entity.Password, user.Salt); verified {
		return token, xerr.UnDefinedError
	}
	token, err = xtoken.GetJwtToken(user.CsdID)
	return
}
