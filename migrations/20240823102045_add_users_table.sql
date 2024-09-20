-- +goose Up

-- SQL in section 'Up' is executed when this migration is applied
CREATE TYPE wallet_type AS ENUM ('personal', 'family');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(60) NOT NULL,
    phone VARCHAR(60) UNIQUE NOT NULL,
    email VARCHAR(60) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    rating DECIMAL(3,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    type wallet_type NOT NULL,
    balance DECIMAL(9,2) NOT NULL,
    main_owner_id INTEGER,
    personal_wallet_id INTEGER, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (main_owner_id) REFERENCES users(id),
    FOREIGN KEY (personal_wallet_id) REFERENCES wallets(id)
);

-- The many-to-many relationship created by the family_wallet_members table allows:
-- Users to belong to multiple family wallets at the same time.
-- Family wallets to have multiple users (members), making it a flexible system for managing shared wallets.
CREATE TABLE family_wallet_members (
    id SERIAL PRIMARY KEY,       -- Auto-incrementing unique ID for each row
    user_id INTEGER,             -- Foreign key referencing users
    wallet_id INTEGER,           -- Foreign key referencing wallets
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Track when the user was added
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (wallet_id) REFERENCES wallets(id)
);

CREATE TABLE trips (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    driver_id VARCHAR(60) NOT NULL,
    taxi_type VARCHAR(60) NOT NULL,
    from_location VARCHAR(100) NOT NULL,
    to_location VARCHAR(100) NOT NULL,
    rating DECIMAL(3,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    amount DECIMAL(9,2) NOT NULL,
    wallet_id INTEGER,
    trip_id INTEGER,
    FOREIGN KEY (wallet_id) REFERENCES wallets(id),
    FOREIGN KEY (trip_id) REFERENCES trips(id)
);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

DROP TABLE transactions;
DROP TABLE trips;
DROP TABLE family_wallet_members;
DROP TABLE wallets;
DROP TABLE users;
DROP TYPE wallet_type;

