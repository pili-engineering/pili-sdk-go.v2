# Pili Streaming Cloud Server-Side Library for Go

## Features

- URL
	- [x] RTMP推流地址: RTMPPublishURL(domain, hub, streamKey, mac, expireAfterSeconds)
	- [x] RTMP直播地址: RTMPPlayURL(domain, hub, streamKey)
	- [x] HLS直播地址: HLSPlayURL(domain, hub, streamKey)
	- [x] HDL直播地址: HDLPlayURL(domain, hub, streamKey)
	- [x] 直播封面地址: SnapshotPlayURL(domain, hub, streamKey)
- Hub
	- [x] 创建流: hub.Create(streamKey)
	- [x] 获得流: hub.Stream(streamKey)
	- [x] 列出流: hub.List(prefix, limit, marker)
	- [x] 列出正在直播的流: hub.ListLive(prefix, limit, marker)
	- [x] 批量查询直播信息: hub.BatchLiveStatus(streams)

- Stream
	- [x] 流信息: stream.Info()
	- [x] 禁用流: stream.DisableTill(till)
	- [x] 启用流: stream.Enable()
 	- [x] 查询直播状态: stream.LiveStatus()
	- [x] 保存直播回放: stream.Saveas(options)
	- [x] 保存直播截图: stream.Snapshot(options)
	- [x] 更改流的实时转码规格: stream.UpdateConverts(profiles)
	- [x] 查询直播历史: stream.HistoryActivity(start, end)

## Contents

- [Installation](#installation)
- [Usage](#usage)
    - [Configuration](#configuration)
	- [URL](#url)
		- [Generate RTMP publish URL](#generate-rtmp-publish-url)
		- [Generate RTMP play URL](#generate-rtmp-play-url)
		- [Generate HLS play URL](#generate-hls-play-url)
		- [Generate HDL play URL](#generate-hdl-play-url)
		- [Generate Snapshot play URL](#generate-snapshot-play-url)
	- [Hub](#hub)
		- [Instantiate a Pili Hub object](#instantiate-a-pili-hub-object)
		- [Create a new Stream](#create-a-new-stream)
		- [Get a Stream](#get-a-stream)
		- [List Streams](#list-streams)
		- [List live Streams](#list-live-streams)
		- [Batch query live Status](#batch-query-live-status)
	- [Stream](#stream)
		- [Get Stream info](#get-stream-info)
		- [Disable a Stream](#disable-a-stream)
		- [Enable a Stream](#enable-a-stream)
		- [Get Stream live status](#get-stream-live-status)
		- [Get Stream history activity](#get-stream-history-activity)
		- [Save Stream live playback](#save-stream-live-playback)
		- [Update Stream converts](#update-stream-converts)
		- [Save Stream snapshot](#save-stream-snapshot)

## Installation

before next step, install git.

```
// install latest version
$ go get github.com/pili-engineering/pili-sdk-go.v2/pili
```

## Usage

### Configuration

```go
package main

import (
	// ...
	"github.com/pili-engineering/pili-sdk-go.v2/pili"
)

var (
	AccessKey = "<QINIU ACCESS KEY>" // 替换成自己 Qiniu 账号的 AccessKey.
	SecretKey = "<QINIU SECRET KEY>" // 替换成自己 Qiniu 账号的 SecretKey.
	HubName   = "<PILI HUB NAME>"    // Hub 必须事先存在.
)

func main() {
	// ...
	mac := &pili.MAC{AccessKey, []byte(SecretKey)}
	client := pili.New(mac, nil)
	// ...
}
```

### URL

#### Generate RTMP publish URL

```go
url := pili.RTMPPublishURL("publish-rtmp.test.com", "PiliSDKTest", "streamkey", mac, 60)
fmt.Println(url)
/*
rtmp://publish-rtmp.test.com/PiliSDKTest/streamkey?e=1463023142&token=7O7hf7Ld1RrC_fpZdFvU8aCgOPuhw2K4eapYOdII:-5IVlpFNNGJHwv-2qKwVIakC0ME=
*/
```

#### Generate RTMP play URL

```go
url := pili.RTMPPlayURL("live-rtmp.test.com", "PiliSDKTest", "streamkey")
fmt.Println(url)
/*
rtmp://live-rtmp.test.com/PiliSDKTest/streamkey
*/
```

#### Generate HLS play URL

```go
url := pili.HLSPlayURL("live-hls.test.com", "PiliSDKTest", "streamkey")
fmt.Println(url)
/*
http://live-hls.test.com/PiliSDKTest/streamkey.m3u8
*/
```

#### Generate HDL play URL

```go
url := pili.HDLPlayURL("live-hdl.test.com", "PiliSDKTest", "streamkey")
fmt.Println(url)
/*
http://live-hdl.test.com/PiliSDKTest/streamkey.flv
*/
```

#### Generate Snapshot play URL

```go
url := pili.SnapshotPlayURL("live-snapshot.test.com", "PiliSDKTest", "streamkey")
fmt.Println(url)
/*
http://live-snapshot.test.com/PiliSDKTest/streamkey.jpg
*/
```

### Hub

#### Instantiate a Pili Hub object

```go
func main() {
	mac := &pili.MAC{AccessKey, []byte(SecretKey)}
	client := pili.New(mac, nil)
	hub := client.Hub("PiliSDKTest")
	// ...
}
```

#### Create a new Stream

```go
stream, err := hub.Create(key)
if err != nil {
	return
}
info, err := stream.Info()
if err != nil {
	return
}
fmt.Println(info)
/*
{hub:PiliSDKTest,key:streamkey,disabled:false}
*/
```

#### Get a Stream

```go
stream := hub.Stream(key)
info, err := stream.Info()
if err != nil {
	return
}
fmt.Println(info)
/*
{hub:PiliSDKTest,key:streamkey,disabled:false}
*/
```

#### List Streams

```go
keys, marker, err := hub.List(prefix, 10, "")
if err != nil {
	return
}
fmt.Printf("keys=%v marker=%v\n", keys, marker)
/*
keys=[streamkey] marker=
*/
```

#### List live Streams

```go
keys, marker, err := hub.ListLive(prefix, 10, "")
if err != nil {
	return
}
fmt.Printf("keys=%v marker=%v\n", keys, marker)
/*
keys=[streamkey] marker=
*/
```

#### Batch query live status
```go
items, err := hub.BatchLiveStatus(streams)
if err != nil {
	return
}
fmt.Println(items)
/*
[{streamKey1 {1487766696 172.21.2.14:63422 1042240 {46 24 0}}} {streamKey2 {1487768638 172.21.2.14:63793 1201352 {51 22 0}}}]
*/
```

### Stream

#### Get Stream info

```go
stream := hub.Stream(key)
info, err := stream.Info()
if err != nil {
	return
}
fmt.Println(info)
/*
{hub:PiliSDKTest,key:streamkey,disabled:false}
*/
```

#### Disable a Stream

```go
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

/*
before disable: {hub:PiliSDKTest,key:sdkexample1487765261435788231A,disabled:false}
after call disable: {hub:PiliSDKTest,key:sdkexample1487765261435788231A,disabled:true}
after time.Minute: {hub:PiliSDKTest,key:sdkexample1487765261435788231A,disabled:false}
*/

```

#### Enable a Stream

```go
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
/*
before enable: {hub:PiliSDKTest,key:streamkey,disabled:true}
after enable: {hub:PiliSDKTest,key:streamkey,disabled:false}
*/
```

#### Get Stream live status

```go
stream := hub.Stream(key)
status, err := stream.LiveStatus()
if err != nil {
	return
}
fmt.Printf("%+v\n", status)
/*
&{StartAt:1463382400 ClientIP:172.21.1.214:52897 BPS:128854 FPS:{Audio:38 Video:23 Data:0}}
*/
```

#### Get Stream history activity

```go
stream := hub.Stream(key)
records, err := stream.HistoryActivity(0, 0)
if err != nil {
	return
}
fmt.Println(records)
/*
[{1463382401 1463382441}]
*/
```

#### Save Stream live playback

```go
stream := hub.Stream(key)
opts := &pili.SaveasOptions{
	Format: "mp4",
}
fname, persistentID, err := stream.Saveas(opts)
if err != nil {
	return
}
fmt.Println(fname, persistentID)
/*
recordings/z1.PiliSDKTest.streamkey/1463156847_1463157463.mp4 z1.58ae57dc254f0e4d3f00000e
*/
```

### Update Stream converts
```go
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
/*
before UpdateConverts: {hub:PiliSDKTest,key:sdkexample1487834862156173949A,disabled:false,converts:[]}
after UpdateConverts: {hub:PiliSDKTest,key:sdkexample1487834862156173949A,disabled:false,converts:[480p 720p]}
*/
```

#### Save Stream snapshot
```go
stream := hub.Stream(key)
opts := &pili.SnapshotOptions{
	Format: "jpg",
}
fname, err := stream.Snapshot(opts)
if err != nil {
	return
}
fmt.Println(fname)
/*
streamkey-2741961933532577110.jpg
*/
```
