package route

import (
	"database/sql"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-cqrs-example/internal/command"
	"github.com/ytake/protoactor-go-cqrs-example/internal/database/mysql"
	"github.com/ytake/protoactor-go-cqrs-example/internal/registration"
	persistencemysql "github.com/ytake/protoactor-go-persistence-mysql"
)

type (
	// CreateUserMessageHandler is a interface to handle CreateUser message
	CreateUserMessageHandler interface {
		// Handle is a method to handle CreateUser message
		Handle(ctx actor.Context, msg *command.CreateUser)
	}
	Actor struct {
		system *actor.ActorSystem
		pid    *actor.PID
	}
	RestAPI struct {
		createUserMessageHandler CreateUserMessageHandler
	}
)

// ActorSystem is a sample
func (a *Actor) ActorSystem() *actor.ActorSystem {
	return a.system
}

// PID is a pid
func (a *Actor) PID() *actor.PID {
	return a.pid
}

// NewRestAPI is a constructor for RestAPI
func NewRestAPI(createUserMessageHandler CreateUserMessageHandler) actor.Actor {
	return &RestAPI{
		createUserMessageHandler: createUserMessageHandler,
	}
}

// NewRestAPIActorSystem is a function to create a new actor system
func NewRestAPIActorSystem(db *sql.DB) (*Actor, error) {
	system := actor.NewActorSystem()
	provider, err := persistencemysql.New(3, persistencemysql.NewTable(), db, system.Logger())
	if err != nil {
		return nil, err
	}
	root, err := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
		return NewRestAPI(
			registration.NewCreateUser(registration.NewUserModelUpdate(mysql.NewUserStore(db)), provider))
	}), "rest-api")
	if err != nil {
		return nil, err
	}
	return &Actor{
		system: system,
		pid:    root,
	}, nil
}

func (a *RestAPI) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *command.CreateUser:
		a.createUserMessageHandler.Handle(ctx, msg)
	}
}
