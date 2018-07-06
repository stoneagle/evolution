package php

type Profile struct {
	UserId        int    `xorm:"not null pk INT(11)"`
	Name          string `xorm:"VARCHAR(255)"`
	PublicEmail   string `xorm:"VARCHAR(255)"`
	GravatarEmail string `xorm:"VARCHAR(255)"`
	GravatarId    string `xorm:"VARCHAR(32)"`
	Location      string `xorm:"VARCHAR(255)"`
	Website       string `xorm:"VARCHAR(255)"`
	Bio           string `xorm:"TEXT"`
	Timezone      string `xorm:"VARCHAR(40)"`
}
