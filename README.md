# Taxi-user-service

This repository contains the backend and frontend implementation of a User Service API that manages user registration, authentication, taxi ordering, wallet management, and transaction processing. The system uses PostgreSQL for data persistence, Redis for caching.

## Features

### User Authentication

Sign up (name, phone number, email, password).
Sign in (phone number, password).
Logout (invalidate user token via Redis)

### Profile Management

View profile details.
Update profile information.
Soft delete profile.

### Taxi Ordering

Order a taxi with the following fields: taxi type, from, to.
System seeks an available driver. If no driver is found, user must wait a configurable amount of time before receiving a rejection message.
Rate the most recent trip with a score from 1 to 5 and an optional comment.

### Wallet Management

Create personal and family wallets.
Link family wallets to a personal wallet.
Add family members by phone number.
Top up personal and family wallets (only the owner can top up a family wallet).
Choose which wallet to use for payments.
View wallet transaction history (only the family wallet owner can view family wallet transactions).

### Transaction Management

Transactions have four statuses: create, blocked, success, canceled.
When an order is created, the transaction status is create.
If the wallet balance is sufficient, the status is blocked; otherwise, it's canceled.
Upon completion of a trip, the transaction status changes to success.
