package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thomaspoignant/go-feature-flag/ffuser"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
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
				//appLauncher(bootstrapApplication)

			}
		},
	}
)

func main() {
	ffclient := config.InitApplication()

	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(cplanecmd)

	cplanecmd.PersistentFlags().StringVar(&bootstrapApplication, "application", "", "An application to be deployed alongside the controlplane")
	rootCmd.Execute()
	user := ffuser.NewUser(String(10))
	hasFlag, _ := ffclient.BoolVariation("test-flag", user, false)
	if hasFlag { // flag "test-flag" is true for the user
		fmt.Println("Flag true")
	} else { // flag "test-flag" is false for the user
		fmt.Println("flag false")
	}
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.LoadHTMLGlob("backend/web/templates/*")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"title": "Unity Repository",
		})
	})
	router.Use(static.Serve("/", static.LocalFile("./build", true)))
	router.NoRoute(func(c *gin.Context) { // fallback
		c.File("./build/index.html")
	})
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

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
		viper.SetConfigName(".unitycp")
		viper.SetDefault("GithubToken", "unset")
		viper.SetDefault("MarketplaceURL", "unset")

	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
		}
	}
}

/*func appLauncher(appname string) {
	//Lookup app from marketplace

	clustername := String(10)
	token := os.Getenv("GHTOKEN")
	//Deploy app via act
	prg := "/home/ubuntu/bin/act"
	arg1 := "-W"
	arg2 := ".github/workflows/test-action.yml"
	arg3 := "--input"
	arg4 := fmt.Sprintf(`METADATA={"metadataVersion":"unity-cs-0.1","deploymentName":"deployment","ghtoken":"%s", "services":[{"name":"unity-sps-prototype","source":"unity-sds/unity-sps-prototype","version":"xxx","branch":"main"}],"extensions":{"kubernetes":{"clustername":"%s","owner":"tom","projectname":"testproject","nodegroups":{"group1":{"instancetype":"m5.xlarge","nodecount":"1"}}}}}`, token, clustername)
	arg5 := "--env"
	arg6 := "WORKFLOWPATH=/home/ubuntu/unity-cs/.github/workflows"
	cmd := exec.Command(prg, arg1, arg2, arg3, arg4, arg5, arg6)
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
*/
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}
