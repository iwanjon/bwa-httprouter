package campaign

import (
	"bwahttprouter/exception"
	"bwahttprouter/helper"
	"context"
	"errors"
	"fmt"
)

// 	"github.com/gosimple/slug"
// )

type Service interface {
	FindCampaigns(ctx context.Context, userId int) ([]Campaign, error)
	GetCampaignById(ctx context.Context, input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(ctx context.Context, input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(ctx context.Context, inputID GetCampaignDetailInput, inputparam CreateCampaignInput) (Campaign, error)
	SaveCampaignImage(ctx context.Context, input CreateCampaignImageInput, filelocation string) (CampaignImage, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) FindCampaigns(ctx context.Context, userId int) ([]Campaign, error) {
	if userId != 0 {
		campaigns, err := s.repo.FindByUserId(ctx, userId)
		helper.PanicIfError(err, "error in find campaign by user id service")

		return campaigns, nil
	}
	campaigns, err := s.repo.FindAll(ctx)
	helper.PanicIfError(err, "error in find campaign by user id service")

	return campaigns, nil

}

func (s *service) GetCampaignById(ctx context.Context, input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repo.FindById(ctx, input.ID)
	helper.PanicIfError(err, " error in getcampaign by id service")

	return campaign, nil

}

func (s *service) CreateCampaign(ctx context.Context, input CreateCampaignInput) (Campaign, error) {
	var campaigninput Campaign

	campaigninput.Name = input.Name
	campaigninput.BackerCount = 0
	campaigninput.Description = input.Description
	campaigninput.GoalAmount = input.GoalAmount
	campaigninput.Perks = input.Perks
	campaigninput.ShortDescription = input.ShortDescription
	campaigninput.UserID = input.User.ID
	campaigninput.CurrentAmount = 0
	slug := fmt.Sprintf("%d-%s", input.User.ID, input.Name)
	campaigninput.Slug = slug

	campaign, err := s.repo.SaveCampaign(ctx, campaigninput)
	helper.PanicIfError(err, "error in create campaign service")
	return campaign, nil
}
func (s *service) UpdateCampaign(ctx context.Context, inputID GetCampaignDetailInput, inputparam CreateCampaignInput) (Campaign, error) {
	selectedCampaign, err := s.repo.FindById(ctx, inputID.ID)
	if selectedCampaign.UserID != inputparam.User.ID {
		exception.PanicIfNotFound(errors.New("invalid user id"), "invalid user id, update campaign service")
	}
	exception.PanicIfNotFound(err, "error not found id of campaign service update campaign")
	selectedCampaign.Name = inputparam.Name
	selectedCampaign.Description = inputparam.Description
	selectedCampaign.ShortDescription = inputparam.ShortDescription
	selectedCampaign.GoalAmount = inputparam.GoalAmount
	selectedCampaign.Perks = inputparam.Perks

	updatedCampaign, err := s.repo.UpdateCampaign(ctx, selectedCampaign)
	helper.PanicIfError(err, "error in updtae campaign service")
	return updatedCampaign, nil

}
func (s *service) SaveCampaignImage(ctx context.Context, input CreateCampaignImageInput, filelocation string) (CampaignImage, error) {
	selectedCamaign, err := s.repo.FindById(ctx, input.CampaignID)
	exception.PanicIfNotFound(err, "not found campaign save image service")
	if selectedCamaign.UserID != input.User.ID {
		exception.PanicIfNotFound(errors.New("invalid user id"), "invalid user id, save campaign service")
	}
	is_primary := 0
	if input.IsPrimary {
		_, err := s.repo.MarkAllImagesAsNonPrimary(ctx, input.CampaignID)
		exception.PanicIfNotFound(err, "error in marl all image to zero status service")
		is_primary = 1
	}
	campaignImage := CampaignImage{}
	campaignImage.FileName = filelocation
	campaignImage.IsPrimary = is_primary
	campaignImage.CampaignID = input.CampaignID
	newcampaingimages, err := s.repo.SaveImage(ctx, campaignImage)
	helper.PanicIfError(err, "errro when saving campaign images , campaign service")

	return newcampaingimages, nil

}

// func (s *service) SaveCampaignImage(input CreateCampaignImageInput, filelocation string) (CampaignImage, error) {
// 	campaign, err := s.repo.FindById(input.CampaignID)
// 	if err != nil {
// 		return CampaignImage{}, err
// 	}

// 	if campaign.UserID != input.User.ID {
// 		return CampaignImage{}, errors.New("error not owner of the campaign")
// 	}

// 	isPrimary := 0
// 	if input.IsPrimary {
// 		_, err := s.repo.MarkAllImagesAsNonPrimary(input.CampaignID)
// 		if err != nil {
// 			return CampaignImage{}, err
// 		}
// 		isPrimary = 1
// 	}
// 	campaignImage := CampaignImage{}
// 	campaignImage.CampaignID = input.CampaignID

// 	campaignImage.IsPrimary = isPrimary

// 	campaignImage.FileName = filelocation

// 	newCAmpaignImage, err := s.repo.SaveImage(campaignImage)

// 	if err != nil {
// 		return newCAmpaignImage, err
// 	}

// 	return newCAmpaignImage, nil

// }

// func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, input CreateCampaignInput) (Campaign, error) {
// 	campaign, err := s.repo.FindById(inputID.ID)
// 	if err != nil {
// 		return campaign, err
// 	}

// 	if campaign.UserID != input.User.ID {
// 		return campaign, errors.New("error not owner of the campaign")
// 	}

// 	campaign.Name = input.Name
// 	campaign.ShortDescription = input.ShortDescription
// 	campaign.Description = input.Description
// 	campaign.Perks = input.Perks
// 	campaign.GoalAmount = input.GoalAmount

// 	updated, err := s.repo.Update(campaign)
// 	if err != nil {
// 		return updated, err
// 	}

// 	return updated, nil
// }

// func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
// 	campaign := Campaign{}
// 	campaign.Name = input.Name
// 	campaign.ShortDescription = input.ShortDescription
// 	campaign.Description = input.Description
// 	campaign.Perks = input.Perks
// 	campaign.GoalAmount = input.GoalAmount
// 	campaign.UserID = input.User.ID

// 	preslug := fmt.Sprintf("%s-%d", input.Name, input.User.ID)
// 	campaign.Slug = slug.Make(preslug)
// 	newcampaign, err := s.repo.Save(campaign)
// 	if err != nil {
// 		return newcampaign, err
// 	}
// 	return newcampaign, nil
// }

// func (s *service) FindCampaigns(userId int) ([]Campaign, error) {
// 	if userId != 0 {
// 		campains, err := s.repo.FindByUserId(userId)
// 		if err != nil {
// 			return campains, err
// 		}
// 		return campains, nil
// 	}
// 	campaigs, err := s.repo.FindAll()
// 	if err != nil {
// 		return campaigs, err
// 	}
// 	return campaigs, nil
// }

// func (s *service) GetCampaignById(input GetCampaignDetailInput) (Campaign, error) {
// 	var campaign Campaign
// 	campaign, err := s.repo.FindById(input.ID)
// 	if err != nil {
// 		return campaign, err
// 	}

// 	return campaign, nil
// }
