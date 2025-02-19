package user

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2025-02-18 21:44
 */

type LoginEntity struct {
	Email    string
	Password string
}

func Login(entity LoginEntity) (err error) {
	// user, err := manager.GetUserEmail(entity.Email)
	return
}
