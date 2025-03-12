package service

import (
	"myApp/model"
	"myApp/repository"
)

type HouseService interface {
	CreateHouse(house *model.House) error
	GetHouseByID(id uint) (*model.House, error)
	GetAllHouses(params map[string]interface{}) ([]model.House, error)
	UpdateHouse(house *model.House) error
	DeleteHouse(id uint) error
	GetHousesByLandlordID(landlordID uint) ([]model.House, error)
	IncrementViewCount(id uint) error
}

type houseService struct {
	repo repository.HouseRepository
}

func NewHouseService(repo repository.HouseRepository) HouseService {
	return &houseService{repo: repo}
}

func (s *houseService) CreateHouse(house *model.House) error {
	return s.repo.Create(house)
}

func (s *houseService) GetHouseByID(id uint) (*model.House, error) {
	return s.repo.GetByID(id)
}

func (s *houseService) GetAllHouses(params map[string]interface{}) ([]model.House, error) {
	return s.repo.GetAll(params)
}

func (s *houseService) UpdateHouse(house *model.House) error {
	return s.repo.Update(house)
}

func (s *houseService) DeleteHouse(id uint) error {
	return s.repo.Delete(id)
}

func (s *houseService) GetHousesByLandlordID(landlordID uint) ([]model.House, error) {
	return s.repo.GetHousesByLandlordID(landlordID)
}

func (s *houseService) IncrementViewCount(id uint) error {
	return s.repo.IncrementViewCount(id)
}