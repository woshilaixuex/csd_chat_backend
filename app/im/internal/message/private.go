package message

import (
	"github.com/jinzhu/copier"
	"github.com/woshilaixuex/csd_chat_backend/app/util/model/chat"
	"github.com/woshilaixuex/csd_chat_backend/app/util/security/xtoken"
	"github.com/woshilaixuex/csd_chat_backend/app/util/xredis"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 私聊业务
 * @Date: 2025-02-22 17:14
 */

type ChatEntity struct {
	Token    string
	MID      uint64
	MContext string
	OtherUId uint64
}

func MessageSave(entity *ChatEntity) (succeed bool, err error) {
	id, err := xtoken.ParseJwtToken(entity.Token)
	if err != nil {
		return
	}
	time, err := xredis.GetTimeFromRedis()
	if err != nil {
		return
	}
	content := new(chat.IMMsgContent)
	err = copier.Copy(content, entity)
	if err != nil {
		return
	}
	content.CreateTime = time
	content.SenderID = id
	err = chat.InsertIMMsg(content)
	if err != nil {
		return
	}
	return !succeed, nil
}
