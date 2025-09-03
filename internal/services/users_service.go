package services

import (
	"context"
	"fmt"
	"time"

	"github.com/carlogy/WorkoutBuilder/internal/auth"
	"github.com/carlogy/WorkoutBuilder/internal/database"
	"github.com/carlogy/WorkoutBuilder/internal/repositories"
	id "github.com/google/uuid"
)

type UserService struct {
	userRepo repositories.UserRepository
	Secret   string
}

type User struct {
	ID            id.UUID    `json:"id"`
	FirstName     *string    `json:"firstName"`
	LastName      *string    `json:"lastName"`
	Email         string     `json:"email"`
	Token         *string    `json:"token"`
	Refresh_Token *string    `json:"refresh_token"`
	CreatedAt     *time.Time `json:"createdAt"`
	ModifiedAt    *time.Time `json:"modifiedAt"`
}

type UserRequestParams struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
}

type EmailAuthenticaton struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (us *UserService) createWorkoutParams(r UserRequestParams) database.CreateUserParams {
	return database.CreateUserParams{
		FirstName: NoneNullToNullString(r.FirstName),
		LastName:  NoneNullToNullString(r.LastName),
		Email:     r.Email,
		Password:  r.Password,
	}
}

func NewUserService(repo repositories.UserRepository, secret string) *UserService {
	return &UserService{userRepo: repo, Secret: secret}
}

func (us *UserService) CreateUser(ctx context.Context, r UserRequestParams) (User, error) {

	hashPW, err := auth.HashPassword(r.Password)
	if err != nil {
		return User{}, err
	}

	dbUserParams := us.createWorkoutParams(UserRequestParams{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Password:  hashPW,
	})

	dbUser, err := us.userRepo.CreateDBUser(ctx, dbUserParams)
	if err != nil {
		return User{}, err
	}

	return us.ConvertDBUserToUser(dbUser), nil
}

func (us *UserService) DeleteUserByID(ctx context.Context, userId id.UUID, token string, as *AuthService) (User, error) {

	err := as.ValidateJWT(token, userId)
	if err != nil {
		return User{}, err
	}

	dbDeletedUser, err := us.userRepo.DeleteDBUserById(ctx, userId)
	if err != nil {
		return User{}, err
	}

	return ConvertDBDeleteUserToUser(dbDeletedUser), nil
}

func (us *UserService) UpdateUserById(ctx context.Context, ur UserRequestParams, token string, id id.UUID) (User, error) {

	hashedPW, err := auth.HashPassword(ur.Password)
	if err != nil {
		fmt.Println("Error hashing password: ", err)
		return User{}, err
	}

	updatedUser, err := us.userRepo.UpdateDBUserById(ctx, database.UpdateUserByIdParams{
		FirstName: NoneNullToNullString(ur.FirstName),
		LastName:  NoneNullToNullString(ur.LastName),
		Email:     ur.Email,
		Password:  hashedPW,
		ID:        id,
	})
	if err != nil {
		fmt.Println("Error updating user: ", err)
		return User{}, err
	}

	return ConvertDBUpdateUserToUser(updatedUser), nil
}

func (us *UserService) ConvertDBUserToUser(u database.CreateUserRow) User {
	user := User{
		ID:         u.ID,
		FirstName:  NullStringToString(u.FirstName),
		LastName:   NullStringToString(u.LastName),
		Email:      u.Email,
		CreatedAt:  NullTimeToTime(u.CreatedAt),
		ModifiedAt: NullTimeToTime(u.ModifiedAt),
	}
	return user
}

func ConvertDBUpdateUserToUser(u database.UpdateUserByIdRow) User {
	user := User{
		ID:         u.ID,
		FirstName:  NullStringToString(u.FirstName),
		LastName:   NullStringToString(u.LastName),
		Email:      u.Email,
		CreatedAt:  NullTimeToTime(u.CreatedAt),
		ModifiedAt: NullTimeToTime(u.ModifiedAt),
	}
	return user
}

func ConvertDBDeleteUserToUser(u database.DeleteUserByIdRow) User {
	user := User{
		ID:         u.ID,
		FirstName:  NullStringToString(u.FirstName),
		LastName:   NullStringToString(u.LastName),
		Email:      u.Email,
		CreatedAt:  NullTimeToTime(u.CreatedAt),
		ModifiedAt: NullTimeToTime(u.ModifiedAt),
	}
	return user
}

func ConvertFullDBUserToUser(u database.User, t *string, rt *string) User {

	user := User{
		ID:            u.ID,
		FirstName:     NullStringToString(u.FirstName),
		LastName:      NullStringToString(u.LastName),
		Email:         u.Email,
		CreatedAt:     NullTimeToTime(u.CreatedAt),
		ModifiedAt:    NullTimeToTime(u.ModifiedAt),
		Token:         t,
		Refresh_Token: rt,
	}
	return user
}
