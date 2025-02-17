package manager

import (
	"fmt"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2025-02-17 21:24
 */

func GetUserByID(id uint) (*UserManager, error) {
	user := &UserManager{}
	has, err := engine.ID(id).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("用户不存在")
	}
	return user, nil
}
