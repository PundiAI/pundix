package main

import (
	"os"

	"github.com/pundix/pundix/app"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	if err := svrcmd.Execute(NewRootCmd(), app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
