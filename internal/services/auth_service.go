package services

import (
	"context"
	"errors"
	"fmt"

	id "github.com/google/uuid"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	db "github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/carlogy/WorkoutBuilder/internal/repositories"
)

type AuthService struct {
	authRepo repositories.AuthRepository
	Secret   string
}

type EmailAuthRequestParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatedRefreshTokenResponse struct {
	RefreshToken string `json:"token"`
}

func NewAuthService(repo repositories.AuthRepository, secret string) *AuthService {
	return &AuthService{authRepo: repo, Secret: secret}
}

func (as *AuthService) AuthenticateByEmail(ctx context.Context, r EmailAuthRequestParams) (User, error) {

	dbUser, err := as.authRepo.GetDBUserByEmail(ctx, r.Email)
	if err != nil {
		return User{}, err
	}

	err = auth.CheckPasswordHash(r.Password, dbUser.Password)
	if err != nil {
		return User{}, err
	}

	token, err := auth.MakeJWT(dbUser.ID, as.Secret)
	if err != nil {
		return User{}, err
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return User{}, err
	}

	tokenParams := db.StoreRefreshTokenParams{
		Token:  refreshToken,
		UserID: dbUser.ID,
	}
	_, err = as.authRepo.StoreRefreshTokenDB(ctx, tokenParams)
	if err != nil {
		return User{}, err
	}

	return ConvertFullDBUserToUser(dbUser, &token, &refreshToken), nil
}

func (as *AuthService) RefreshToken(ctx context.Context, oldToken string) (UpdatedRefreshTokenResponse, error) {

	newToken, err := auth.MakeRefreshToken()
	if err != nil {
		return UpdatedRefreshTokenResponse{}, err
	}

	dbToken, err := as.authRepo.UpdateRefreshTokenDB(ctx, db.UpdateRefreshTokenParams{
		Token:   newToken,
		Token_2: oldToken,
	})
	if err != nil {
		return UpdatedRefreshTokenResponse{}, err
	}
	if dbToken.Token != newToken {
		return UpdatedRefreshTokenResponse{}, errors.New("failed to update token")
	}

	return UpdatedRefreshTokenResponse{
		RefreshToken: dbToken.Token,
	}, nil

}

func (as *AuthService) RevokeToken(ctx context.Context, token string) (bool, error) {

	revokedToken, err := as.authRepo.RevokeTokenDB(ctx, token)
	if err != nil {
		return false, err
	}

	if !revokedToken.RevokedAt.Valid {
		fmt.Println("Token returned is not revoked:  ", revokedToken.RevokedAt)
		return false, errors.New("failed to revoke token")

	}

	return true, nil
}

func (as *AuthService) ValidateJWT(token string, userID id.UUID) error {

	tokenUserID, err := auth.ValidateJWT(token, as.Secret)
	if err != nil {
		fmt.Println("Error validating token: ", err)
		return err
	}

	if tokenUserID != userID {
		fmt.Println("Token user and userID do not match: ", err)
		return errors.New("invalid userId provided for update")
	}

	return nil
}
