package services

import (
	"time"

	"github.com/carlogy/WorkoutBuilder/internal/database"
	id "github.com/google/uuid"
)

type User struct {
	ID         id.UUID    `json:"id"`
	FirstName  *string    `json:"firstName"`
	LastName   *string    `json:"lastName"`
	Email      string     `json:"email"`
	CreatedAt  *time.Time `json:"createdAt"`
	ModifiedAt *time.Time `json:"modifiedAt"`
}

func ConvertDBUserToUser(u database.CreateUserRow) User {
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

func ConvertFullDBUserToUser(u database.User) User {

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
