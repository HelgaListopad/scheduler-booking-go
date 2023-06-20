package main

import "scheduler-booking/data"

type ConfigServer struct {
	URL  string
	Port string
	Cors []string
}

type AppConfig struct {
	Server ConfigServer
	DB     data.DBConfig
}
