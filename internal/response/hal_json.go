package response

import (
	"time"

	"github.com/nvellon/hal"
)

type UserCounter struct {
	Count int
	Total int
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (c User) GetMap() hal.Entry {
	return hal.Entry{
		"id":         c.ID,
		"name":       c.Name,
		"email":      c.Email,
		"created_at": c.CreatedAt,
	}
}

func (u UserCounter) GetMap() hal.Entry {
	return hal.Entry{
		"count": u.Count,
		"total": u.Total,
	}
}
