package service

import (
	"errors"
	"myApp/model"
	"myApp/repository"
)

type LandlordService interface {
	CreateLandlord(landlord *model.Landlord) error
	GetLandlordByID(id uint) (*model.Landlord, error)
	GetLandlordByUserID(userID uint) (*model.Landlord, error)
	UpdateLandlord(landlord *model.Landlord) error
	DeleteLandlord(id uint) error
	VerifyLandlord(id uint) error
}

type landlordService struct {
	repo repository.LandlordRepository
	userRepo repository.UserRepository
}

func NewLandlordService(repo repository.LandlordRepository, userRepo repository.UserRepository) LandlordService {
	return &landlordService{repo: repo, userRepo: userRepo}
}

func (s *landlordService) CreateLandlord(landlord *model.Landlord) error {
	// 检查用户是否存在
	user, err := s.userRepo.FindByID(landlord.UserID)
	if err != nil {
		return errors.New("用户不存在")
	}
	
	// 检查用户是否已经是房东
	existingLandlord, _ := s.repo.FindByUserID(landlord.UserID)
	if existingLandlord != nil {
		return errors.New("该用户已经是房东")
	}
	
	// 更新用户类型为房东
	user.UserType = 1 // 1-房东
	if err := s.userRepo.Update(user); err != nil {
		return err
	}
	
	return s.repo.Create(landlord)
}

func (s *landlordService) GetLandlordByID(id uint) (*model.Landlord, error) {
	return s.repo.FindByID(id)
}

func (s *landlordService) GetLandlordByUserID(userID uint) (*model.Landlord, error) {
	return s.repo.FindByUserID(userID)
}

func (s *landlordService) UpdateLandlord(landlord *model.Landlord) error {
	// 检查房东是否存在
	existingLandlord, err := s.repo.FindByID(landlord.ID)
	if err != nil {
		return errors.New("房东不存在")
	}
	
	// 保持用户ID不变
	landlord.UserID = existingLandlord.UserID
	
	return s.repo.Update(landlord)
}

func (s *landlordService) DeleteLandlord(id uint) error {
	// 获取房东信息
	landlord, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("房东不存在")
	}
	
	// 更新用户类型为普通用户
	user, err := s.userRepo.FindByID(landlord.UserID)
	if err == nil {
		user.UserType = 0 // 0-普通用户
		s.userRepo.Update(user)
	}
	
	return s.repo.Delete(id)
}

func (s *landlordService) VerifyLandlord(id uint) error {
	// 获取房东信息
	landlord, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("房东不存在")
	}
	
	// 更新认证状态
	landlord.Verified = true
	return s.repo.Update(landlord)
}