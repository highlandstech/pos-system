package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// parseConfig parses the configuration either from the same folder, or
// from an explicit path.
func parseConfig(customConfig *string) (*Config, error) {
	if customConfig == nil || len(*customConfig) == 0 {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigFile(*customConfig)
	}

	var config Config

	// Try to read the configuration file.
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Default config.
	viper.SetDefault("server.port", 8088)

	// Parse the configuration into the config object.
	err := viper.Unmarshal(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	customConfig := flag.String("config", "", "The path to a custom config file")
	flag.Parse()

	db, err := ConnectDB()

	if err != nil {
		panic(err)
	}

	config, err := parseConfig(customConfig)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Instanciate all of the route handlers here.
	productHandlers := &ProductHandlers{DB: db}
	orderHandlers := &OrderHandlers{DB: db}
	userHandlers := &UserHandlers{
		DB:     db,
		Config: config,
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// Apply the sessions middleware.
	store := cookie.NewStore([]byte(config.Secret))
	router.Use(sessions.Sessions("pos-sessions", store))

	// Set the static/public path.
	router.Use(static.Serve("/", static.LocalFile("./public", false)))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
			"admin": false,
		})
	})

	router.GET("/admin", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user").(*UserStruct)

		// Prevent anyone who is not logged in to view this page.
		if user == nil || !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{})
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Admin",
			"admin": true,
		})
	})

	router.POST("/login", userHandlers.Login)

	router.POST("/api/order", orderHandlers.CreateOrder)
	router.GET("/api/order/:orderId", orderHandlers.PrintOrder)

	router.POST("/api/product", productHandlers.CreateProduct)
	router.GET("/api/products", productHandlers.ListProducts)

	router.Run(":" + strconv.Itoa(config.Server.Port))
}
