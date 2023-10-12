package tlscertreloader

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReload(t *testing.T) {
	key := "./.vscode/test.key"
	cert := "./.vscode/test.crt"
	r := MustNewCertReloader(cert, key, WithPeriod(1*time.Second))
	v, err := r.GetCertificate(nil)
	assert.NoError(t, err)
	assert.NotNil(t, v)
}
