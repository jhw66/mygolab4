package commentfavorite

import (
	"context"

	"github.com/jhw66/myvideo_lab4/pkg/db/model"
	"gorm.io/gorm"
)

type CommentFavoriteRepository interface {
	FindByUserAndComment(ctx context.Context, userID, commentID string) (*model.CommentFavorite, error)
	FindByUserAndCommentWithTx(ctx context.Context, tx *gorm.DB, userID, commentID string) (*model.CommentFavorite, error)
	Create(ctx context.Context, favorite *model.CommentFavorite) error
	CreateWithTx(ctx context.Context, tx *gorm.DB, favorite *model.CommentFavorite) error
	DeleteByUserAndComment(ctx context.Context, userID, commentID string) error
	DeleteByUserAndCommentWithTx(ctx context.Context, tx *gorm.DB, userID, commentID string) error
	DeleteByCommentID(ctx context.Context, commentID string) error
	DeleteByCommentIDWithTx(ctx context.Context, tx *gorm.DB, commentID string) error
	DeleteByRootIDWithTx(ctx context.Context, tx *gorm.DB, rootID string) error
	DeleteByVideoIDWithTx(ctx context.Context, tx *gorm.DB, videoID string) error
}

type commentFavoriteRepository struct {
	db *gorm.DB
}

func NewCommentFavoriteRepository(db *gorm.DB) CommentFavoriteRepository {
	return &commentFavoriteRepository{db: db}
}

func (r *commentFavoriteRepository) FindByUserAndComment(ctx context.Context, userID, commentID string) (*model.CommentFavorite, error) {
	var favorite model.CommentFavorite
	if err := r.db.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).Take(&favorite).Error; err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *commentFavoriteRepository) FindByUserAndCommentWithTx(ctx context.Context, tx *gorm.DB, userID, commentID string) (*model.CommentFavorite, error) {
	var favorite model.CommentFavorite
	if err := tx.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).Take(&favorite).Error; err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *commentFavoriteRepository) Create(ctx context.Context, favorite *model.CommentFavorite) error {
	return r.db.WithContext(ctx).Create(favorite).Error
}

func (r *commentFavoriteRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, favorite *model.CommentFavorite) error {
	return tx.WithContext(ctx).Create(favorite).Error
}

func (r *commentFavoriteRepository) DeleteByUserAndComment(ctx context.Context, userID, commentID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		Delete(&model.CommentFavorite{}).Error
}

func (r *commentFavoriteRepository) DeleteByUserAndCommentWithTx(ctx context.Context, tx *gorm.DB, userID, commentID string) error {
	return tx.WithContext(ctx).
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		Delete(&model.CommentFavorite{}).Error
}

func (r *commentFavoriteRepository) DeleteByCommentID(ctx context.Context, commentID string) error {
	return r.db.WithContext(ctx).
		Where("comment_id = ?", commentID).
		Delete(&model.CommentFavorite{}).Error
}

func (r *commentFavoriteRepository) DeleteByCommentIDWithTx(ctx context.Context, tx *gorm.DB, commentID string) error {
	return tx.WithContext(ctx).
		Where("comment_id = ?", commentID).
		Delete(&model.CommentFavorite{}).Error
}

func (r *commentFavoriteRepository) DeleteByRootIDWithTx(ctx context.Context, tx *gorm.DB, rootID string) error {
	return tx.WithContext(ctx).
		Where("comment_id IN (?)", tx.Model(&model.Comment{}).Select("id").Where("root_id = ?", rootID)).
		Delete(&model.CommentFavorite{}).Error
}

func (r *commentFavoriteRepository) DeleteByVideoIDWithTx(ctx context.Context, tx *gorm.DB, videoID string) error {
	return tx.WithContext(ctx).
		Where("comment_id IN (?)", tx.Model(&model.Comment{}).Select("id").Where("video_id = ?", videoID)).
		Delete(&model.CommentFavorite{}).Error
}
