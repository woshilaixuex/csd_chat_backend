package user

import (
	"github.com/jinzhu/copier"
	"github.com/woshilaixuex/csd_chat_backend/app/util/model/manager"
	"github.com/woshilaixuex/csd_chat_backend/app/util/security/encryption"
	"github.com/woshilaixuex/csd_chat_backend/app/util/security/xtoken"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xerr"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xredis"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 用户注册业务
 * @Date: 2025-02-17 21:44
 */
type RegisterEntity struct {
	UserName    string
	StudentID   string
	RealName    string
	PhoneNumber string
	Email       string
	Salt        string
	Password    string
}

func (entity *RegisterEntity) HashPassword() string {
	return entity.Password
}

// 项目初始化时重置root
func RootInit() {

}
func Register(entity *RegisterEntity) (token string, err error) {
	found, user, err := manager.GetUserEmail(entity.Email)
	if err != nil {
		return
	}
	if found {
		return token, xerr.UserExists
	}
	if !checkInvitCode("") {
		return token, xerr.InviteError
	}
	salt, hash, err := encryption.EncryptPassword(entity.Password)
	if err != nil {
		return
	}
	entity.Salt = salt
	entity.Password = hash
	user = new(manager.UserManager)
	user.CsdID, err = xredis.GetNewGlobalCsdID()
	if err != nil {
		return
	}
	err = copier.Copy(user, &entity)
	if err != nil {
		return
	}
	err = manager.InsertUser(user)
	if err != nil {
		return
	}
	token, err = xtoken.GetJwtToken(user.CsdID)
	return
}
func checkInvitCode(code string) bool {
	return code == ""
}
