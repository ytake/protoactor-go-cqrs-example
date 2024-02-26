package command

import "github.com/asynkron/protoactor-go/actor"

// CreateUser はユーザーを生成を指示するコマンド
type CreateUser struct {
	UserName string
	Email    string
	Stream   *actor.PID
}
