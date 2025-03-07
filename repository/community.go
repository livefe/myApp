package repository

import (
	"myApp/model"

	"gorm.io/gorm"
)

type CommunityRepository interface {
	Create(community *model.Community) error
	FindByID(id uint) (*model.Community, error)
	FindAll() ([]model.Community, error)
	AddMember(member *model.CommunityMember) error
	IncrementMemberCount(communityID uint, increment int) error
	FindMembersByCommunity(communityID uint) ([]model.CommunityMember, error)
}

type communityRepository struct {}

func NewCommunityRepository() CommunityRepository {
	return &communityRepository{}
}

func (r *communityRepository) Create(community *model.Community) error {
	return model.GetDB().Create(community).Error
}

func (r *communityRepository) FindByID(id uint) (*model.Community, error) {
	var community model.Community
	if err := model.GetDB().First(&community, id).Error; err != nil {
		return nil, err
	}
	return &community, nil
}

func (r *communityRepository) FindAll() ([]model.Community, error) {
	var communities []model.Community
	if err := model.GetDB().Find(&communities).Error; err != nil {
		return nil, err
	}
	return communities, nil
}

func (r *communityRepository) AddMember(member *model.CommunityMember) error {
	return model.GetDB().Create(member).Error
}

func (r *communityRepository) IncrementMemberCount(communityID uint, increment int) error {
	return model.GetDB().Model(&model.Community{}).Where("id = ?", communityID).UpdateColumn("members_count", gorm.Expr("members_count + ?", increment)).Error
}

func (r *communityRepository) FindMembersByCommunity(communityID uint) ([]model.CommunityMember, error) {
	var members []model.CommunityMember
	if err := model.GetDB().Where("community_id = ?", communityID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}
