package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/pundix/pundix/app"
	"github.com/pundix/pundix/cmd"
)

func main() {
	if err := svrcmd.Execute(cmd.NewRootCmd(), app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
