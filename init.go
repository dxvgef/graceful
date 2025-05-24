package graceful

import (
	"context"
)

func init() {
	mainContext, mainCancel = context.WithCancel(context.Background())
}
