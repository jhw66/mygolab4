package model

func Migrate() {
	Db.AutoMigrate(&User{}, &Vedio{})
}
