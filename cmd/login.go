package cmd

import (
	"context"
	"fmt"

	"github.com/eirsyl/feedy/pkg/client"
	"github.com/eirsyl/feedy/pkg/pocket"
	"github.com/jedib0t/go-pretty/text"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/eirsyl/feedy/pkg/config"
)

var loginCmd = &cobra.Command{
	Use:   "login [consumer-key]",
	Short: "Authenticate with pocket",
	Long: `
Authenticate with pocket and store authentication token.
	`,
	Args: cobra.ExactArgs(1),
	PreRun: func(_ *cobra.Command, args []string) {
	},
	RunE: func(_ *cobra.Command, args []string) error {

		var consumerKey string
		{
			consumerKey = args[0]
		}

		c, err := config.GetConfigProvider()
		if err != nil {
			return errors.Wrap(err, "could not create config backend")
		}
		defer c.Close() // nolint: errcheck

		cc, err := client.New()
		if err != nil {
			return errors.Wrap(err, "could not create http client")
		}

		p, err := pocket.New(cc)
		if err != nil {
			return errors.Wrap(err, "could not create pocket client")
		}

		token, err := p.Login(context.TODO(), consumerKey)
		if err != nil {
			return errors.Wrap(err, "login failed")
		}

		if err = c.SaveUser(config.User{
			ConsumerKey: consumerKey,
			Token:       token,
		}); err != nil {
			return errors.Wrap(err, "could not store access token")
		}

		fmt.Println(text.FgGreen.Sprint("Login successful"))

		return nil
	},
}
