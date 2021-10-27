package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envPrefix         = "IMAGE_SERVER"
	defaultConfigFile = "conf.yml"
	// Config paths
	config     = "config.path"
	listenAddr = "http.listen"
	dbDriver   = "db.driver"
	dbDSN      = "db.url"
	s3Bucket   = "s3.bucketName"
)

var rootCmd = &cobra.Command{
	Use: "image-server",

	Short: "Image server is a server for images of dynamic size and quality",
	Long:  `Image server is a server for images of dynamic size and quality`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Root flags
	rootCmd.PersistentFlags().StringP(config, "c", "", "Path to config file")
	viper.BindPFlag(config, rootCmd.PersistentFlags().Lookup(config))
	rootCmd.PersistentFlags().StringP(s3Bucket, "n", "", "Name of S3 bucket")
	viper.BindPFlag(s3Bucket, rootCmd.PersistentFlags().Lookup(s3Bucket))

	// Server flags
	serverCmd.PersistentFlags().StringP(listenAddr, "l", ":8080", "Address for server to listen to in the format ip:port")
	viper.BindPFlag(listenAddr, serverCmd.PersistentFlags().Lookup(listenAddr))
	serverCmd.PersistentFlags().String(dbDriver, "sqlite", "Driver to use for connecting to db")
	viper.BindPFlag(dbDriver, serverCmd.PersistentFlags().Lookup(dbDriver))
	serverCmd.PersistentFlags().String(dbDSN, "./test.db", "DSN to connect to database with")
	viper.BindPFlag(dbDSN, serverCmd.PersistentFlags().Lookup(dbDSN))

	rootCmd.AddCommand(serverCmd)
}

func initConfig() {

	if viper.GetString(config) != "" {
		viper.SetConfigFile(viper.GetString(config))
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("conf")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
