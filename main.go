package main

import (
	"bkio-zipper/fs"
	"github.com/spf13/viper"
	"log"
	"sync"
)

func main() {
	log.Println("Begin database backup zip program")
	ch := make(chan string)

	go func() {
		wg := sync.WaitGroup{}

		for _, c := range fs.C {
			wg.Add(1)
			go fs.ZipBackups(c.Backup, c.Output, ch, &wg)
		}

		wg.Wait()

		close(ch)
	}()

	for val := range ch {
		log.Println(val)
	}
}

func init() {
	viper.AddConfigPath("./setup")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	err = viper.UnmarshalKey("dbs", &fs.C)
	if err != nil {
		log.Fatal(err)
	}
}
