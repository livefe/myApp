package service

import (
	"myApp/model"
	"myApp/repository"
	"time"
)

type ViewingService interface {
	CreateViewing(viewing *model.Viewing) error
	GetViewingByID(id uint) (*model.Viewing, error)
	GetAllViewings(params map[string]interface{}) ([]model.Viewing, error)
	UpdateViewing(viewing *model.Viewing) error
	DeleteViewing(id uint) error
	GetViewingsByUserID(userID uint) ([]model.Viewing, error)
	GetViewingsByHouseID(houseID uint) ([]model.Viewing, error)
	ConfirmViewing(id uint) error
	CompleteViewing(id uint) error
	CancelViewing(id uint, reason string) error
}

type viewingService struct {
	repo repository.ViewingRepository
}

func NewViewingService(repo repository.ViewingRepository) ViewingService {
	return &viewingService{repo: repo}
}

func (s *viewingService) CreateViewing(viewing *model.Viewing) error {
	// 设置初始状态为待确认
	viewing.Status = model.ViewingPending
	return s.repo.Create(viewing)
}

func (s *viewingService) GetViewingByID(id uint) (*model.Viewing, error) {
	return s.repo.GetByID(id)
}

func (s *viewingService) GetAllViewings(params map[string]interface{}) ([]model.Viewing, error) {
	return s.repo.GetAll(params)
}

func (s *viewingService) UpdateViewing(viewing *model.Viewing) error {
	return s.repo.Update(viewing)
}

func (s *viewingService) DeleteViewing(id uint) error {
	return s.repo.Delete(id)
}

func (s *viewingService) GetViewingsByUserID(userID uint) ([]model.Viewing, error) {
	return s.repo.GetViewingsByUserID(userID)
}

func (s *viewingService) GetViewingsByHouseID(houseID uint) ([]model.Viewing, error) {
	return s.repo.GetViewingsByHouseID(houseID)
}

func (s *viewingService) ConfirmViewing(id uint) error {
	// 获取预约看房记录
	viewing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	
	// 更新状态为已确认
	viewing.Status = model.ViewingConfirmed
	
	// 设置确认时间
	now := time.Now()
	viewing.ConfirmTime = &now
	
	return s.repo.Update(viewing)
}

func (s *viewingService) CompleteViewing(id uint) error {
	// 获取预约看房记录
	viewing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	
	// 更新状态为已完成
	viewing.Status = model.ViewingCompleted
	
	return s.repo.Update(viewing)
}

func (s *viewingService) CancelViewing(id uint, reason string) error {
	// 获取预约看房记录
	viewing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	
	// 更新状态为已取消
	viewing.Status = model.ViewingCancelled
	
	// 设置取消时间和原因
	now := time.Now()
	viewing.CancelTime = &now
	viewing.CancelReason = reason
	
	return s.repo.Update(viewing)
}