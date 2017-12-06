package pili

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Meeting struct {
	OwnerId string
	Client  *Client
}

func newMeeting(ownerId string, client *Client) *Meeting {
	return &Meeting{OwnerId: ownerId, Client: client}
}

func (c *Meeting) RoomStatus(room string) (status RoomStatusReturn, err error) {
	err = c.Client.Call(&status, "GET", fmt.Sprintf("%s%s/%s/rooms/%s", APIHTTPScheme, RTCAPIHost, RTCVersion, room))
	return
}

func (c *Meeting) RoomCreate(args CreateRoomArgs) (ret CreateRoomReturn, err error) {
	err = c.Client.CallWithJSON(&ret, "POST", fmt.Sprintf("%s%s/%s/rooms", APIHTTPScheme, RTCAPIHost, RTCVersion), args)
	return
}

func (c *Meeting) RoomDelete(room string) (ret DeleteRoomReturn, err error) {
	err = c.Client.Call(&ret, "DELETE", fmt.Sprintf("%s%s/%s/rooms/%s", APIHTTPScheme, RTCAPIHost, RTCVersion, room))
	return
}

func (c *Meeting) CreateToken(option RoomAccessPolicy) (string, error) {
	b, err := json.Marshal(option)
	if err != nil {
		return "", err
	}
	return signWithData(b, c.Client.mac), nil
}

func signWithData(b []byte, mac *MAC) (token string) {

	blen := base64.URLEncoding.EncodedLen(len(b))

	key := mac.AccessKey
	nkey := len(key)
	ret := make([]byte, nkey+30+blen)

	base64.URLEncoding.Encode(ret[nkey+30:], b)

	h := hmac.New(sha1.New, mac.SecretKey)
	h.Write(ret[nkey+30:])
	digest := h.Sum(nil)

	copy(ret, key)
	ret[nkey] = ':'
	base64.URLEncoding.Encode(ret[nkey+1:], digest)
	ret[nkey+29] = ':'

	return string(ret)
}

type CreateRoomReturn struct {
	Room string `json:"room_name"`
}

type DeleteRoomReturn struct {
}

type RoomStatusReturn struct {
	Room        string `json:"room_name"`
	OwnerUserID string `json:"owner_id"`
	UserMax     int    `json:"user_max"`
	// Status      int    `json:"room_status"`
}

type RoomAccessPolicy struct {
	Room     string `json:"room_name"`
	User     string `json:"user_id"`
	Perm     string `json:"perm"`
	Version  string `json:"version"`
	ExpireAt int64  `json:"expire_at"`
}

type CreateRoomArgs struct {
	User    string `json:"owner_id"`
	Room    string `json:"room_name"`
	UserMax int    `json:"user_max"`
}
