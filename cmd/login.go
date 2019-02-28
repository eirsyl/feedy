package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/eirsyl/feedy/pkg/config"
	"github.com/eirsyl/flexit/cmd"
)

func init() {
	cmd.BoolConfig(loginCmd, "force", "", false, "override existing authentication token if present")
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with pocket",
	Long: `
Authenticate with pocket and store authentication token.
	`,
	Args: cobra.NoArgs,
	PreRun: func(_ *cobra.Command, args []string) {
	},
	RunE: func(_ *cobra.Command, args []string) error {

		var force bool
		{
			force = viper.GetBool("force")
		}

		c, err := config.NewFileConfig()
		if err != nil {
			return err
		}

		return c.SaveToken("", force)
	},
}
