package store

import (
	"context"
	"fmt"
	"time"

	"readygo/pkg/db"
	"readygo/pkg/settings"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
)

var CaptchaStore base64Captcha.Store

var ctx = context.Background()

// RedisCaptchaStore
type RedisCaptchaStore struct {
	rdb     *redis.Client
	prefix  string
	expires time.Duration
}

func (s RedisCaptchaStore) Set(id string, value string) error {
	key := fmt.Sprintf("%s%s", s.prefix, id)
	if err := s.rdb.Set(ctx, key, value, s.expires).Err(); err != nil {
		return err
	}
	return nil
}
func (s RedisCaptchaStore) Get(id string, clear bool) (val string) {
	key := fmt.Sprintf("%s%s", s.prefix, id)
	val, _ = s.rdb.Get(ctx, key).Result()
	// if err != nil {
	// 	fmt.Println("rdb.Get error:", err)
	// }
	if clear {
		s.rdb.Del(ctx, key).Err()
		// if err != nil {
		// 	fmt.Println("rdb.Del error:", err)
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
		CaptchaStore = RedisCaptchaStore{
			rdb:     db.RDB,
			prefix:  settings.Captcha.Prefix,
			expires: settings.Captcha.Expires,
		}
	case "Memory":
		CaptchaStore = base64Captcha.DefaultMemStore
	default:
		return fmt.Errorf("%s captcha store not implemented", settings.Captcha.Store)
	}
	return nil
}
