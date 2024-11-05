# Web Template
This tool is built to help developers create a web application quickly with essential initializations. It includes the following features:
- Go
- Gin
- zap
- viper
- Redis
- MySQL

Congif.yaml file contains the configuration for the application.

The project structure is as follows:
```
web_template
├── README.md
├── main.go
├── config.yaml
├── settings
│   └── settings.go
├── controllers
├── models
├── routes
│   └── routes.go
├── utils
│   └── logger
│       └── logger.go
|   └── mysql
│       └── mysql.go
|   └── redis
│       └── redis.go


How to use:
1. Clone the repository
2. Run `go mod download` to download the required dependencies
3. Update the config.yaml file with the required configuration
4. Run `go run main.go` to start the application
5. Open the browser and go to `http://localhost:XXXX` to access the application
Note: Replace `XXXX` with the port number specified in the config.yaml file.