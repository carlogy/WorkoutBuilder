package repositories

import (
	"context"

	db "github.com/carlogy/WorkoutBuilder/internal/database"
	id "github.com/google/uuid"
)

type UserRepository struct {
	db *db.Queries
}

// type UserRepository interface {
// 	CreateDBUser(ctx context.Context, user db.CreateUserParams) (db.CreateUserRow, error)
// 	UpdateDBUserById(ctx context.Context, user db.UpdateUserByIdParams) (db.UpdateUserByIdRow, error)
// 	DeleteDBUserById(ctx context.Context, id id.UUID) (db.DeleteUserByIdRow, error)
// }

func NewUserRepository(db *db.Queries) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) CreateDBUser(ctx context.Context, user db.CreateUserParams) (db.CreateUserRow, error) {

	createdUser, err := ur.db.CreateUser(ctx, user)
	if err != nil {
		return db.CreateUserRow{}, err
	}

	return createdUser, nil
}

func (ur *UserRepository) UpdateDBUserById(ctx context.Context, user db.UpdateUserByIdParams) (db.UpdateUserByIdRow, error) {

	updatedUser, err := ur.db.UpdateUserById(ctx, user)
	if err != nil {
		return db.UpdateUserByIdRow{}, err
	}

	return updatedUser, nil
}

func (ur *UserRepository) DeleteDBUserById(ctx context.Context, id id.UUID) (db.DeleteUserByIdRow, error) {

	deletedUser, err := ur.db.DeleteUserById(ctx, id)
	if err != nil {
		return db.DeleteUserByIdRow{}, err
	}

	return deletedUser, nil
}

func (ur *UserRepository) GetDBUserByEmail(ctx context.Context, email string) (db.User, error) {

	dbUser, err := ur.db.GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, err
	}

	return dbUser, nil
}
