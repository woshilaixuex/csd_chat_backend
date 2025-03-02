package chat

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2025-02-22 21:15
 */

func InsertIMMsg(content *IMMsgContent) (err error) {
	session := engine.NewSession()
	defer session.Close()
	_, err = engine.Insert(content)
	if err != nil {
		return
	}
	err = session.Commit()
	return
}
