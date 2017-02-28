package pili

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStreamInfo(t *testing.T) {

	info := StreamInfo{"a", "b", 0, nil}
	require.False(t, info.Disabled())
	require.Equal(t, info.String(), "{hub:a,key:b,disabled:false,converts:[]}")

	info = StreamInfo{"a", "b", -1, nil}
	require.True(t, info.Disabled())
	require.Equal(t, info.String(), "{hub:a,key:b,disabled:true,converts:[]}")

	info = StreamInfo{"a", "b", 12345, nil}
	require.False(t, info.Disabled())
	require.Equal(t, info.String(), "{hub:a,key:b,disabled:false,converts:[]}")

	info = StreamInfo{"a", "b", time.Now().Unix() + 1000, nil}
	require.True(t, info.Disabled())
	require.Equal(t, info.String(), "{hub:a,key:b,disabled:true,converts:[]}")
}

func TestStream(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	mac := &MAC{testAccessKey, []byte(testSecretKey)}
	client := New(mac, nil)
	hub := client.Hub(testHub)

	// Create.
	key := testStreamPrefix + "TestStream"
	_, err := hub.Create(key)
	require.NoError(t, err)
	stream := hub.Stream(key)
	require.NoError(t, err)
	info, err := stream.Info()
	require.NoError(t, err)
	require.True(t, checkStream(info, testHub, key, false))

	// Disable.
	stream = hub.Stream(key)
	err = stream.Disable()
	require.NoError(t, err)
	info, err = stream.Info()
	require.NoError(t, err)
	require.True(t, checkStream(info, testHub, key, true))

	// DisableTill
	stream = hub.Stream(key)
	err = stream.DisableTill(time.Now().Unix() + 1)
	require.NoError(t, err)
	info, err = stream.Info()
	require.True(t, checkStream(info, testHub, key, true))
	time.Sleep(time.Second)
	info, err = stream.Info()
	require.True(t, checkStream(info, testHub, key, false))

	// Converts
	stream = hub.Stream(key)
	err = stream.UpdateConverts([]string{"480p", "720p"})
	require.NoError(t, err)
	info, err = stream.Info()
	require.Equal(t, []string{"480p", "720p"}, info.Converts)
	err = stream.UpdateConverts(nil)
	require.NoError(t, err)
	info, err = stream.Info()
	require.Equal(t, []string{}, info.Converts)

	// Enable.
	err = stream.Enable()
	require.NoError(t, err)
	info, err = stream.Info()
	require.NoError(t, err)
	require.True(t, checkStream(info, testHub, key, false))

	// LiveStatus, no live.
	_, err = stream.LiveStatus()
	require.Equal(t, err, ErrNoLive)

	// Save, not found.
	_, err = stream.Save(0, 0)
	require.Error(t, err)
	require.Equal(t, err, ErrNoData)

	// HistoryActivity, empty.
	records, err := stream.HistoryActivity(0, 0)
	require.NoError(t, err)
	require.True(t, len(records) == 0)
}

// 这个测试case需要保持推流(test1)
func TestStreamSave(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	mac := &MAC{testAccessKey, []byte(testSecretKey)}
	client := New(mac, nil)
	hub := client.Hub(testHub)

	stream := hub.Stream("test1")
	fname, _, err := stream.Saveas(nil)
	require.NoError(t, err)
	require.True(t, len(fname) > 0)

	opts := &SaveasOptions{
		Format: "mp4",
		Fname:  "test1.mp4",
	}
	fname, pid, err := stream.Saveas(opts)
	require.NoError(t, err)
	require.Equal(t, "test1.mp4", fname)
	require.True(t, pid != "")

	opts.Pipeline = "notexist"
	_, _, err = stream.Saveas(opts)
	require.Equal(t, "no such pipeline", err.Error())
}

// 这个测试case需要保持推流(test1)
func TestStreamSnapshot(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	mac := &MAC{testAccessKey, []byte(testSecretKey)}
	client := New(mac, nil)
	hub := client.Hub(testHub)

	stream := hub.Stream("test1")
	fname, err := stream.Snapshot(nil)
	require.NoError(t, err)
	require.True(t, len(fname) > 0)

	opts := &SnapshotOptions{
		Fname: "test1.jpg",
	}
	fname, err = stream.Snapshot(opts)
	require.NoError(t, err)
	require.Equal(t, "test1.jpg", fname)
}
