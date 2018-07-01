package models

type SocialAccount struct {
	Id        int    `xorm:"not null pk autoincr INT(11)"`
	UserId    int    `xorm:"index INT(11)"`
	Provider  string `xorm:"not null unique(account_unique) VARCHAR(255)"`
	ClientId  string `xorm:"not null unique(account_unique) VARCHAR(255)"`
	Data      string `xorm:"TEXT"`
	Code      string `xorm:"unique VARCHAR(32)"`
	CreatedAt int    `xorm:"INT(11)"`
	Email     string `xorm:"VARCHAR(255)"`
	Username  string `xorm:"VARCHAR(255)"`
}
