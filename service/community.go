package service

import (
	"myApp/model"
	"myApp/repository"
)

type CommunityService interface {
	CreateCommunity(community *model.Community) (*model.Community, error)
	AddMember(communityID uint, userID uint, role int) error
	GetCommunityByID(id uint) (*model.Community, error)
	GetCommunityMembers(communityID uint) ([]model.CommunityMember, error)
	GetAllCommunities() ([]model.Community, error)
}

type communityService struct {
	communityRepo repository.CommunityRepository
}

func NewCommunityService(repo repository.CommunityRepository) CommunityService {
	return &communityService{communityRepo: repo}
}

func (s *communityService) CreateCommunity(community *model.Community) (*model.Community, error) {
	community.MembersCount = 1
	community.Status = 1
	if err := s.communityRepo.Create(community); err != nil {
		return nil, err
	}

	// 自动添加创建者为管理员
	member := &model.CommunityMember{
		CommunityID: community.ID,
		UserID:      community.CreatorID,
		Role:        2, // 管理员角色
	}
	if err := s.communityRepo.AddMember(member); err != nil {
		return nil, err
	}
	return community, nil
}

func (s *communityService) AddMember(communityID uint, userID uint, role int) error {
	member := &model.CommunityMember{
		CommunityID: communityID,
		UserID:      userID,
		Role:        role,
	}

	if err := s.communityRepo.AddMember(member); err != nil {
		return err
	}

	// 更新成员数量
	return s.communityRepo.IncrementMemberCount(communityID, 1)
}

func (s *communityService) GetCommunityByID(id uint) (*model.Community, error) {
	return s.communityRepo.FindByID(id)
}

func (s *communityService) GetCommunityMembers(communityID uint) ([]model.CommunityMember, error) {
	return s.communityRepo.FindMembersByCommunity(communityID)
}

func (s *communityService) GetAllCommunities() ([]model.Community, error) {
	return s.communityRepo.FindAll()
}
