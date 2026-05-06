package relation

import (
	"context"

	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"gorm.io/gorm"
)

type RelationRepository interface {
	ExistsByUserAndTarget(ctx context.Context, userID, targetUserID string) (bool, error)
	Create(ctx context.Context, userID, targetUserID string) error
	DeleteByUserAndTarget(ctx context.Context, userID, targetUserID string) error
	ListFollowingUsers(ctx context.Context, userID string) ([]model.User, error)
	ListFollowerUsers(ctx context.Context, userID string) ([]model.User, error)
	ListFriendUsers(ctx context.Context, userID string) ([]model.User, error)
}

type relationRepository struct {
	db *gorm.DB
}

func NewRelationRepository(db *gorm.DB) RelationRepository {
	return &relationRepository{
		db: db,
	}
}

func (r *relationRepository) ExistsByUserAndTarget(ctx context.Context, userID, targetUserID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Relation{}).
		Where("user_id = ? and target_user_id = ?", userID, targetUserID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *relationRepository) Create(ctx context.Context, userId, targetUserId string) error {
	return r.db.WithContext(ctx).Create(&model.Relation{
		UserID:       userId,
		TargetUserID: targetUserId,
	}).Error
}

func (r *relationRepository) DeleteByUserAndTarget(ctx context.Context, userID, targetUserID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? and target_user_id = ?", userID, targetUserID).
		Delete(&model.Relation{}).Error
}

func (r *relationRepository) ListFollowingUsers(ctx context.Context, userID string) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).
		Table("user").
		Joins("join relation on relation.target_user_id = user.id").
		Where("relation.user_id = ? AND relation.deleted_at is null", userID).
		Find(&users).Error
	return users, err
}

func (r *relationRepository) ListFollowerUsers(ctx context.Context, userID string) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).
		Table("user").
		Joins("join relation on relation.user_id = user.id").
		Where("relation.target_user_id = ?", userID).
		Find(&users).Error
	return users, err
}

func (r *relationRepository) ListFriendUsers(ctx context.Context, userID string) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).
		Table("user").
		Where("id in (?)", r.db.Table("relation").Select("target_user_id").Where("user_id = ?", userID)).
		Where("id in (?)", r.db.Table("relation").Select("user_id").Where("target_user_id = ?", userID)).
		Find(&users).Error
	return users, err
}
