package model

type User struct {
	BaseModel
	UserName          string `gorm:"unique;not null"`
	NickName          string `gorm:"not null"`
	PasswordDigest    string
	Avatar            string
	TotpSecret        string
	TotpPendingSecret string
	TotpEnabled       bool
	Videos            []Video `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (User) TableName() string {
	return "user"
}
