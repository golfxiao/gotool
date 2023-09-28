package upgrade

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := Init("./doc/us.conf")
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestDownload(t *testing.T) {
	usdown := &USDownload{
		RootPath:  "share",
		ReleaseId: 5,
		Path:      []string{"/6/60000/api-ms-win-core-console-l1-1-0.dll", "/6/60000/api-ms-win-core-handle-l1-1-0.dll", "/resources/html/login/script/109_14103ef2f44c08c77550.js", "/resources/pages/userList.html", "/VideoEngineCore.dll", "/VideoMixerEngine.dll"},
	}
	patch, err := usdown.Download()
	assert.Nil(t, err)
	assert.NotNil(t, patch)
	t.Logf("patch: %+v", patch)

}
