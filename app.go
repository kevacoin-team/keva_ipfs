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
	log.Printf("Loading config: %v", configFileName)
	configFile, err := os.Open(configFileName)

	if err != nil {
		log.Fatal("File error: ", err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&cfg); err != nil {
		log.Fatal("Config error: ", err.Error())
	}

	log.Printf("- config Payment address: %v", cfg.Payment_address)
	if len(cfg.Payment_address) == 0 {
		log.Fatalln("Config variable Payment_address required.")
	}

	log.Printf("- config Min_payment: %v", cfg.Min_payment)
	if cfg.Min_payment < 0 {
		log.Fatalln("Invalid config variable Min_payment value.")
	}

	log.Printf("- config Electrum_port: %v", cfg.Electrum_port)
	if len(cfg.Electrum_host) == 0 {
		log.Fatalln("Config variable Electrum_host required.")
	}

	log.Printf("- config Electrum_host: %v", cfg.Electrum_host)
	if cfg.Electrum_port == 0 || cfg.Electrum_port < 0 {
		log.Fatalln("Config variable Electrum_port required.")
	}

	log.Printf("- config Tls_enabled: %v", cfg.Tls_enabled)
	if cfg.Tls_enabled == true {
		if len(cfg.Tls_cert) == 0 || len(cfg.Tls_key) == 0 {
			log.Fatalln("Config variables Tls_cert and Tls_key required.")
		}
	}
}

func setupServer() *gin.Engine {
	var err error = nil
	electrumServer := electrum.NewServer()
	portStr := strconv.Itoa(cfg.Electrum_port)

	if cfg.Tls_enabled {
		conf := &tls.Config{}
		if err = electrumServer.ConnectSSL(cfg.Electrum_host + ":" + portStr, conf); err != nil {
			log.Fatal("Electrum connection error: ", err)
		}
	} else {
		if err = electrumServer.ConnectTCP(cfg.Electrum_host + ":" + portStr); err != nil {
			log.Fatal("Electrum connection error: ", err)
		}
	}

	// Timed "server.ping" call to prevent disconnection.
	go func() {
		for {
			if err = electrumServer.Ping(); err != nil {
				log.Fatal("Electrum keep alive error: ", err)
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
	// Setting gin to release mode; comment to use environment variable or systems default
	gin.SetMode(gin.ReleaseMode)

	loadConfig(&cfg)
	router = setupServer()
	// The port used by the server is the eletrumx port plus 10.
	portStr := strconv.Itoa(cfg.Electrum_port + 10)

	if cfg.Tls_enabled == true {
		log.Println("Using TLS/SSL")
		router.RunTLS(":" + portStr, cfg.Tls_cert, cfg.Tls_key)
	} else {
		log.Println("**Warning: TLS/SSL not enabled. Set tls_enabled to true to enable TLS/SSL.")
		router.Run(":" + portStr)
	}
	log.Println("Listening on port: " + portStr)
}
