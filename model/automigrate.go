package model

func Migrate() {
	Db.SetupJoinTable(&User{}, "video", &Favorite{})
	Db.SetupJoinTable(&Video{}, "user", &Favorite{})
	Db.SetupJoinTable(&User{}, "video", &Comment{})
	Db.SetupJoinTable(&Video{}, "user", &Comment{})
	Db.AutoMigrate(&User{}, &Video{}, &Favorite{}, &Comment{})
}
