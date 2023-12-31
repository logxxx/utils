package netutil

import "testing"

func TestTryGetFileSize(t *testing.T) {

	reqURL := "https://tpc.googlesyndication.com/simgad/8234315780176092875?sqp=4sqPyQQrQikqJwhfEAEdAAC0QiABKAEwCTgDQPCTCUgAUAFYAWBfcAJ4AcUBLbKdPg&rs=AOga4qmpD7I1cH0fqLhFrdparrgjaxiaag"
	resp := TryGetFileSize(reqURL)
	t.Logf("resp:%v", resp)

	err := DownloadToFile(reqURL, "test.jpg")
	if err != nil {
		panic(err)
	}

}
