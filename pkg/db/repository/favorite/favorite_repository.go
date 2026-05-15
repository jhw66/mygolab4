package favorite

import (
	"context"

	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"gorm.io/gorm"
)

type FavoriteRepository interface {
	FindByUserAndVideo(ctx context.Context, userID, videoID string) (*model.Favorite, error)
	FindByUserAndVideoWithTx(ctx context.Context, tx *gorm.DB, userID, videoID string) (*model.Favorite, error)
	Create(ctx context.Context, favorite *model.Favorite) error
	CreateWithTx(ctx context.Context, tx *gorm.DB, favorite *model.Favorite) error
	DeleteByUserAndVideo(ctx context.Context, userID, videoID string) error
	DeleteByUserAndVideoWithTx(ctx context.Context, tx *gorm.DB, userID, videoID string) error
	DeleteByVideoIDWithTx(ctx context.Context, tx *gorm.DB, videoID string) error
	ListVideosByUserID(ctx context.Context, userID string) ([]model.Video, error)
	CountByVideoID(ctx context.Context, videoID string) (int64, error)
}

type favoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{
		db: db,
	}
}

func (r *favoriteRepository) FindByUserAndVideo(ctx context.Context, userID, videoID string) (*model.Favorite, error) {
	var favorite model.Favorite
	if err := r.db.WithContext(ctx).Where("user_id = ? AND video_id = ?", userID, videoID).Take(&favorite).Error; err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *favoriteRepository) FindByUserAndVideoWithTx(ctx context.Context, tx *gorm.DB, userID, videoID string) (*model.Favorite, error) {
	var favorite model.Favorite
	if err := tx.WithContext(ctx).Where("user_id = ? AND video_id = ?", userID, videoID).Take(&favorite).Error; err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *favoriteRepository) Create(ctx context.Context, favorite *model.Favorite) error {
	return r.db.WithContext(ctx).Create(favorite).Error
}

func (r *favoriteRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, favorite *model.Favorite) error {
	return tx.WithContext(ctx).Create(favorite).Error
}

func (r *favoriteRepository) DeleteByUserAndVideo(ctx context.Context, userID, videoID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Delete(&model.Favorite{}).Error
}

func (r *favoriteRepository) DeleteByUserAndVideoWithTx(ctx context.Context, tx *gorm.DB, userID, videoID string) error {
	return tx.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Delete(&model.Favorite{}).Error
}

func (r *favoriteRepository) DeleteByVideoIDWithTx(ctx context.Context, tx *gorm.DB, videoID string) error {
	return tx.WithContext(ctx).
		Where("video_id = ?", videoID).
		Delete(&model.Favorite{}).Error
}

func (r *favoriteRepository) ListVideosByUserID(ctx context.Context, userID string) ([]model.Video, error) {
	var videos []model.Video
	err := r.db.WithContext(ctx).Joins("JOIN favorite ON favorite.video_id = video.id").
		Where("favorite.user_id = ?", userID).Find(&videos).Error
	return videos, err
}

func (r *favoriteRepository) CountByVideoID(ctx context.Context, videoID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Favorite{}).Where("video_id = ?", videoID).Count(&count).Error
	return count, err
}
