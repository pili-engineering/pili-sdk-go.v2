package pili

import (
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	testAccessKey    = ""
	testSecretKey    = ""
	testHub          = "PiliSDKTest"
	testStreamPrefix = "sdktest" + strconv.FormatInt(time.Now().UnixNano(), 10)
)

func init() {
	if v := os.Getenv("QINIU_ACCESS_KEY"); v != "" {
		testAccessKey = v
	}

	if v := os.Getenv("QINIU_SECRET_KEY"); v != "" {
		testSecretKey = v
	}

	if testAccessKey == "" || testSecretKey == "" {
		log.Fatal("need set access key and secret key")
	}

	if v := os.Getenv("PILI_API_HOST"); v != "" {
		APIHost = v
	}
}

func skipTest() bool {
	return testAccessKey == "" || testSecretKey == "" || testHub == ""
}

func TestURL(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}

	mac := &MAC{testAccessKey, []byte(testSecretKey)}
	expect := "rtmp://publish-rtmp.test.com/" + testHub + "/key?e="
	result := RTMPPublishURL("publish-rtmp.test.com", testHub, "key", mac, 10)
	require.True(t, strings.HasPrefix(result, expect))

	expect = "rtmp://live-rtmp.test.com/" + testHub + "/key"
	result = RTMPPlayURL("live-rtmp.test.com", testHub, "key")
	require.Equal(t, result, expect)

	expect = "http://live-hls.test.com/" + testHub + "/key.m3u8"
	result = HLSPlayURL("live-hls.test.com", testHub, "key")
	require.Equal(t, result, expect)

	expect = "http://live-hdl.test.com/" + testHub + "/key.flv"
	result = HDLPlayURL("live-hdl.test.com", testHub, "key")
	require.Equal(t, result, expect)

	expect = "http://live-snapshot.test.com/" + testHub + "/key.jpg"
	result = SnapshotPlayURL("live-snapshot.test.com", testHub, "key")
	require.Equal(t, result, expect)
}
