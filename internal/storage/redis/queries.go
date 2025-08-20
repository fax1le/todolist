package redis

import (
	"todo/internal/models"
	"context"
	"encoding/json"
	"time"
)

func StoreSession(session_uuid string, user_id int, ip string, ua string) error {
	ctx := context.Background()

	var session models.Session

	session.UID = user_id
	session.IAT = time.Now().Unix()
	session.IP = ip
	session.UA = ua

	val, err := json.Marshal(session)	

	if err != nil {
		return err
	}

	err = Client.Set(ctx, "session:"+session_uuid, val, time.Hour).Err()

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

func GetUID(session_uuid string) (int, error) {
	ctx := context.Background()

	res, err := Client.Get(ctx, "session:"+session_uuid).Result()

	if err != nil {
		return -1, err
	}

	var val models.Session

	err = json.Unmarshal([]byte(res), &val)
	
	return val.UID, err
}
