package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"github.com/checksum0/go-electrum/electrum"
	"github.com/gin-gonic/gin"
)

var cfg ServerConfig

func loadConfig(cfg *ServerConfig) {
	configFileName := "config.json"
	configFileName, _ = filepath.Abs(configFileName)
	log.Println("Loading config: ", configFileName)
	configFile, err := os.Open(configFileName)

	if err != nil {
		log.Fatalln("File error: ", err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&cfg); err != nil {
		log.Fatalln("Config error: ", err.Error())
	}

	log.Println("- config Payment address: ", cfg.Payment_address)
	if len(cfg.Payment_address) == 0 {
		log.Fatalln("Config variable Payment_address required.")
	}

	log.Println("- config Min_payment: ", cfg.Min_payment)
	if cfg.Min_payment < 0 {
		log.Fatalln("Invalid config variable Min_payment value.")
	}

	log.Println("- config Electrum_port: ", cfg.Electrum_port)
	if len(cfg.Electrum_host) == 0 {
		log.Fatalln("Config variable Electrum_host required.")
	}

	log.Println("- config Electrum_host: ", cfg.Electrum_host)
	if cfg.Electrum_port == 0 || cfg.Electrum_port < 0 {
		log.Fatalln("Config variable Electrum_port required.")
	}

	log.Println("- config Tls_enabled: ", cfg.Tls_enabled)
	if cfg.Tls_enabled == true {
		if len(cfg.Tls_cert) == 0 || len(cfg.Tls_key) == 0 {
			log.Fatalln("Config variables Tls_cert and Tls_key required.")
		}
	} else {
		log.Println("- config *Warning*: TLS/SSL not enabled. Set tls_enabled to true to enable TLS/SSL.")
	}
}

func setupServer() *gin.Engine {
	var err error = nil
	electrumServer := electrum.NewServer()
	portStr := strconv.Itoa(cfg.Electrum_port)

	if cfg.Tls_enabled {
		conf := &tls.Config{}
		if err = electrumServer.ConnectSSL(cfg.Electrum_host + ":" + portStr, conf); err != nil {
			log.Fatalln("Electrum connection error: ", err.Error())
		}
	} else {
		if err = electrumServer.ConnectTCP(cfg.Electrum_host + ":" + portStr); err != nil {
			log.Fatalln("Electrum connection error: ", err.Error())
		}
	}

	// Timed "server.ping" call to prevent disconnection.
	go func() {
		for {
			if err = electrumServer.Ping(); err != nil {
				// Log error, don't treat as fatal
				log.Println("Electrum keep alive error: ", err.Error())
			}
			time.Sleep(60 * time.Second)
		}
	}()

	router := gin.New()
	v1 := router.Group("/v1")
	{
		// Get payment info
		v1.GET("/payment_info", func(c *gin.Context) {
			getPaymentInfo(c, cfg.Payment_address, cfg.Min_payment)
		})

		// Upload media
		v1.POST("/media", func(c *gin.Context) {
			uploadMedia(c)
		})

		// Add to IPFS
		v1.POST("/pin", func(c *gin.Context) {
			publishMediaIPFS(electrumServer, c, cfg.Payment_address, cfg.Min_payment)
		})
	}

	return router
}

func main() {
	var router *gin.Engine

	// Setting log file
	file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Error creating log file: ", err.Error())
	}

	defer file.Close()
	log.SetOutput(file)

	// Setting gin to release mode; comment to use environment variable or systems default
	gin.SetMode(gin.ReleaseMode)

	loadConfig(&cfg)
	router = setupServer()
	// The port used by the server is the eletrumx port plus 10.
	portStr := strconv.Itoa(cfg.Electrum_port + 10)

	if cfg.Tls_enabled == true {
		router.RunTLS(":" + portStr, cfg.Tls_cert, cfg.Tls_key)
	} else {
		router.Run(":" + portStr)
	}
	log.Println("Listening on port: " + portStr)
}
