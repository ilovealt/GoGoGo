package response

import (
	"time"

	"gogogo/internal/model"
)

type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Role      string    `json:"role"`
	Status    uint8     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(user model.User) User {
	return User{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
