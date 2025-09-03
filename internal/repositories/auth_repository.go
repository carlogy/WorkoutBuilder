package repositories

import (
	"context"

	db "github.com/carlogy/WorkoutBuilder/internal/database"
)

type AuthRepository struct {
	db *db.Queries
}

func NewAuthRepository(db *db.Queries) *AuthRepository {
	return &AuthRepository{db: db}
}

func (ar *AuthRepository) GetDBUserByEmail(ctx context.Context, email string) (db.User, error) {

	dbUser, err := ar.db.GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, err
	}

	return dbUser, nil
}

func (ar *AuthRepository) StoreRefreshTokenDB(ctx context.Context, rt db.StoreRefreshTokenParams) (db.RefreshToken, error) {

	storeToken, err := ar.db.StoreRefreshToken(ctx, rt)
	if err != nil {
		return db.RefreshToken{}, err
	}

	return storeToken, nil
}

func (ar *AuthRepository) UpdateRefreshTokenDB(ctx context.Context, rt db.UpdateRefreshTokenParams) (db.RefreshToken, error) {

	dbRfToken, err := ar.db.UpdateRefreshToken(ctx, rt)
	if err != nil {
		return db.RefreshToken{}, err
	}
	return dbRfToken, nil
}

func (ar *AuthRepository) RevokeTokenDB(ctx context.Context, token string) (db.RefreshToken, error) {

	dbRT, err := ar.db.RevokeRefreshToken(ctx, token)
	if err != nil {
		return db.RefreshToken{}, err
	}
	return dbRT, nil
}
