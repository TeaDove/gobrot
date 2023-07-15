package cli

import (
	"runtime"

	"github.com/urfave/cli/v2"
)

func drawVideo(cCtx *cli.Context) error {
	runtime.GOMAXPROCS(cCtx.Int(maxprocsFlag.Name))

	return nil
}
