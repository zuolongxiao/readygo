package db

import (
	"context"
	"fmt"
	"time"
)

var ctx = context.Background()
var prefix = "captcha_"

// RedisCaptchaStore
type RedisCaptchaStore struct{}

func (RedisCaptchaStore) Set(id string, value string) error {
	exp := time.Minute * 10
	key := fmt.Sprintf("%s%s", prefix, id)
	if err := RDB.Set(ctx, key, value, exp).Err(); err != nil {
		return err
	}
	return nil
}
func (RedisCaptchaStore) Get(id string, clear bool) (val string) {
	key := fmt.Sprintf("%s%s", prefix, id)
	val, _ = RDB.Get(ctx, key).Result()
	// if err != nil {
	// 	fmt.Println("RDB.Get error:", err)
	// }
	if clear {
		RDB.Del(ctx, key).Err()
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
