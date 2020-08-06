package pili

import (
	"fmt"
	"pili-sdk-go.v2/pili"
	"testing"
	"time"
)

const (
	AccessKey = ""
	SecretKey = ""
)

func testCreateRoom(t *testing.T) {
	clien := pili.New(&pili.MAC{AccessKey, []byte(SecretKey)}, nil)
	meeting := clien.Meeting("123")
	//创建房间
	args := pili.CreateRoomArgs{
		User:    "123",
		Room:    "roomName",
		UserMax: 10,
	}
	ret, err := meeting.RoomCreate(args)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v\n", ret)
}

func testGetRoom() {
	clien := pili.New(&pili.MAC{AccessKey, []byte(SecretKey)}, nil)
	meeting := clien.Meeting("123")
	roomName := "testRoom"
	ret, err := meeting.RoomStatus(roomName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v\n", ret)
}

func testDeleteRoom() {
	clien := pili.New(&pili.MAC{AccessKey, []byte(SecretKey)}, nil)
	meeting := clien.Meeting("123")
	roomName := "testRoomName"
	ret, err := meeting.RoomDelete(roomName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v\n", ret)
}

func testRoomToken() {
	client := pili.New(&pili.MAC{AccessKey, []byte(SecretKey)}, nil)
	meeting := client.Meeting("123")
	args := pili.RoomAccessPolicy{
		Room:     "roomName",
		User:     "123",
		Perm:     "admin",
		Version:  "2.0",
		ExpireAt: time.Now().Add(time.Hour * 24).UnixNano(),
	}
	token, err := meeting.CreateToken(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)
}
