package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	ccResources "github.com/lgsvl/data-marketplace-chaincode/resources"
	"github.com/lgsvl/data-marketplace-esearchFeeder/resources"
)

func LoadConfig(logger *log.Logger) (resources.Config, error) {
	logger.Println("entering-utils-loadconfig")
	defer logger.Println("exiting-utils-loadconfig")
	config := resources.Config{}
	indicesConfig := resources.IndicesConfig{}
	esConfig := resources.ElasticSearchConfig{}
	ccConfig := resources.ChainRestConfig{}

	jsonIndexTimes := os.Getenv("INDEX_TIMES")
	indexTimes := make(map[string]int)
	err := json.Unmarshal([]byte(jsonIndexTimes), &indexTimes)

	if err != nil {
		errorMsg := fmt.Sprintf("error-unmarshalling-jsonIndexTimes-%s-error-%s", jsonIndexTimes, err.Error())
		logger.Println(errorMsg)
		return config, fmt.Errorf(errorMsg)
	}

	indicesConfig.IndicesToTimeMap = indexTimes
	jsonIndexPaths := os.Getenv("INDEX_PATHS")

	indexPaths := make(map[string]string)
	err = json.Unmarshal([]byte(jsonIndexPaths), &indexPaths)
	if err != nil {
		errorMsg := fmt.Sprintf("error-unmarshalling-jsonIndexPaths-%s-error-%s", jsonIndexPaths, err.Error())
		logger.Println(errorMsg)
		return config, fmt.Errorf(errorMsg)
	}
	indicesConfig.IndicesToCCQueryMap = indexPaths
	esConfig.Endpoint = os.Getenv("ELASTIC_SEARCH_URL")
	if esConfig.Endpoint == "" {
		logger.Println("ELASTIC_SEARCH_URL-env-variable-not-set-using-default")
		esConfig.Endpoint = resources.DefaultElasticSearchEndpoint
	}
	ccConfig.Endpoint = os.Getenv("CHAINCODE_REST_URL")

	if ccConfig.Endpoint == "" {
		logger.Println("CHAINCODE_REST_URL-env-variable-not-set-using-default")
		ccConfig.Endpoint = resources.DefaultChaincodeRESTEndpoint

	}

	auth := os.Getenv("AUTHORIZATION")
	config.Authorization = auth
	config.ESConfig = esConfig
	config.CCConfig = ccConfig
	config.IdxConfig = indicesConfig

	return config, nil
}

func PrepareDataContractTypesForFeeding(logger *log.Logger, data []ccResources.DataContractType) ([]resources.Data, error) {
	logger.Println("entering-utils-PrepareDataContractTypesForFeeding")
	defer logger.Println("exiting-utils-PrepareDataContractTypesForFeeding")
	if data == nil {
		logger.Println("data-is-empty")
		return nil, nil
	}
	res := make([]resources.Data, len(data))
	for _, d := range data {
		logger.Printf("appending-data-id-%s\n", d.ID)
		temp := resources.Data{
			ID:      d.ID,
			Index:   resources.DataContractTypesIndex,
			Type:    d.DocType,
			Content: d,
		}
		res = append(res, temp)
	}

	return res, nil
}

func PeristsTimestamp(logger *log.Logger, timestamp string, file string) error {
	logger.Println("entering-utils-PeristsTimestamp")
	defer logger.Println("exiting-utils-PeristsTimestamp")

	err := ioutil.WriteFile(file, []byte(timestamp), 0644)
	if err != nil {
		log.Printf("error-writing-to-timestamp-file-%s", err.Error())
		return err
	}
	return nil

}

func GetTimestamp(logger *log.Logger, timestampFile string) (string, error) {
	logger.Println("entering-utils-GetTimestamp")
	defer logger.Println("exiting-utils-GetTimestamp")

	input, err := ioutil.ReadFile(timestampFile)
	if err != nil {
		logger.Printf("error-reading-timestamp-file-%s", err.Error())
		return "", err
	}

	return string(input), nil

}
