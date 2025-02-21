package manager

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2025-02-17 21:24
 */

func GetUserByID(id uint64) (found bool, user *UserManager, err error) {
	user = new(UserManager)
	found, err = engine.Where("csd_id = ?", id).Get(user)
	if err != nil {
		return
	}
	return
}
func GetUserEmail(email string) (found bool, user *UserManager, err error) {
	user = new(UserManager)
	found, err = engine.Where("email = ?", email).Get(user)
	if err != nil {
		return
	}
	return
}

func InsertUser(user *UserManager) (err error) {
	_, err = engine.Insert(user)
	return
}
