package store

import (
	"context"
	"fmt"
	"time"

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
var prefix = "captcha_"

// RedisCaptchaStore
type RedisCaptchaStore struct{}

func (RedisCaptchaStore) Set(id string, value string) error {
	exp := time.Minute * 10
	key := fmt.Sprintf("%s%s", prefix, id)
	if err := db.RDB.Set(ctx, key, value, exp).Err(); err != nil {
		return err
	}
	return nil
}
func (RedisCaptchaStore) Get(id string, clear bool) (val string) {
	key := fmt.Sprintf("%s%s", prefix, id)
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

func Setup() {
	switch settings.Captcha.Store {
	case "Redis":
		CaptchaStore = RedisCaptchaStore{}
	default:
		CaptchaStore = base64Captcha.DefaultMemStore
	}
}
