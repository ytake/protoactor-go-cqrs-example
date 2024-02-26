package route

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/asynkron/protoactor-go/stream"
	"github.com/oklog/ulid/v2"
	"github.com/ytake/protoactor-go-cqrs-example/internal/command"
	"github.com/ytake/protoactor-go-cqrs-example/internal/database/mysql"
	"github.com/ytake/protoactor-go-cqrs-example/internal/message"
	"github.com/ytake/protoactor-go-cqrs-example/internal/registration"
	"github.com/ytake/protoactor-go-cqrs-example/pkg/event"
)

type Provider struct {
	providerState persistence.ProviderState
}

func NewProvider(snapshotInterval int) *Provider {
	return &Provider{
		providerState: persistence.NewInMemoryProvider(snapshotInterval),
	}
}

func (p *Provider) GetState() persistence.ProviderState {
	return p.providerState
}

func (p *Provider) InitState(actorName string, eventNum, eventIndexAfterSnapshot int) *Provider {
	for i := 0; i < eventNum; i++ {
		p.providerState.PersistEvent(
			actorName,
			i,
			&event.UserCreated{
				UserID:   ulid.Make().String(),
				UserName: "test" + strconv.Itoa(i),
				Email:    "test" + strconv.Itoa(i) + "@example.com",
			},
		)
	}
	p.providerState.PersistSnapshot(
		actorName,
		eventIndexAfterSnapshot,
		&event.UserCreated{
			UserID:   ulid.Make().String(),
			UserName: "test" + strconv.Itoa(eventIndexAfterSnapshot-1),
			Email:    "test" + strconv.Itoa(eventIndexAfterSnapshot-1) + "@example.com",
		},
	)
	return p
}

type Success struct{}

func (s *Success) AddUserIfNotExists(ctx context.Context, param mysql.AddUserParams) error {
	return nil
}

type Fail struct{}

func (s *Fail) AddUserIfNotExists(ctx context.Context, param mysql.AddUserParams) error {
	return errors.New("failed to add user /for tests")
}

func TestRestAPI_Receive(t *testing.T) {
	tests := []struct {
		name  string
		actor actor.Actor
		in    struct {
			stream  *stream.TypedStream[message.UserCreateMessenger]
			command command.CreateUser
		}
		want bool
	}{
		{
			name:  "success",
			actor: NewRestAPI(registration.NewCreateUser(registration.NewUserModelUpdate(&Success{}), NewProvider(3))),
			want:  true,
		},
		{
			name: "fail",
			actor: NewRestAPI(
				registration.NewCreateUser(
					registration.NewUserModelUpdate(&Fail{}), NewProvider(3).InitState("rest-api/user-test2@example.com", 4, 3))),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			system := actor.NewActorSystem()
			root, err := system.Root.SpawnNamed(actor.PropsFromProducer(func() actor.Actor {
				return tt.actor
			}), "rest-api")
			if err != nil {
				t.Errorf("failed to spawn actor: %v", err)
			}
			typedStream := stream.NewTypedStream[message.UserCreateMessenger](system)
			system.Root.Send(root, &command.CreateUser{
				UserName: "test2",
				Email:    "test2@example.com",
				Stream:   typedStream.PID(),
			})
			res := <-typedStream.C()
			if res.IsSuccess() != tt.want {
				t.Errorf("failed to create user: %v", res)
			}
		})
	}
}
