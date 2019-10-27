package di

import (
	app "github.com/stepupgo/stepupgo2-1/pkg/application"
	"github.com/stepupgo/stepupgo2-1/pkg/infra/api"
)

var (
	Top app.ATop
)

func Init() {
	initTop()
}

func initTop() {
	t := api.NewITop()
	Top = app.NewTopApp(t)
}
