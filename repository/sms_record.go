package repository

import (
	"myApp/model"

	"gorm.io/gorm"
)

// SMSRecordRepository 短信记录仓库接口
type SMSRecordRepository interface {
	Create(record *model.SMSRecord) error
	FindByPhone(phone string, limit, offset int) ([]*model.SMSRecord, error)
	CountByPhone(phone string) (int64, error)
}

// smsRecordRepository 短信记录仓库实现
type smsRecordRepository struct {
	db *gorm.DB
}

// NewSMSRecordRepository 创建短信记录仓库实例
func NewSMSRecordRepository() SMSRecordRepository {
	return &smsRecordRepository{
		db: model.GetDB(),
	}
}

// Create 创建短信记录
func (r *smsRecordRepository) Create(record *model.SMSRecord) error {
	return r.db.Create(record).Error
}

// FindByPhone 根据手机号查询短信记录
func (r *smsRecordRepository) FindByPhone(phone string, limit, offset int) ([]*model.SMSRecord, error) {
	var records []*model.SMSRecord
	query := r.db.Where("phone = ?", phone).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// CountByPhone 统计指定手机号的短信记录数量
func (r *smsRecordRepository) CountByPhone(phone string) (int64, error) {
	var count int64
	if err := r.db.Model(&model.SMSRecord{}).Where("phone = ?", phone).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
