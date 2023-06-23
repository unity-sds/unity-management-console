package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"github.com/unity-sds/unity-control-plane/backend/internal/aws"
	"github.com/unity-sds/unity-control-plane/backend/internal/database"
	"github.com/unity-sds/unity-control-plane/backend/internal/processes"
	"github.com/unity-sds/unity-control-plane/backend/internal/web"
	"os"
	"path/filepath"
)

var (
	conf config.AppConfig

	cfgFile   string
	bootstrap bool
	rootCmd   = &cobra.Command{Use: "Unity", Short: "Unity Command Line Tool", Long: ""}
	cplanecmd = &cobra.Command{
		Use:   "webapp",
		Short: "Execute management console commands",
		Long:  `Management console startup configuration commands`,
		Run: func(cmd *cobra.Command, args []string) {
			store, err := database.NewGormDatastore()
			if bootstrap == true {
				//appLauncher(bootstrapApplication)
				storeDefaultSSMParameters(conf, store)
				r := processes.ActRunnerImpl{}
				err := processes.UpdateCoreConfig(nil, store, conf, &r)
				if err != nil {
					log.WithError(err).Error("Problem updating ssm config")
				}
				//provisionS3(conf)
				//installGateway(conf)

			}
			router := web.DefineRoutes(conf)

			err = router.Run()
			if err != nil {
				log.Errorf("Failed to launch API. %v", err)
				return
			}
		},
	}
)

func storeDefaultSSMParameters(appConfig config.AppConfig, store database.Datastore) {

	err := store.StoreSSMParams(appConfig.DefaultSSMParameters, "bootstrap")
	if err != nil {
		log.WithError(err).Error("Problem storing parameters in database.")
	}
}

func provisionS3(appConfig config.AppConfig) {
	aws.CreateBucket(appConfig)
}

func installGateway(appConfig config.AppConfig) {
	runner := processes.ActRunnerImpl{}
	meta := ""
	err := processes.InstallMarketplaceApplication(nil, meta, appConfig, "", runner)
	if err != nil {
		return
	}
}

func main() {
	log.Info("Launching Unity Management Console")

	cobra.OnInitialize(initConfig)

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

	if cfgFile != "" { //
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(configdir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("unity")
		viper.SetDefault("GithubToken", "unset")
		viper.SetDefault("MarketplaceURL", "unset")
		viper.SetDefault("WorkflowBasePath", "unset")
		viper.SetDefault("Workdir", "unset")

	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			file, createErr := os.Create(filepath.Join(configdir, "unity.yaml"))
			if createErr != nil {
				log.Fatalf("Failed to create config file: %v", createErr)
			}
			defer file.Close()
			log.Infof("Created config file: %s", viper.ConfigFileUsed())
		} else {
			// Config file was found but another error was produced
			log.Errorf("Failed to read config file: %v", err)
		}
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Errorf("Unable to decode into struct, %v", err)
	}

}
