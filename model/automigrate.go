package model

func Migrate() {
	Db.AutoMigrate(&User{}, &Video{})
}
