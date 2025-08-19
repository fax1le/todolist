package redis

import (
	"context"
	"time"
)

func StoreSession(session_uuid string, user_id int) error {
	ctx := context.Background()

	err := Client.Set(ctx, "session:"+session_uuid, user_id, time.Hour).Err()

	return err
}

func GetDeleteSession(session_uuid string) (string, error) {
	ctx := context.Background()

	res, err := Client.GetDel(ctx, "session:"+session_uuid).Result()

	return res, err
}

func SessionExists(session_uuid string) bool {
	ctx := context.Background()

	err := Client.Get(ctx, "session:"+session_uuid).Err()

	return err == nil
}

func GetUID(session_uuid string) (string, error) {
	ctx := context.Background()

	res, err := Client.Get(ctx, "session:"+session_uuid).Result()

	return res, err
}
