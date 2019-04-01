package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lgsvl/data-marketplace-esearchFeeder/resources"
	"github.com/lgsvl/data-marketplace-esearchFeeder/utils"
	"github.com/olivere/elastic"
)

type Controller interface {
	GetDataContractTypes(resources.GetDataContractTypesRequest) resources.GetDataContractTypesResponse
	FeedData(resources.FeedDataRequest) resources.FeedDataResponse
}

type controller struct {
	logger *log.Logger

	cCConfig        resources.ChainRestConfig
	eSConfig        resources.ElasticSearchConfig
	frequencyConfig resources.IndicesConfig

	esClient *elastic.Client
}

func NewController(l *log.Logger, config resources.Config) (Controller, error) {

	return &controller{
		logger:          l,
		cCConfig:        config.CCConfig,
		eSConfig:        config.ESConfig,
		frequencyConfig: config.IdxConfig,
	}, nil
}

func (c *controller) GetDataContractTypes(req resources.GetDataContractTypesRequest) resources.GetDataContractTypesResponse {
	c.logger.Printf("entering-controller-GetDataContractTypesData")
	defer c.logger.Printf("exiting-controller-GetDataContractTypesData")

	timestamp, err := utils.GetTimestamp(c.logger, resources.DefaultTimestampFile)
	if err != nil {
		c.logger.Printf("timestamp-not-read-using-default")
		timestamp = resources.DefaultTimestamp
	}
	if timestamp == "" {
		timestamp = resources.DefaultTimestamp
	}
	getDataURL := utils.FormatURL(c.cCConfig.Endpoint, req.QueryPath)
	getDataURL = fmt.Sprintf("%s?timestamp=%s", getDataURL, timestamp)

	now := time.Now()

	response, err := utils.HttpExecuteWithHeader(c.logger, &http.Client{}, "GET", getDataURL, req.Authorization)
	if err != nil {
		c.logger.Printf("utils-HttpExecute-failed-%s\n", err.Error())
		return resources.GetDataContractTypesResponse{Error: "error-getting-data-from-chaincode"}
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.logger.Printf("error-finding-data")
		return resources.GetDataContractTypesResponse{Error: fmt.Sprintf("error-finding-data-response-code-%#v", response.StatusCode)}
	}

	resp := resources.GetDataContractTypesResponse{}

	err = utils.UnmarshalResponse(context.Background(), c.logger, response, &resp)
	if err != nil {
		errorMsg := fmt.Sprintf("failed-to-unmarshal-response")
		c.logger.Printf(errorMsg)
		return resources.GetDataContractTypesResponse{Error: errorMsg}
	}

	c.logger.Printf("response-is-%#v\n", resp)

	err = utils.PeristsTimestamp(c.logger, now.Format("2006-01-02T15:04:05.000Z"), resources.DefaultTimestampFile)
	if err != nil {
		errorMsg := fmt.Sprintf("failed-to-persist-timestamp-%s", err.Error())
		c.logger.Printf(errorMsg)
	}

	return resp
}

func (c *controller) FeedData(req resources.FeedDataRequest) resources.FeedDataResponse {
	c.logger.Printf("entering-controller-feedData")
	defer c.logger.Printf("exiting-controller-feedData")

	if len(req.Bulk) == 0 {
		return resources.FeedDataResponse{Status: "success"}
	}

	if c.esClient == nil {
		esClient, err := elastic.NewClient(
			elastic.SetURL(c.eSConfig.Endpoint))

		if err != nil {
			c.logger.Printf("error-creating-esClient-%s\n", err.Error())
			return resources.FeedDataResponse{Error: err.Error()}
		}
		c.esClient = esClient
	}

	bulkReq := c.esClient.Bulk()

	for _, dataInstance := range req.Bulk {
		if dataInstance.Index == "" {
			continue
		}
		c.logger.Printf("adding-to-bulk-%#v\n", dataInstance)
		req := elastic.NewBulkIndexRequest().Index(strings.ToLower(dataInstance.Index)).Type(dataInstance.Type).Id(dataInstance.ID).Doc(dataInstance.Content)
		bulkReq = bulkReq.Add(req)
	}

	bulkResponse, err := bulkReq.Do(context.Background())
	if err != nil {
		errorMsg := fmt.Sprintf("error-feeding-data-%s", err.Error())
		c.logger.Println(errorMsg)
		return resources.FeedDataResponse{Error: errorMsg}
	}
	if bulkResponse != nil {
		c.logger.Printf("bulk-response-%#v\n", bulkResponse)
	}
	return resources.FeedDataResponse{Status: "success"}
}
