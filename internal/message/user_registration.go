package message

type UserCreateMessenger interface {
	// IsSuccess is success or not
	IsSuccess() bool
}

type UserCreateResponse struct {
	UserID string `json:"user_id"`
	// 失敗成功
	Success bool `json:"-"`
}

func (u *UserCreateResponse) IsSuccess() bool {
	return u.Success
}

type UserCreateError struct {
	// エラーメッセージ
	Message string
}

func (u *UserCreateError) IsSuccess() bool {
	return false
}
