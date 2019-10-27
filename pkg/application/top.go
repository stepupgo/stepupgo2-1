package application

import (
	"net/http"

	"github.com/stepupgo/stepupgo2-1/pkg/infra/api"
)

type top struct {
	api api.ITop
}

type ATop interface {
	GetAvailable() (*http.Response, error)
}

func NewTopApp(t api.ITop) ATop {
	return &top{t}
}

func (t *top) GetAvailable() (*http.Response, error) {
	resp, err := t.api.GetAvailableLotteries()
	return resp, err
}
