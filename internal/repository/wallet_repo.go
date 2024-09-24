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

func (repo *Repository) GetFamilyWalletByOwnerID(ctx context.Context, userID int, walletType string) (int, error) {
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

func (repo *Repository) DeductAmountFromWallet(ctx context.Context, walletID, amount int) error {
	query := `UPDATE wallets
			  SET balance = balance - $1, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $2 AND balance >= $1`
	_, err := repo.db.ExecContext(ctx, query, amount, walletID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) GetAllUserWallets(ctx context.Context, userID int) ([]entity.Wallet, error) {
	query := `
	WITH owned_wallets AS (
		SELECT id, type, balance, main_owner_id, personal_wallet_id, created_at, updated_at, deleted_at
		FROM wallets
		WHERE main_owner_id = $1
	),
	family_member_wallets AS (
		SELECT w.id, w.type, w.balance, w.main_owner_id, w.personal_wallet_id, w.created_at, w.updated_at, w.deleted_at
		FROM family_wallet_members fwm
		JOIN wallets w ON fwm.wallet_id = w.id
		WHERE fwm.user_id = $1
	)
	SELECT * FROM owned_wallets
	UNION
	SELECT * FROM family_member_wallets;
	`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []entity.Wallet

	for rows.Next() {
		var wallet entity.Wallet
		err := rows.Scan(&wallet.ID, &wallet.WalletType, &wallet.Balance, &wallet.OwnerID, &wallet.PersonalWalletID, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return wallets, nil
}
