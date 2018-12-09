package realdebrid

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_AddMagnetLinkSimple(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "https://api.real-debrid.com/rest/1.0/torrents/addMagnet", req.URL.String())
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "magnet-url", req.FormValue("magnet"))

		return &http.Response{
			StatusCode: http.StatusCreated,
			Body: ioutil.NopCloser(bytes.NewBufferString(
				`{ "id": "MNREAKNMGAG7C", "uri": "https://api.real-debrid.com/rest/1.0/torrents/info/MNREAKNMGAG7C" }`,
			)),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		}
	})

	urlInfo, err := client.AddMagnetLinkSimple("magnet-url")
	assert.NoError(t, err)
	assert.Equal(t, TorrentUrlInfo{ID: "MNREAKNMGAG7C", URI: "https://api.real-debrid.com/rest/1.0/torrents/info/MNREAKNMGAG7C"}, urlInfo)
}

func TestClient_GetTorrent(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "https://api.real-debrid.com/rest/1.0/torrents/info/MNREAKNMGAG7C", req.URL.String())
		assert.Equal(t, "GET", req.Method)

		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(bytes.NewBufferString(`
{ "id": "XCBYL4ZIYPU42", "filename": "test", "original_filename": "test", "hash": "05d9df877f471dc4418fe1160cd8ff51b5258f55", "bytes": 957874366,
"original_bytes": 957874454, "host": "real-debrid.com", "split": 2000, "progress": 100, "status": "downloaded", "added": "2018-12-08T21:57:33.000Z", "files": [
{ "id": 1, "path": "/testfile.dat", "bytes": 30, "selected": 0 } ], "links": [ "https://real-debrid.com/d/XN7FEFLLQWJJS" ], "ended": "2018-11-28T04:23:19.000Z" } `,
			)),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		}
	})

	info, err := client.GetTorrent("MNREAKNMGAG7C")
	assert.NoError(t, err)
	assert.Equal(t, TorrentInfo{
		Added:            time.Date(2018, 12, 8, 21, 57, 33, 0, time.UTC),
		Ended:            time.Date(2018, 11, 28, 04, 23, 19, 0, time.UTC),
		ID:               "XCBYL4ZIYPU42",
		Filename:         "test",
		OriginalFilename: "test",
		Hash:             "05d9df877f471dc4418fe1160cd8ff51b5258f55",
		Bytes:            957874366,
		OriginalBytes:    957874454,
		Host:             "real-debrid.com",
		Split:            2000,
		Progress:         100,
		Status:           StatusDownloaded,
		Files: []File{
			{
				ID:       1,
				Path:     "/testfile.dat",
				Bytes:    30,
				Selected: 0,
			},
		},
		Links: []string{"https://real-debrid.com/d/XN7FEFLLQWJJS"},
	}, info)
}

func TestClient_SelectFilesFromTorrent(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "https://api.real-debrid.com/rest/1.0/torrents/selectFiles/XCBYL4ZIYPU42", req.URL.String())
		assert.Equal(t, "POST", req.Method)

		return &http.Response{
			StatusCode: http.StatusCreated,
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		}
	})

	err := client.SelectFilesFromTorrent("XCBYL4ZIYPU42", []int{1, 2, 3})
	assert.NoError(t, err)
}