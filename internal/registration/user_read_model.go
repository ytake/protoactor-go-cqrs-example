package registration

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/ytake/protoactor-go-cqrs-example/internal/database/mysql"
	"github.com/ytake/protoactor-go-cqrs-example/pkg/event"
)

// UserModelUpdate is actor to update read model
type UserModelUpdate struct {
	query mysql.RegistrationUserExecutor
}

// NewUserModelUpdate is constructor for UserModelUpdate
func NewUserModelUpdate(query mysql.RegistrationUserExecutor) actor.Actor {
	return &UserModelUpdate{
		query: query,
	}
}

// Receive is sent messages to be processed from the mailbox associated with the instance of the actor
func (u *UserModelUpdate) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *event.UserCreated:
		// イベントソーシングのイベントを読み込んで、Read Modelを更新する
		// ここではRead Modelにユーザーが存在しない場合はユーザーを生成する
		err := u.query.AddUserIfNotExists(context.Background(), mysql.AddUserParams{
			Email: msg.Email,
			Name:  msg.UserName,
			ID:    msg.UserID,
		})
		if err != nil {
			// エラーが発生した場合はログを出力する
			ctx.Logger().Error(err.Error())
			return
		}
	}
}
