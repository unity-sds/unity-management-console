package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/processes"
	"github.com/unity-sds/unity-management-console/backend/internal/web"
	"os"
	"path/filepath"
)

var (
	appConfig config.AppConfig

	cfgFile     string
	bootstrap   bool
	initialised bool
	rootCmd     = &cobra.Command{Use: "Unity", Short: "Unity Command Line Tool", Long: ""}
	cplanecmd   = &cobra.Command{
		Use:   "webapp",
		Short: "Execute management console commands",
		Long:  `Management console startup configuration commands`,
		Run: func(cmd *cobra.Command, args []string) {
			if bootstrap == true || !initialised {
				log.Info("Bootstrap flag set or uninitialised workdir, bootstrapping")
				processes.BootstrapEnv(&appConfig)
			}
			router := web.DefineRoutes(appConfig)

			err := router.Run()
			if err != nil {
				log.Errorf("Failed to launch API. %v", err)
				return
			}
		},
	}
)

func main() {
	log.Info("Launching Unity Management Console")

	cobra.OnInitialize(initConfig)

	filename := filepath.Join(appConfig.Workdir, "workspace", "provider")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		initialised = false
	} else if err == nil {
		initialised = true
	} else {
		// There was some other error when trying to check the file
		log.Errorf("Error occurred while checking file: %s", err)
	}

	rootCmd.AddCommand(cplanecmd)

	cplanecmd.PersistentFlags().BoolVar(&bootstrap, "bootstrap", false, "Provision an S3 bucket, Bootstrap an API Gateway for access to the management console")
	err := rootCmd.Execute()
	if err != nil {
		log.Errorf("Failed to parse CLI. %v", err)
		return
	}

	config.InitApplication()

}

func initConfig() {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("Error fetching home directory: %v", err)
		return
	}

	configdir := filepath.Join(dir, ".unity")

	if _, err := os.Stat(configdir); os.IsNotExist(err) {
		errDir := os.MkdirAll(configdir, 0755)
		if errDir != nil {
			log.Errorf("Error creating directory: %v", errDir)
			return
		}
	}
	path, err := os.Getwd()
	if cfgFile != "" { //
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(configdir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("unity")
		viper.SetDefault("GithubToken", "unset")
		viper.SetDefault("MarketplaceOwner", "unity-sds")
		viper.SetDefault("MarketplaceRepo", "unity-marketplace")
		viper.SetDefault("Workdir", filepath.Join(path, "workdir"))
		viper.SetDefault("AWSRegion", "us-west-2")
		viper.SetDefault("MarketplaceBaseUrl", "https://raw.githubusercontent.com/")
		viper.SetDefault("BasePath", "")
		viper.SetDefault("ConsoleHost", "")

	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			file, createErr := os.Create(filepath.Join(configdir, "unity.yaml"))
			if createErr != nil {
				log.WithError(createErr).Panicf("Failed to create config file")
			}
			defer file.Close()
			log.Infof("Created config file: %s", viper.ConfigFileUsed())
		} else {
			// Config file was found but another error was produced
			log.WithError(err).Panicf("Failed to read config file")
		}
	}
	if err := viper.Unmarshal(&appConfig); err != nil {
		log.WithError(err).Panicf("Unable to decode into struct")
	}

}
