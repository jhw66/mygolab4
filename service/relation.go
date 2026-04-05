package service

import (
	"github.com/jhw66/myvideo_lab4/model"
	"github.com/jhw66/myvideo_lab4/serializer"
)

type Relation struct{}

func (Relation) RelationAction(tid string, uid string) *serializer.Response {
	relation := model.Relation{
		UserID:       uid,
		TargetUserID: tid,
	}

	var count int64
	model.Db.Model(&model.Relation{}).Where("user_id = ? and target_user_id = ?", uid, tid).Count(&count)
	if count != 0 {
		if err := model.Db.Where("user_id = ? and target_user_id = ?", uid, tid).Delete(&relation).Error; err != nil {
			return &serializer.Response{
				Status: 500,
				Msg:    "取消关注失败",
			}
		}
		return &serializer.Response{
			Status: 200,
			Msg:    "取消关注成功",
		}
	}
	if err := model.Db.Create(&relation).Error; err != nil {
		return &serializer.Response{
			Status: 500,
			Msg:    "关注失败",
		}
	}
	return &serializer.Response{
		Status: 200,
		Msg:    "关注成功",
	}
}

func (Relation) FollowingList(uid string) *serializer.Response {
	var users []model.User
	if err := model.Db.Table("user").Joins("join relation on relation.target_user_id = user.id").Where("relation.user_id = ?", uid).
		Find(&users).Error; err != nil {
		return &serializer.Response{
			Status: 500,
			Msg:    "查询关注列表失败",
		}
	}

	return &serializer.Response{
		Status: 200,
		Data:   serializer.BuildUserList(&users),
		Msg:    "查询关注列表成功",
	}
}

func (Relation) FollowerList(uid string) *serializer.Response {
	var users []model.User
	if err := model.Db.Table("user").Joins("join relation on relation.user_id = user.id").Where("relation.target_user_id = ?", uid).
		Find(&users).Error; err != nil {
		return &serializer.Response{
			Status: 500,
			Msg:    "查询粉丝列表失败",
		}
	}

	return &serializer.Response{
		Status: 200,
		Data:   serializer.BuildUserList(&users),
		Msg:    "查询粉丝列表成功",
	}
}

func (Relation) FriendList(uid string) *serializer.Response {
	var users []model.User

	if err := model.Db.Table("user").
		Where("id in (?)", model.Db.Table("relation").Select("target_user_id").Where("user_id = ?", uid)).
		Where("id in (?)", model.Db.Table("relation").Select("user_id").Where("target_user_id = ?", uid)).
		Find(&users).Error; err != nil {
		return &serializer.Response{
			Status: 500,
			Msg:    "查询好友列表失败",
		}
	}

	return &serializer.Response{
		Status: 200,
		Data:   serializer.BuildUserList(&users),
		Msg:    "查询好友列表成功",
	}
}
