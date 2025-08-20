package redis

import (
	"context"
	"encoding/json"
	"time"
	"todo/internal/models"
)

func StoreSession(session_uuid string, user_id int, ip string, ua string) error {
	ctx := context.Background()

	var session models.Session

	session.UID = user_id
	session.IAT = time.Now().Unix()
	session.EXP = time.Now().Add(time.Hour).Unix()
	session.IP = ip
	session.UA = ua

	val, err := json.Marshal(session)

	if err != nil {
		return err
	}

	err = Client.Set(ctx, "session:"+session_uuid, val, time.Hour).Err()

	return err
}

func GetSession(session_uuid string) (models.Session, error) {
	ctx := context.Background()

	var session models.Session

	res, err := Client.Get(ctx, "session:"+session_uuid).Result()

	if err != nil {
		return session, err
	}

	err = json.Unmarshal([]byte(res), &session)

	return session, err
}

func GetDeleteSession(session_uuid string) (string, error) {
	ctx := context.Background()

	res, err := Client.GetDel(ctx, "session:"+session_uuid).Result()

	return res, err
}

func RenewSession(session_uuid string) error {
	ctx := context.Background()

	err := Client.Expire(ctx, "session:"+session_uuid, time.Hour).Err()

	return err
}
