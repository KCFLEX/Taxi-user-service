package entity

import "time"

type Wallet struct {
	ID               int        // Unique ID for the wallet
	WalletType       string     // Type of wallet ('personal' or 'family')
	Balance          float64    // Balance in the wallet
	OwnerID          int        // User ID of the main owner (if applicable)
	PersonalWalletID *int       // ID of the personal wallet that this family wallet is attached to (nullable)
	CreatedAt        time.Time  // When the wallet was created
	UpdatedAt        time.Time  // When the wallet was last updated
	DeletedAt        *time.Time // Soft delete field (nullable)
}

type FamilyWalletMember struct {
	ID        int // Optional unique ID
	UserID    int /*Maps to user_id (foreign key to users.id)*/
	WalletID  int // Maps to wallet_id (foreign key to wallets.id)
	CreatedAt time.Time
}
