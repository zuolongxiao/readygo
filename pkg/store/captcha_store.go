package store

import (
	"context"
	"fmt"

	"readygo/pkg/db"
	"readygo/pkg/settings"

	"github.com/mojocn/base64Captcha"
)

type CaptchaStorer interface {
	Get(string, bool) string
	Set(string, string) error
	Verify(string, string, bool) bool
}

var CaptchaStore CaptchaStorer

var ctx = context.Background()

// RedisCaptchaStore
type RedisCaptchaStore struct{}

func (RedisCaptchaStore) Set(id string, value string) error {
	key := fmt.Sprintf("%s%s", settings.Captcha.Prefix, id)
	if err := db.RDB.Set(ctx, key, value, settings.Captcha.Expires).Err(); err != nil {
		return err
	}
	return nil
}
func (RedisCaptchaStore) Get(id string, clear bool) (val string) {
	key := fmt.Sprintf("%s%s", settings.Captcha.Prefix, id)
	val, _ = db.RDB.Get(ctx, key).Result()
	// if err != nil {
	// 	fmt.Println("RDB.Get error:", err)
	// }
	if clear {
		db.RDB.Del(ctx, key).Err()
		// if err != nil {
		// 	fmt.Println("RDB.Del error:", err)
		// }
	}
	return
}
func (s RedisCaptchaStore) Verify(id, answer string, clear bool) bool {
	code := s.Get(id, clear)
	return code == answer
}

func Setup() error {
	switch settings.Captcha.Store {
	case "Redis":
		CaptchaStore = RedisCaptchaStore{}
	case "Memory":
		CaptchaStore = base64Captcha.DefaultMemStore
	default:
		return fmt.Errorf("%s captcha store not implemented", settings.Captcha.Store)
	}
	return nil
}
