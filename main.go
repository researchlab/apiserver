package main

import (
	"apiserver/config"
	"apiserver/model"
	"apiserver/pkg/version"
	"apiserver/router"
	"apiserver/router/middleware"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/lexkong/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path.")
	v   = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	pflag.Parse()
	if *v {
		_version := version.Get()
		marshalled, err := json.MarshalIndent(&_version, "", " ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(marshalled))
		return
	}
	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	//Set gin mode
	gin.SetMode(viper.GetString("runmode"))
	// init db
	model.DB.Init()
	defer model.DB.Close()
	//Create the Gin engine
	g := gin.New()

	//Routes
	g = router.Load(
		//Cores
		g,
		//Middlewares.
		middleware.RequestId(),
		middleware.Logging(),
	)
	//Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()
	//Start to listening the incoming requests.
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")
	if len(cert) != 0 && len(key) != 0 {
		go func() {
			log.Infof("Start to listening the incoming requests on https address: %s", viper.GetString("tls.addr"))
			log.Info(http.ListenAndServeTLS(viper.GetString("tls.addr"), cert, key, g).Error())
		}()
	}
	log.Infof("Start to listening the incoming requests on http address %s", viper.GetString("addr"))
	log.Infof(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		//Ping the server by sending a GET request to '/health'.
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		//Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}
