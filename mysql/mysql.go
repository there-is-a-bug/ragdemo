package mysql

import (
	"MyMemory/model"
	"context"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserDocumentRepo struct {
	db *gorm.DB
}

var (
	so   sync.Once
	Repo *UserDocumentRepo
)

func init() {
	so.Do(func() {
		dsn := model.SqlDsn
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			// 开发环境打印全部SQL；生产改成 logger.Warn / logger.Error
			Logger: logger.Default.LogMode(logger.Error),
		})
		if err != nil {
			panic(err)
		}
		Repo = &UserDocumentRepo{
			db: db,
		}
	})
}

func GetUserDocumentRepo() *UserDocumentRepo {
	return Repo
}

// Create 创建文档
func (r *UserDocumentRepo) Create(ctx context.Context, doc *model.UserDocument) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

// GetByID 根据数据库自增id查询（自动过滤软删除）
func (r *UserDocumentRepo) GetByID(ctx context.Context, id uint64) (*model.UserDocument, error) {
	var res model.UserDocument
	err := r.db.WithContext(ctx).First(&res, id).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetByUserDocID 根据 user_id + doc_id 查询业务文档
func (r *UserDocumentRepo) GetByUserDocID(ctx context.Context, userID int64, docID string) (*model.UserDocument, error) {
	var res model.UserDocument
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND doc_id = ?", userID, docID).
		First(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// Update 更新非零字段；doc传入要修改的字段
func (r *UserDocumentRepo) Update(ctx context.Context, doc *model.UserDocument) error {
	return r.db.WithContext(ctx).Save(doc).Error
}

// UpdateColumns 部分字段更新，推荐，只更新指定字段
func (r *UserDocumentRepo) UpdateColumns(ctx context.Context, id uint64, data map[string]any) error {
	return r.db.WithContext(ctx).
		Model(&model.UserDocument{ID: id}).
		Updates(data).Error
}

// Delete 软删除，gorm.DeletedAt，只打deleted_at标记
// ⚠️业务层调用完这里，**务必同步删除ES对应文档**
func (r *UserDocumentRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.UserDocument{ID: id}).Error
}

// DeleteByUserDocID 通过user_id+doc_id软删除
func (r *UserDocumentRepo) DeleteByUserDocID(ctx context.Context, userID int64, docID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND doc_id = ?", userID, docID).
		Delete(&model.UserDocument{}).Error
}

// ListByUserID 用户文档分页列表，page从1开始，pageSize每页条数
func (r *UserDocumentRepo) ListByUserID(ctx context.Context, userID int64, page int, pageSize int) ([]model.UserDocument, int64, error) {
	var list []model.UserDocument
	var total int64

	query := r.db.WithContext(ctx).Model(&model.UserDocument{}).Where("user_id = ?", userID)

	// 总条数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Limit(pageSize).Offset(offset).Find(&list).Error
	return list, total, err
}

// UnscopedGetByID 查询包含已软删除的数据（恢复/排查数据用）
func (r *UserDocumentRepo) UnscopedGetByID(ctx context.Context, id uint64) (*model.UserDocument, error) {
	var res *model.UserDocument
	err := r.db.WithContext(ctx).Unscoped().First(&res, id).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
