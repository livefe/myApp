package service

import (
	"myApp/model"
	"myApp/repository"
)

type FavoriteService interface {
	AddFavorite(favorite *model.Favorite) error
	RemoveFavorite(id uint) error
	GetFavoriteByID(id uint) (*model.Favorite, error)
	GetUserFavorites(userID uint) ([]model.Favorite, error)
	IsFavorite(userID, houseID uint) (bool, error)
	ToggleFavorite(userID, houseID uint, notes string) error
}

type favoriteService struct {
	repo repository.FavoriteRepository
}

func NewFavoriteService(repo repository.FavoriteRepository) FavoriteService {
	return &favoriteService{repo: repo}
}

func (s *favoriteService) AddFavorite(favorite *model.Favorite) error {
	return s.repo.Create(favorite)
}

func (s *favoriteService) RemoveFavorite(id uint) error {
	return s.repo.Delete(id)
}

func (s *favoriteService) GetFavoriteByID(id uint) (*model.Favorite, error) {
	return s.repo.GetByID(id)
}

func (s *favoriteService) GetUserFavorites(userID uint) ([]model.Favorite, error) {
	return s.repo.GetFavoritesByUserID(userID)
}

func (s *favoriteService) IsFavorite(userID, houseID uint) (bool, error) {
	return s.repo.IsFavorite(userID, houseID)
}

func (s *favoriteService) ToggleFavorite(userID, houseID uint, notes string) error {
	// 检查是否已收藏
	isFav, err := s.repo.IsFavorite(userID, houseID)
	if err != nil {
		return err
	}
	
	// 如果已收藏，则取消收藏
	if isFav {
		return s.repo.DeleteByUserAndHouse(userID, houseID)
	}
	
	// 如果未收藏，则添加收藏
	favorite := &model.Favorite{
		UserID:  userID,
		HouseID: houseID,
		Notes:   notes,
	}
	return s.repo.Create(favorite)
}