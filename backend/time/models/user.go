package models

type User struct {
	Id               int    `xorm:"not null pk autoincr INT(11)"`
	Username         string `xorm:"not null unique VARCHAR(255)"`
	Email            string `xorm:"not null unique VARCHAR(255)"`
	PasswordHash     string `xorm:"not null VARCHAR(60)"`
	AuthKey          string `xorm:"not null VARCHAR(32)"`
	ConfirmedAt      int    `xorm:"INT(11)"`
	UnconfirmedEmail string `xorm:"VARCHAR(255)"`
	BlockedAt        int    `xorm:"INT(11)"`
	RegistrationIp   string `xorm:"VARCHAR(45)"`
	CreatedAt        int    `xorm:"not null INT(11)"`
	UpdatedAt        int    `xorm:"not null INT(11)"`
	Flags            int    `xorm:"not null default 0 INT(11)"`
	LastLoginAt      int    `xorm:"INT(11)"`
}
