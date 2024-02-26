package action

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nvellon/hal"
	"github.com/ytake/protoactor-go-cqrs-example/internal/database/mysql"
	"github.com/ytake/protoactor-go-cqrs-example/internal/response"
)

type UserList struct {
	query mysql.RegistrationUserQueryExecutor
}

func NewUserList(query mysql.RegistrationUserQueryExecutor) *UserList {
	return &UserList{
		query: query,
	}
}

// Handle GraphQLやporotobufなど好きなフォーマットでご利用ください
// sampleではわかりやすいようにhal+jsonにしています
func (u *UserList) Handle(c echo.Context) error {
	users, err := u.query.GetRegistrationUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	hr := hal.NewResource(
		response.UserCounter{Count: len(users), Total: len(users)},
		c.Request().URL.String())
	for _, r := range users {
		hr.Embed("users",
			hal.NewResource(
				response.User{ID: r.ID, Name: r.Name, Email: r.Email, CreatedAt: r.CreatedAt.Time},
				"/users/1"))
	}
	c.Response().Header().Set(echo.HeaderContentType, "application/hal+json")
	return c.JSON(http.StatusOK, hr)
}
