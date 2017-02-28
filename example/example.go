package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pili-engineering/pili-sdk-go.v2/pili"
)

var (
	AccessKey = ""            // Qiniu 账号的 AccessKey.
	SecretKey = ""            // Qiniu 账号的 SecretKey.
	HubName   = "PiliSDKTest" // Hub 必须事先存在.
)

func init() {
	if v := os.Getenv("QINIU_ACCESS_KEY"); v != "" {
		AccessKey = v
	}

	if v := os.Getenv("QINIU_SECRET_KEY"); v != "" {
		SecretKey = v
	}

	if AccessKey == "" || SecretKey == "" {
		log.Fatal("need set access key and secret key")
	}

	if v := os.Getenv("PILI_API_HOST"); v != "" {
		pili.APIHost = v
	}
}

func createStream(hub *pili.Hub, key string) {
	stream, err := hub.Create(key)
	if err != nil {
		return
	}
	info, err := stream.Info()
	if err != nil {
		return
	}
	fmt.Println(info)
}

func getStream(hub *pili.Hub, key string) {
	stream := hub.Stream(key)
	info, err := stream.Info()
	if err != nil {
		return
	}
	fmt.Println(info)
}

func listStreams(hub *pili.Hub, prefix string) {
	keys, marker, err := hub.List(prefix, 10, "")
	if err != nil {
		return
	}
	fmt.Printf("keys=%v marker=%v\n", keys, marker)
}

func listLiveStreams(hub *pili.Hub, prefix string) {
	keys, marker, err := hub.ListLive(prefix, 10, "")
	if err != nil {
		return
	}
	fmt.Printf("keys=%v marker=%v\n", keys, marker)
}

func batchQueryLiveStreams(hub *pili.Hub, streams []string) {
	items, err := hub.BatchLiveStatus(streams)
	if err != nil {
		return
	}
	fmt.Println(items)
}

func updateStreamConverts(hub *pili.Hub, key string) {
	stream := hub.Stream(key)
	info, err := stream.Info()
	if err != nil {
		return
	}
	fmt.Println("before UpdateConverts:", info)
	err = stream.UpdateConverts([]string{"480p", "720p"})
	if err != nil {
		return
	}
	info, err = stream.Info()
	if err != nil {
		return
	}
	fmt.Println("after UpdateConverts:", info)
}

func disableStream(hub *pili.Hub, key string) {
	stream := hub.Stream(key)
	info, err := stream.Info()
	if err != nil {
		return
	}
	fmt.Println("before disable:", info)

	err = stream.DisableTill(time.Now().Add(time.Minute).Unix())
	if err != nil {
		return
	}

	info, err = stream.Info()
	if err != nil {
		return
	}
	fmt.Println("after call disable:", info)

	time.Sleep(time.Minute)

	info, err = stream.Info()
	if err != nil {
		return
	}
	fmt.Println("after time.Minute:", info)
}

func enableStream(hub *pili.Hub, key string) {
	stream := hub.Stream(key)
	info, err := stream.Info()
	if err != nil {
		return
	}
	fmt.Println("before enable:", info)

	err = stream.Enable()
	if err != nil {
		return
	}

	info, err = stream.Info()
	if err != nil {
		return
	}
	fmt.Println("after enable:", info)
}

func liveStatus(hub *pili.Hub, key string) {
	stream := hub.Stream(key)
	status, err := stream.LiveStatus()
	if err != nil {
		return
	}
	fmt.Printf("%+v\n", status)
}

func historyActivity(hub *pili.Hub, key string) {
	stream := hub.Stream(key)
	records, err := stream.HistoryActivity(0, 0)
	if err != nil {
		return
	}
	fmt.Println(records)
}

func savePlayback(hub *pili.Hub, key string) {
	stream := hub.Stream(key)
	opts := &pili.SaveasOptions{
		Format: "mp4",
	}
	fname, persistentID, err := stream.Saveas(opts)
	if err != nil {
		return
	}
	fmt.Println(fname, persistentID)
}

func saveSnapshot(hub *pili.Hub, key string) {
	stream := hub.Stream(key)
	opts := &pili.SnapshotOptions{
		Format: "jpg",
	}
	fname, err := stream.Snapshot(opts)
	if err != nil {
		return
	}
	fmt.Println(fname)
}

func main() {
	streamKeyPrefix := "sdkexample" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// 初始化 client & hub.
	mac := &pili.MAC{AccessKey: AccessKey, SecretKey: []byte(SecretKey)}
	client := pili.New(mac, nil)
	hub := client.Hub(HubName)

	keyA := streamKeyPrefix + "A"
	fmt.Println("获得不存在的流A:")
	streamA := hub.Stream(keyA)
	_, err := streamA.Info()
	fmt.Println(err, "IsNotExists", pili.IsNotExists(err))

	fmt.Println("创建流:")
	createStream(hub, keyA)

	fmt.Println("获得流:")
	getStream(hub, keyA)

	fmt.Println("创建重复流:")
	_, err = hub.Create(keyA)
	fmt.Println(err, "IsExists", pili.IsExists(err))

	keyB := streamKeyPrefix + "B"
	fmt.Println("创建另一路流:")
	createStream(hub, keyB)

	fmt.Println("列出流:")
	listStreams(hub, "carter")

	fmt.Println("列出正在直播的流:")
	listLiveStreams(hub, "carter")

	fmt.Println("批量查询直播信息:")
	batchQueryLiveStreams(hub, []string{keyA, keyB})

	fmt.Println("更改流的实时转码规格:")
	updateStreamConverts(hub, keyA)

	fmt.Println("禁用流:")
	disableStream(hub, keyA)

	fmt.Println("启用流:")
	enableStream(hub, keyA)

	fmt.Println("查询直播状态:")
	liveStatus(hub, keyA)

	fmt.Println("查询推流历史:")
	historyActivity(hub, keyA)

	fmt.Println("保存直播数据:")
	savePlayback(hub, keyA)

	fmt.Println("保存直播截图:")
	saveSnapshot(hub, keyA)

	fmt.Println("RTMP 推流地址:")
	url := pili.RTMPPublishURL("publish-rtmp.test.com", HubName, keyA, mac, 3600)
	fmt.Println(url)

	fmt.Println("RTMP 直播放址:")
	url = pili.RTMPPlayURL("live-rtmp.test.com", HubName, keyA)
	fmt.Println(url)

	fmt.Println("HLS 直播地址:")
	url = pili.HLSPlayURL("live-hls.test.com", HubName, keyA)
	fmt.Println(url)

	fmt.Println("HDL 直播地址:")
	url = pili.HDLPlayURL("live-hdl.test.com", HubName, keyA)
	fmt.Println(url)

	fmt.Println("截图直播地址:")
	url = pili.SnapshotPlayURL("live-snapshot.test.com", HubName, keyA)
	fmt.Println(url)
}
