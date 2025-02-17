package manager

import "time"

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 数据表单模型
 * @Date: 2025-02-17 07:26
 */

type UserManager struct {
	CsdID        uint      `xorm:"'csd_id' pk autoincr"`        // 用户账号自增ID
	Username     string    `xorm:"'username' unique notnull"`   // 用户名
	StudentID    string    `xorm:"'student_id' unique notnull"` // 学号
	RealName     string    `xorm:"'real_name' notnull"`         // 真实姓名
	PhoneNumber  string    `xorm:"'phone_number'"`              // 手机号
	Email        string    `xorm:"'email' unique"`              // 邮箱
	Salt         string    `xorm:"'salt' notnull"`              // 密码盐
	HashPassword string    `xorm:"'hash_password' notnull"`     // 加密后的密码
	CreatedAt    time.Time `xorm:"'created_at' created"`        // 账号创建时间
	UpdatedAt    time.Time `xorm:"'updated_at' updated"`        // 账号更新时间
	InviteBy     uint      `xorm:"'invite_by'"`                 // 邀请人ID
	Status       string    `xorm:"'status' default('active')"`  // 账号状态
}

func (UserManager) TableName() string {
	return "user_manager"
}

type AdminManager struct {
	CsdID    uint `xorm:"'csd_id' notnull"`    // 用户账号自增ID
	RoleID   uint `xorm:"'role_id' notnull"`   // 身份的自增ID
	InviteBy uint `xorm:"'invite_by' notnull"` // 邀请人：邀请码
}

func (AdminManager) TableName() string {
	return "admin_manager"
}

type RoleManager struct {
	RoleID      uint   `xorm:"'role_id' pk autoincr"`      // 身份的自增ID
	RoleName    string `xorm:"'role_name' unique notnull"` // 身份名
	Authorities string `xorm:"'authorities' notnull"`      // 权限的集合
}

func (RoleManager) TableName() string {
	return "role_manager"
}

type AuthoritiesManager struct {
	AuthoritiesID   uint   `xorm:"'authorities_id' pk autoincr"`      // 权限的ID
	AuthoritiesName string `xorm:"'authorities_name' unique notnull"` // 权限名
	AuthoritiesDesc string `xorm:"'authorities_desc' unique notnull"` // 权限解释
}

func (AuthoritiesManager) TableName() string {
	return "authorities_manager"
}
