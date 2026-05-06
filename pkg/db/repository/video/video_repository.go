package video

import (
	"context"

	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"gorm.io/gorm"
)

type VideoRepository interface {
	SearchByTitle(ctx context.Context, keyword string) ([]model.Video, error)
	FindByID(ctx context.Context, id string) (*model.Video, error)
	FindByIDWithTx(ctx context.Context, tx *gorm.DB, id string) (*model.Video, error)
	FindByUserID(ctx context.Context, userID string) ([]model.Video, error)
	FindByIDs(ctx context.Context, ids []string) ([]model.Video, error)
	ListRankSeed(ctx context.Context) ([]model.Video, error)
	Create(ctx context.Context, video *model.Video) error
	CreateWithTx(ctx context.Context, tx *gorm.DB, video *model.Video) error
	Update(ctx context.Context, video *model.Video) error
	UpdateWithTx(ctx context.Context, tx *gorm.DB, video *model.Video) error
	DeleteByIDWithTx(ctx context.Context, tx *gorm.DB, id string) (*model.Video, error)
	UpdateByID(ctx context.Context, id string, favoriteCount uint64, commentCount uint64, hotScore uint) error
}

type videoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) VideoRepository {
	return &videoRepository{
		db: db,
	}
}

func (r *videoRepository) SearchByTitle(ctx context.Context, keyword string) ([]model.Video, error) {
	var videos []model.Video
	err := r.db.WithContext(ctx).Where("title like ?", "%"+keyword+"%").Find(&videos).Error
	return videos, err
}

func (r *videoRepository) FindByID(ctx context.Context, id string) (*model.Video, error) {
	var video model.Video
	if err := r.db.WithContext(ctx).Where("id = ?", id).Take(&video).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) FindByIDWithTx(ctx context.Context, tx *gorm.DB, id string) (*model.Video, error) {
	var video model.Video
	if err := tx.WithContext(ctx).Where("id = ?", id).Take(&video).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) FindByUserID(ctx context.Context, userID string) ([]model.Video, error) {
	var videos []model.Video
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&videos).Error
	return videos, err
}

func (r *videoRepository) FindByIDs(ctx context.Context, ids []string) ([]model.Video, error) {
	var videos []model.Video
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&videos).Error
	return videos, err
}

func (r *videoRepository) ListRankSeed(ctx context.Context) ([]model.Video, error) {
	var videos []model.Video
	err := r.db.WithContext(ctx).Select("id,favorite_count,comment_count").Find(&videos).Error
	return videos, err
}

func (r *videoRepository) Create(ctx context.Context, video *model.Video) error {
	return r.db.WithContext(ctx).Create(video).Error
}

func (r *videoRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, video *model.Video) error {
	return tx.WithContext(ctx).Create(video).Error
}

func (r *videoRepository) Update(ctx context.Context, video *model.Video) error {
	return r.db.WithContext(ctx).Save(video).Error
}

func (r *videoRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, video *model.Video) error {
	return tx.WithContext(ctx).Save(video).Error
}

func (r *videoRepository) DeleteByIDWithTx(ctx context.Context, tx *gorm.DB, id string) (*model.Video, error) {
	var video model.Video
	if err := tx.WithContext(ctx).Where("id = ?", id).Take(&video).Error; err != nil {
		return nil, err
	}
	if err := tx.WithContext(ctx).Delete(&video).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) UpdateByID(ctx context.Context, id string, favoriteCount uint64, commentCount uint64, hotScore uint) error {
	return r.db.WithContext(ctx).Model(&model.Video{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"favorite_count": favoriteCount,
			"comment_count":  commentCount,
			"hot_score":      hotScore,
		}).Error
}
