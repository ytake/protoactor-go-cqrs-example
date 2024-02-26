package action

import (
	"net/http"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/stream"
	"github.com/labstack/echo/v4"
	"github.com/ytake/protoactor-go-cqrs-example/internal/command"
	"github.com/ytake/protoactor-go-cqrs-example/internal/message"
	"github.com/ytake/protoactor-go-cqrs-example/internal/route"
)

type (
	userInput struct {
		UserName string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
	}
	UserRegistration struct {
		system *actor.ActorSystem
		ref    *actor.PID
		stream *stream.TypedStream[message.UserCreateMessenger]
	}
)

// NewUserRegistration is a sample
func NewUserRegistration(ref *route.Actor) (*UserRegistration, error) {
	typedStream := stream.NewTypedStream[message.UserCreateMessenger](ref.ActorSystem())
	return &UserRegistration{
		system: ref.ActorSystem(),
		ref:    ref.PID(),
		stream: typedStream,
	}, nil
}

func (u *UserRegistration) Handle(c echo.Context) error {
	ui := new(userInput)
	if err := c.Bind(ui); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(ui); err != nil {
		return err
	}
	go func() {
		u.system.Root.Send(u.ref, &command.CreateUser{
			UserName: ui.UserName,
			Email:    ui.Email,
			Stream:   u.stream.PID(),
		})
	}()
	res := <-u.stream.C()
	if res.IsSuccess() {
		return c.JSON(http.StatusOK, res)
	}
	return c.JSON(http.StatusBadRequest, res)
}
