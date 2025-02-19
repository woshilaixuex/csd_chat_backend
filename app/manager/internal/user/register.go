package user

import (
	"github.com/jinzhu/copier"
	"github.com/woshilaixuex/csd_chat_backend/app/util/model/manager"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xerr"
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
	Password    string
}

func (entity *RegisterEntity) HashPassword() string {
	return entity.Password
}

func Register(entity RegisterEntity) (token string, err error) {
	user, err := manager.GetUserEmail(entity.Email)
	if err != nil {
		return
	}
	if user != nil {
		return token, xerr.UserExists
	}
	user = new(manager.UserManager)
	err = copier.Copy(user, &entity)
	if err != nil {
		return
	}
	err = manager.InsertUser(user)
	if err != nil {
		return
	}

	return
}
func CheckInvitCode(code string) bool {
	return true
}
