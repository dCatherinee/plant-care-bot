package domain

import "time"

type User struct {
	ID             int64
	TelegramUserID int64
	CreatedAt      time.Time
}

func NewUser(telegramUserID int64) (User, error) {
	if telegramUserID <= 0 {
		return User{}, ValidationError{
			Field:   "telegramUserID",
			Problem: "must be positive",
		}
	}

	return User{
		TelegramUserID: telegramUserID,
	}, nil
}
