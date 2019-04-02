package main

import (
	"os"
	"time"

	"github.com/lgsvl/data-marketplace-esearchFeeder/controller"
	"github.com/lgsvl/data-marketplace-esearchFeeder/resources"
	"github.com/lgsvl/data-marketplace-esearchFeeder/utils"
)

func main() {

	//Get the list of indexes to pupulate, their times and urls
	//For each index create a go routine with a specific ticker
	// Whenever the ticker happens, call get data and feed data methods
	logger := utils.CreateLogger("esearFeeder")
	config, err := utils.LoadConfig(logger)

	ctl, err := controller.NewController(logger, config)
	if err != nil {
		logger.Fatal("failed-to-create-controller")
	}

	if _, err := os.Stat(resources.DefaultTimestampFile); os.IsNotExist(err) {
		_, err := os.Create(resources.DefaultTimestampFile)
		if err != nil {
			logger.Fatal("failed-to-create-bookmark-file")
		}
	}

	for index, frequency := range config.IdxConfig.IndicesToTimeMap {
		queryURL, ok := config.IdxConfig.IndicesToCCQueryMap[index]
		if !ok {
			logger.Fatal("index-url-entry-not-found")
		}
		ticker := time.NewTicker(time.Duration(frequency) * time.Second)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					req := resources.GetDataContractTypesRequest{
						Authorization: map[string]string{"authorization": config.Authorization},
						QueryPath:     queryURL,
						Container:     []interface{}{},
					}
					resp := ctl.GetDataContractTypes(req)
					if resp.Error != "" {
						logger.Printf(resp.Error)
					} else {
						cleanData, err := utils.PrepareDataContractTypesForFeeding(logger, resp.Response)
						if err != nil {
							logger.Printf("error-cleaning-data")
						}
						feedDataRequest := resources.FeedDataRequest{Bulk: cleanData}
						feedDataResp := ctl.FeedData(feedDataRequest)
						if feedDataResp.Error != "" {
							logger.Printf("error-feeding-data")
						}
					}
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
	}

	select {}
}
