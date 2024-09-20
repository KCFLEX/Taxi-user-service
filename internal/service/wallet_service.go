package services

import (
	"context"
	"errors"

	"github.com/KCFLEX/Taxi-user-service/errorpac"
	"github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
	"github.com/KCFLEX/Taxi-user-service/internal/repository/entity"
)

func (srv *Service) AddPersonalWallet(ctx context.Context, userID int, walletInfo models.Wallet) error {

	newWallet := entity.Wallet{
		WalletType: walletInfo.WalletType,
		Balance:    walletInfo.Balance,
		OwnerID:    userID,
	}

	err := srv.repo.CreatePersonalWallet(ctx, newWallet)
	if err != nil {
		return &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errors.New("failed to create new wallet"),
		}
	}

	return nil
}

func (srv *Service) AddFamilyWallet(ctx context.Context, userID int, walletInfo models.Wallet) error {
	// get wallet ID inorder to attach the family wallet to the owner personal wallet
	walletID, err := srv.repo.GetPersonalWalletBYID(ctx, userID)
	if err != nil {
		return &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errorpac.ErrRetrieveWalletIDFail,
		}
	}

	// now put the wallet id in the foreign key of the new family wallet we are going to create to signify ownership
	newFamilyWallet := entity.Wallet{
		WalletType:       walletInfo.WalletType,
		Balance:          walletInfo.Balance,
		OwnerID:          userID,
		PersonalWalletID: &walletID,
	}
	err = srv.repo.AddFamilyWallet(ctx, newFamilyWallet)
	if err != nil {
		return &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errors.New("failed to add Family wallet"),
		}
	}

	return nil
}

func (srv *Service) AddUserToFamilyByPhone(ctx context.Context, userID int, phone models.Phone) error {
	// checking if user exists
	checkPhone := entity.User{
		Phone: phone.PhoneNO,
	}

	getUser, err := srv.repo.UserPhoneExists(ctx, checkPhone)
	if err != nil {
		return err
	}

	// retrieving wallet id
	getWalletID, err := srv.repo.GetfamilyWalletByOwnerID(ctx, userID, "family")
	if err != nil {
		return &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errorpac.ErrRetrieveWalletIDFail,
		}
	}

	// now we add a member to the family wallet
	newFamilyMember := entity.FamilyWalletMember{
		UserID:   getUser.ID,
		WalletID: getWalletID,
	}

	err = srv.repo.AddUserToFamilyWallet(ctx, newFamilyMember)
	if err != nil {
		return &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errors.New("failed to add user to family wallet"),
		}
	}

	return nil

}
