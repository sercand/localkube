package localkubectl

import (
	"io"
	"os"
)

var (
	Out io.Writer = os.Stdout
)
