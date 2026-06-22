package cmd

import (
	"github.com/imohamedsheta/xapp/app/database/seeders"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/spf13/cobra"
)

var SeedCommand = &cobra.Command{
	Use:   "seed",
	Short: "Test helper command while creating feature",
	Run:   handleSeedCmd(),
}

func handleSeedCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := seeders.SeedDatabase(); err != nil {
			utils.PrintErr("Seed Database Failed %v", err)
		} else {
			utils.PrintSuccess("Seed Database Successfully")
		}
	}
}
