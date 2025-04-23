package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/processes"
	"github.com/unity-sds/unity-management-console/backend/internal/version"
	"github.com/unity-sds/unity-management-console/backend/internal/web"
	"math/rand"
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
			filename := filepath.Join(appConfig.Workdir, "workspace", "provider.tf")
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				log.Infof("File %s doesn't exist", filename)
				initialised = false
			} else if err == nil {
				log.Infof("File %s exists", filename)
				initialised = true
			} else {
				// There was some other error when trying to check the file
				log.Errorf("Error occurred while checking file: %s", err)
			}
			if bootstrap == true || !initialised {
				log.Info("Bootstrap flag set or uninitialised workdir, bootstrapping")
				err := processes.BootstrapEnv(&appConfig)
				if err != nil {
					log.Errorf("Errors were encountered during bootstrap process. Please check the errorlog for more details.")
				}
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
	log.Infof("Launching Unity Management Console v%s", version.GetVersion())

	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(cplanecmd)

	// Add version flag to root command
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Show version information")
	
	// Add custom pre-run function to check for version flag
	oldPreRun := rootCmd.PersistentPreRun
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			log.Infof("Unity Management Console v%s", version.GetVersion())
			os.Exit(0)
		}
		
		if oldPreRun != nil {
			oldPreRun(cmd, args)
		}
	}

	cplanecmd.PersistentFlags().BoolVar(&bootstrap, "bootstrap", false, "Provision an S3 bucket, Bootstrap an API Gateway for access to the management console")
	err := rootCmd.Execute()
	if err != nil {
		log.Errorf("Failed to parse CLI. %v", err)
		return
	}

	config.InitApplication()
}

func generateRandomString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
func initConfig() {
	uniqueString, err := generateRandomString(6)

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
		viper.SetDefault("MarketplaceOwner", "unity-sds")
		viper.SetDefault("MarketplaceRepo", "unity-marketplace")
		viper.SetDefault("Workdir", filepath.Join(path, "workdir"))
		viper.SetDefault("AWSRegion", "us-west-2")
		viper.SetDefault("MarketplaceBaseUrl", "https://raw.githubusercontent.com/")
		viper.SetDefault("BasePath", "")
		viper.SetDefault("ConsoleHost", "")
		viper.SetDefault("InstallPrefix", uniqueString)

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
	
	// Set the version in the app config
	appConfig.Version = version.GetVersion()

}
