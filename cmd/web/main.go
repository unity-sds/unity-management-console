package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"os/exec"
)

var (
	cfgFile              string
	bootstrapApplication string
	rootCmd              = &cobra.Command{Use: "Unity", Short: "Unity Command Line Tool", Long: ""}
	cplanecmd            = &cobra.Command{
		Use:   "bootstrap",
		Short: "Execute control plane commands",
		Long:  `Control plane startup configuration commands`,
		Run: func(cmd *cobra.Command, args []string) {
			if bootstrapApplication != "" {
				appLauncher(bootstrapApplication)

			}
		},
	}
)

func main() {
	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	//   c.JSON(http.StatusOK, gin.H{
	//     "message": "pong",
	//   })
	// })
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(cplanecmd)

	cplanecmd.PersistentFlags().StringVar(&bootstrapApplication, "application", "", "An application to be deployed alongside the controlplane")
	rootCmd.Execute()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

func initConfig() {
	if cfgFile != "" { //
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func appLauncher(appname string) {
	//Lookup app from marketplace

	token := os.Getenv("GHTOKEN")
	//Deploy app via act
	prg := "/home/ubuntu/bin/act"
	arg4 := "-W"
	arg45 := ".github/workflows/test-action.yml"
	arg5 := "--input"
	arg55 := fmt.Sprintf(`METADATA={"metadataVersion":"unity-cs-0.1","deploymentName":"deployment","ghtoken":"%s", "services":[{"name":"unity-sps-prototype","source":"unity-sds/unity-sps-prototype","version":"xxx","branch":"main"}],"extensions":{"kubernetes":{"clustername":"testclustertomtues4","owner":"tom","projectname":"testproject","nodegroups":{"group1":{"instancetype":"m5.xlarge","nodecount":"1"}}}}}`, token)
	arg6 := "--env"
	arg65 := "WORKFLOWPATH=/home/ubuntu/unity-cs/.github/workflows"
	cmd := exec.Command(prg, arg4, arg45, arg5, arg55, arg6, arg65)
	cmd.Dir = "/home/ubuntu/unity-cs"
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	_ = cmd.Start()

	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	_ = cmd.Wait()
}
