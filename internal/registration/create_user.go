package registration

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"github.com/oklog/ulid/v2"
	"github.com/ytake/protoactor-go-cqrs-example/internal/command"
	"github.com/ytake/protoactor-go-cqrs-example/internal/message"
	"github.com/ytake/protoactor-go-cqrs-example/pkg/event"
	"google.golang.org/protobuf/proto"
)

// User is an actor to create user
type User struct {
	persistence.Mixin
	stream *actor.PID
	state  *event.UserCreated
	rmu    *actor.PID
}

func NewUser(stream *actor.PID, rmu *actor.PID) actor.Actor {
	return &User{
		stream: stream,
		rmu:    rmu,
	}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor
// このアクターは永続化を行うため、persistence.Mixinを埋め込む
// 永続化を行い、Read Model更新アクターを起動する
func (u *User) Receive(context actor.Context) {
	defer context.Poison(context.Self())
	switch msg := context.Message().(type) {
	case *persistence.RequestSnapshot:
		u.PersistSnapshot(u.state)
	case *persistence.ReplayComplete:
		// リプレイが完了したら内部状態を変更する
		context.Logger().Info(
			fmt.Sprintf("replay completed, internal state changed to '%v'", u.state))
	case *command.CreateUser:
		// スナップショットにユーザーが存在する場合は生成済みエラーを返す
		if u.IsStateExists(msg.Email) {
			context.Send(u.stream, &message.UserCreateError{Message: "user already exists"})
			return
		}
		// ユーザーが存在しない場合はユーザーを生成する
		ev := &event.UserCreated{
			UserName: msg.UserName,
			Email:    msg.Email,
			UserID:   ulid.Make().String(),
		}
		u.createUser(context, ev)
		// ユーザー生成成功を通知する
		context.Send(u.stream, &message.UserCreateResponse{UserID: ev.UserID, Success: true})
	case *event.UserCreated:
		if msg.String() != "" {
			// event がリプレイされた場合は状態を更新する
			u.state = msg
			u.sendToReadModelUpdater(context, msg)
		}
	}
}

// IsStateExists is a method to check state exists or not
func (u *User) IsStateExists(email string) bool {
	if u.state == nil {
		return false
	}
	return u.state.Email == email
}

// createUser is a method to create user
func (u *User) createUser(context actor.Context, msg *event.UserCreated) {
	// バージョンを利用して楽観的ロックを実現することができます
	// この例では、バージョンをインクリメントして永続化し、Read Model更新アクターに送信します
	// 実際に実装する場合は、HTTPリクエストにversionを含めて送信するなどで、永続化されたデータとの競合を防ぐことができます
	msg.Version++
	u.persist(msg)
	u.sendToReadModelUpdater(context, msg)
}

func (u *User) sendToReadModelUpdater(context actor.Context, ev *event.UserCreated) {
	context.RequestWithCustomSender(u.rmu, ev, context.Self())
}

func (u *User) persist(msg proto.Message) {
	if !u.Recovering() {
		u.PersistReceive(msg)
	}
	switch ev := msg.(type) {
	case *event.UserCreated:
		u.state = ev
	}
}
