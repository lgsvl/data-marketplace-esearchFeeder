package resources

import (
	"github.com/lgsvl/data-marketplace-chaincode/resources"
)

const (
	DefaultElasticSearchEndpoint        string = "elasticsearch:9200"
	DefaultChaincodeRESTEndpoint        string = "chaincode-rest-service:9090"
	DefaultGetDataContractTypesPath     string = "/queries/getDataContractTypesAfterTimeStamp"
	DefaultGetDataContractTypesPageSize int    = 100
	DataContractTypesIndex              string = "dataContractType"
	DefaultTimestampFile                string = "/var/run/esfeeder/datacontracttypes"
	DefaultTimestamp                    string = "0000-01-01T00:00:00.000Z"
)

type Config struct {
	Authorization string
	IdxConfig     IndicesConfig
	ESConfig      ElasticSearchConfig
	CCConfig      ChainRestConfig
}
type IndicesConfig struct {
	IndicesToTimeMap    map[string]int
	IndicesToCCQueryMap map[string]string
}

type ElasticSearchConfig struct {
	Endpoint string
}

type ChainRestConfig struct {
	Endpoint string
}

type GetDataContractTypesRequest struct {
	Authorization map[string]string
	QueryPath     string
	BookMark      string
	PageSize      int
	Container     []interface{}
}

type FeedDataRequest struct {
	Bulk []Data `json: "bulk"`
}
type FeedDataResponse struct {
	Status string
	Error  string
}

type Data struct {
	Index   string      `json:"index"`
	Type    string      `json:"type"`
	ID      string      `json:"id"`
	Content interface{} `json:"content"`
}

type GetDataContractTypesResponse struct {
	Response         []resources.DataContractType `json:"response"`
	ResponseMetadata []ResponseMetadata           `json:"ResponseMetadata`
	Error            string                       `json:"error"`
}

type ResponseMetadata struct {
	ResponseMetadata Metadata `json:"ResponseMetadata"`
}

type Metadata struct {
	RecordsCount string `json:"RecordsCount"`
	Bookmark     string `json:"Bookmark"`
}
