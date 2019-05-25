package main

import (
	"flag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

func main() {

	log.Printf("Start")

	//viper.SetConfigName("config")
	//viper.AddConfigPath(".")
	//viper.MergeInConfig()

	// using standard library "flag" package
	flag.Int("flagname", 1234, "help message for flagname")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	i := viper.GetInt("flagname") // retrieve value from viper

	log.Printf("Value: %d", i)
}
