package repository

import (
	"context"

	"github.com/KCFLEX/Taxi-user-service/internal/repository/entity"
)

func (repo *Repository) CreatePersonalWallet(ctx context.Context, walletInfo entity.Wallet) error {
	query := `INSERT INTO wallets (type, balance, main_owner_id) VALUES ($1, $2, $3)`

	_, err := repo.db.ExecContext(ctx, query, walletInfo.WalletType, walletInfo.Balance, walletInfo.OwnerID)
	if err != nil {
		return err
	}

	return nil

}

func (repo *Repository) GetPersonalWalletBYID(ctx context.Context, userID int) (int, error) {
	query := `SELECT  id FROM wallets WHERE main_owner_id = $1`

	var walletID int

	err := repo.db.QueryRowContext(ctx, query, userID).Scan(&walletID)

	if err != nil {
		return 0, err
	}

	return walletID, nil
}

func (repo *Repository) AddFamilyWallet(ctx context.Context, walletInfo entity.Wallet) error {
	query := `INSERT INTO wallets (type, balance, main_owner_id, personal_wallet_id) VALUES ($1, $2, $3, $4)`

	_, err := repo.db.ExecContext(ctx, query, walletInfo.WalletType, walletInfo.Balance, walletInfo.OwnerID, walletInfo.PersonalWalletID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) GetfamilyWalletByOwnerID(ctx context.Context, userID int, walletType string) (int, error) {
	query := `SELECT  id FROM wallets WHERE main_owner_id = $1 AND type = $2`
	var walletID int
	err := repo.db.QueryRowContext(ctx, query, userID, walletType).Scan(&walletID)
	if err != nil {
		return 0, err
	}

	return walletID, nil
}

// for adding user to family wallet
func (repo *Repository) AddUserToFamilyWallet(ctx context.Context, newMember entity.FamilyWalletMember) error {
	query := `INSERT INTO family_wallet_members (user_id, wallet_id) VALUES ($1, $2)`

	_, err := repo.db.ExecContext(ctx, query, newMember.UserID, newMember.WalletID)
	if err != nil {
		return err
	}

	return nil
}
