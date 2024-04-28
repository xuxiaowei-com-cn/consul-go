package command

import (
	"github.com/urfave/cli/v2"
)

func KvCommand() *cli.Command {
	return &cli.Command{
		Name:  "kv",
		Usage: "Key / Value",
		Flags: append(Common()),
		Subcommands: []*cli.Command{
			KvExportCommand(),
			KvImportCommand(),
		},
	}
}
