package comment

import (
	"context"

	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	CreateWithTx(ctx context.Context, tx *gorm.DB, comment *model.Comment) error
	FindByIDAndVideo(ctx context.Context, commentID, videoID string) (*model.Comment, error)
	FindByIDAndVideoWithTx(ctx context.Context, tx *gorm.DB, commentID, videoID string) (*model.Comment, error)
	Delete(ctx context.Context, comment *model.Comment) error
	DeleteWithTx(ctx context.Context, tx *gorm.DB, comment *model.Comment) error
	DeleteByVideoIDWithTx(ctx context.Context, tx *gorm.DB, videoID string) error
	ListByVideoID(ctx context.Context, videoID string, page int, pageSize int) ([]model.Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (r *commentRepository) Create(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, comment *model.Comment) error {
	return tx.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) FindByIDAndVideo(ctx context.Context, commentID, videoID string) (*model.Comment, error) {
	var comment model.Comment
	if err := r.db.WithContext(ctx).Where("id = ? AND video_id = ?", commentID, videoID).Take(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) FindByIDAndVideoWithTx(ctx context.Context, tx *gorm.DB, commentID, videoID string) (*model.Comment, error) {
	var comment model.Comment
	if err := tx.WithContext(ctx).Where("id = ? AND video_id = ?", commentID, videoID).Take(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) Delete(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Delete(comment).Error
}

func (r *commentRepository) DeleteWithTx(ctx context.Context, tx *gorm.DB, comment *model.Comment) error {
	return tx.WithContext(ctx).Delete(comment).Error
}

func (r *commentRepository) DeleteByVideoIDWithTx(ctx context.Context, tx *gorm.DB, videoID string) error {
	return tx.WithContext(ctx).
		Where("video_id = ?", videoID).
		Delete(&model.Comment{}).Error
}

func (r *commentRepository) ListByVideoID(ctx context.Context, videoID string, page int, pageSize int) ([]model.Comment, error) {
	offset := (page - 1) * pageSize
	var comments []model.Comment
	err := r.db.WithContext(ctx).Preload("User").Where("video_id = ?", videoID).Order("created_at desc").
		Limit(pageSize).Offset(offset).Find(&comments).Error
	return comments, err
}
