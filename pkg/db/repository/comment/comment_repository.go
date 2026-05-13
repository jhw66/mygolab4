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
	FindByIDAndVideoUnscoped(ctx context.Context, commentID, videoID string) (*model.Comment, error)
	FindByIDAndVideoWithTx(ctx context.Context, tx *gorm.DB, commentID, videoID string) (*model.Comment, error)
	Delete(ctx context.Context, comment *model.Comment) error
	DeleteWithTx(ctx context.Context, tx *gorm.DB, comment *model.Comment) error
	DeleteByRootIDWithTx(ctx context.Context, tx *gorm.DB, rootID string) error
	DeleteByVideoIDWithTx(ctx context.Context, tx *gorm.DB, videoID string) error
	ListByVideoID(ctx context.Context, videoID string, page int, pageSize int) ([]model.Comment, error)
	ListRootByVideoID(ctx context.Context, videoID string, page int, pageSize int) ([]model.Comment, error)
	ListRepliesByRootID(ctx context.Context, videoID, rootID string, page int, pageSize int) ([]model.Comment, error)
	CountRootByVideoID(ctx context.Context, videoID string) (int64, error)
	CountRepliesByRootID(ctx context.Context, videoID, rootID string) (int64, error)
	CountRepliesByRootIDWithTx(ctx context.Context, tx *gorm.DB, videoID, rootID string) (int64, error)
	IncrementFavoriteCountWithTx(ctx context.Context, tx *gorm.DB, commentID string) error
	DecrementFavoriteCountWithTx(ctx context.Context, tx *gorm.DB, commentID string) error
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

func (r *commentRepository) FindByIDAndVideoUnscoped(ctx context.Context, commentID, videoID string) (*model.Comment, error) {
	var comment model.Comment
	if err := r.db.WithContext(ctx).Unscoped().Where("id = ? AND video_id = ?", commentID, videoID).Take(&comment).Error; err != nil {
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

func (r *commentRepository) DeleteByRootIDWithTx(ctx context.Context, tx *gorm.DB, rootID string) error {
	return tx.WithContext(ctx).
		Where("root_id = ?", rootID).
		Delete(&model.Comment{}).Error
}

func (r *commentRepository) DeleteByVideoIDWithTx(ctx context.Context, tx *gorm.DB, videoID string) error {
	return tx.WithContext(ctx).
		Where("video_id = ?", videoID).
		Delete(&model.Comment{}).Error
}

func (r *commentRepository) ListByVideoID(ctx context.Context, videoID string, page int, pageSize int) ([]model.Comment, error) {
	var comments []model.Comment
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Unscoped().Preload("User").
		Where("video_id = ?", videoID).
		Order("created_at desc").Limit(pageSize).Offset(offset).Find(&comments).Error
	return comments, err
}

func (r *commentRepository) ListRootByVideoID(ctx context.Context, videoID string, page int, pageSize int) ([]model.Comment, error) {
	var comments []model.Comment
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Unscoped().Preload("User").
		Where("video_id = ? AND comment_id IS NULL", videoID).
		Order("created_at desc").Limit(pageSize).Offset(offset).Find(&comments).Error
	return comments, err
}

func (r *commentRepository) ListRepliesByRootID(ctx context.Context, videoID, rootID string, page int, pageSize int) ([]model.Comment, error) {
	var comments []model.Comment
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Unscoped().Preload("User").
		Where("video_id = ? AND root_id = ?", videoID, rootID).
		Order("created_at asc").Limit(pageSize).Offset(offset).Find(&comments).Error
	return comments, err
}

func (r *commentRepository) CountRootByVideoID(ctx context.Context, videoID string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Unscoped().Model(&model.Comment{}).
		Where("video_id = ? AND comment_id IS NULL", videoID).
		Count(&total).Error
	return total, err
}

func (r *commentRepository) CountRepliesByRootID(ctx context.Context, videoID, rootID string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Unscoped().Model(&model.Comment{}).
		Where("video_id = ? AND root_id = ?", videoID, rootID).
		Count(&total).Error
	return total, err
}

func (r *commentRepository) CountRepliesByRootIDWithTx(ctx context.Context, tx *gorm.DB, videoID, rootID string) (int64, error) {
	var total int64
	err := tx.WithContext(ctx).Unscoped().Model(&model.Comment{}).
		Where("video_id = ? AND root_id = ?", videoID, rootID).
		Count(&total).Error
	return total, err
}

func (r *commentRepository) IncrementFavoriteCountWithTx(ctx context.Context, tx *gorm.DB, commentID string) error {
	return tx.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ?", commentID).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
}

func (r *commentRepository) DecrementFavoriteCountWithTx(ctx context.Context, tx *gorm.DB, commentID string) error {
	return tx.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ? AND favorite_count > 0", commentID).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
}
