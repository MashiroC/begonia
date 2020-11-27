package dispatch

import (
	"github.com/MashiroC/begonia/dispatch/frame"
)

// dispatch_default.go something

type recvMsg struct {
	connID string
	f      frame.Frame
}