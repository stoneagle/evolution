package models

type Token struct {
	UserId    int    `xorm:"not null pk unique(token_unique) INT(11)"`
	Code      string `xorm:"not null pk unique(token_unique) VARCHAR(32)"`
	CreatedAt int    `xorm:"not null INT(11)"`
	Type      int    `xorm:"not null pk unique(token_unique) SMALLINT(6)"`
}
