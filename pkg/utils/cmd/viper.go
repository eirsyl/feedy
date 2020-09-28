package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/camelcase"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Flag binding (attach flags to viper)

func uppercaseName(name string) string {
	s := camelcase.Split(name)
	snake := strings.Join(s, "_")

	return strings.ToUpper(snake)
}

func bindFlag(flags *pflag.FlagSet, name string) error {
	err := viper.BindPFlag(name, flags.Lookup(name))
	if err != nil {
		return err
	}

	err = viper.BindEnv(name, uppercaseName(name))
	if err != nil {
		return err
	}

	return nil
}

// Getters

func flagsetForGetter(cmd *cobra.Command, name string) *pflag.FlagSet {
	if cmd.PersistentFlags().Lookup(name) != nil {
		return cmd.PersistentFlags()
	}

	return cmd.Flags()
}

// GetString retrieves a string from the config.
func GetString(cmd *cobra.Command, name string) string {
	flags := flagsetForGetter(cmd, name)

	if err := bindFlag(flags, name); err != nil {
		panic(err)
	}

	return viper.GetString(name)
}

// GetBool retrieves a bool from the config.
func GetBool(cmd *cobra.Command, name string) bool {
	flags := flagsetForGetter(cmd, name)

	if err := bindFlag(flags, name); err != nil {
		panic(err)
	}

	return viper.GetBool(name)
}

// GetInt retrieves an int from the config.
func GetInt(cmd *cobra.Command, name string) int {
	flags := flagsetForGetter(cmd, name)

	if err := bindFlag(flags, name); err != nil {
		panic(err)
	}

	return viper.GetInt(name)
}

// GetDuration retrieves a duration from the config.
func GetDuration(cmd *cobra.Command, name string) time.Duration {
	flags := flagsetForGetter(cmd, name)

	if err := bindFlag(flags, name); err != nil {
		panic(err)
	}

	return viper.GetDuration(name)
}

// GetStringSlice retrieves a string slice from the config.
func GetStringSlice(cmd *cobra.Command, name string) []string {
	flags := flagsetForGetter(cmd, name)

	if err := bindFlag(flags, name); err != nil {
		panic(err)
	}

	return viper.GetStringSlice(name)
}

// GetStringToString retrieves a string map from the config.
func GetStringToString(cmd *cobra.Command, name string) map[string]string {
	flags := flagsetForGetter(cmd, name)

	if err := bindFlag(flags, name); err != nil {
		panic(err)
	}

	return viper.GetStringMapString(name)
}

// FlagChecker defines the flag check function.
type FlagChecker func() error

// CheckFlags takes multiple flag checkers and validates them.
func CheckFlags(checkers ...FlagChecker) {
	var fails []string

	for _, checker := range checkers {
		if err := checker(); err != nil {
			fails = append(fails, err.Error())
		}
	}

	if len(fails) > 0 {
		fmt.Println(strings.Join(fails, "\n"))
		os.Exit(1)
	}
}

// RequireString returns an error if the given setting is not a string.
func RequireString(cmd *cobra.Command, name string) FlagChecker {
	return func() error {
		v := GetString(cmd, name)
		if v == "" {
			return fmt.Errorf("flag %s can not be an empty string", name)
		}

		return nil
	}
}
