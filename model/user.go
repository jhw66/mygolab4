package model

type User struct {
	BaseModel
	UserName       string `gorm:"unique;not null"`
	NickName       string `gorm:"not null"`
	PasswordDigest string
	Avatar         string
	Videos         []Video `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (User) TableName() string {
	return "user"
}

func GetUserByID(id string) (*User, error) {
	var user User
	err := Db.Where("id = ?", id).Take(&user).Error
	return &user, err
}
