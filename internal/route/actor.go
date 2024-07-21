package route

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/ytake/protoactor-go-cqrs-example/internal/command"
	"github.com/ytake/protoactor-go-cqrs-example/internal/database/mysql"
	"github.com/ytake/protoactor-go-cqrs-example/internal/message"
	"github.com/ytake/protoactor-go-cqrs-example/internal/registration"
	persistencemysql "github.com/ytake/protoactor-go-persistence-mysql"
)

type (
	Actor struct {
		system *actor.ActorSystem
		pid    *actor.PID
	}
	RestAPI struct {
		db       mysql.RegistrationUserExecutor
		provider persistence.Provider
		rmu      *actor.PID
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
func NewRestAPI(db mysql.RegistrationUserExecutor, provider persistence.Provider) actor.Actor {
	return &RestAPI{
		db:       db,
		provider: provider,
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
		return NewRestAPI(mysql.NewUserStore(db), provider)
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
	case *actor.Started:
		a.rmu = ctx.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return registration.NewUserModelUpdate(a.db)
		}))
	case *command.CreateUser:
		ref, err := ctx.SpawnNamed(
			actor.PropsFromProducer(func() actor.Actor {
				return registration.NewUser(msg.Stream, a.rmu)
			}, actor.WithReceiverMiddleware(persistence.Using(a.provider))), "user-"+msg.Email)
		// 登録ユーザーのメールアドレスが既に存在する場合はエラーを返す
		// メッセージ送信時に現在のバージョンを送信することで、永続化されたデータとの競合を防ぐことができます
		// 詳しくはprotobufを参照してください
		if errors.Is(err, actor.ErrNameExists) {
			ctx.Send(msg.Stream, &message.UserCreateError{Message: fmt.Sprintf("user %s already exists", msg.Email)})
			return
		}
		if err != nil {
			ctx.Send(msg.Stream, &message.UserCreateError{Message: fmt.Sprintf("failed error %s", err.Error())})
			return
		}
		ctx.Send(ref, msg)
	}
}
