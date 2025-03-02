package chat

import "time"

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 数据表单模型
 * @Date: 2025-02-22 15:57
 */

type IMMsgContent struct { // 消息
	Mid         uint64    `xorm:"pk autoincr 'mid'"`
	Content     string    `xorm:"'content'"`
	SenderID    uint64    `xorm:"'sender_id'"`
	RecipientID uint64    `xorm:"'recipient_id'"`
	MsgType     int       `xorm:"'msg_type'"`
	CreateTime  time.Time `xorm:"'create_time' created"`
}

type IMMsgRelation struct { // 消息索引表
	OwnerUID   uint64    `xorm:"'owner_uid'"`
	OtherUID   uint64    `xorm:"'other_uid'"`
	Mid        uint64    `xorm:"'mid'"`
	Type       int       `xorm:"'type'"`
	CreateTime time.Time `xorm:"'create_time' created"`
}

type IMMsgContact struct { // 联系表
	OwnerUID   uint64    `xorm:"'owner_uid'"`
	OtherUID   uint64    `xorm:"'other_uid'"`
	Mid        uint64    `xorm:"'mid'"`
	Type       int       `xorm:"'type'"`
	CreateTime time.Time `xorm:"'create_time' created"`
}
