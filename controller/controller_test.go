//
// Copyright (c) 2019 LG Electronics Inc.
// SPDX-License-Identifier: Apache-2.0
//

package controller_test

import (
	"log"
	"os"

	"github.com/jarcoal/httpmock"
	"github.com/lgsvl/data-marketplace-esearchFeeder/controller"
	"github.com/lgsvl/data-marketplace-esearchFeeder/resources"
	"github.com/lgsvl/data-marketplace-esearchFeeder/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("esfeeder-controller", func() {
	var (
		logger *log.Logger
		ctl    controller.Controller
		config resources.Config
		cc     resources.ChainRestConfig
		es     resources.ElasticSearchConfig
		idx    resources.IndicesConfig
		err    error
	)
	BeforeEach(func() {
		logger = log.New(os.Stdout, "esfeeder-controller-test: ", log.Lshortfile|log.LstdFlags)
		cc = resources.ChainRestConfig{Endpoint: "http://fake-endpoint:9999"}
		es = resources.ElasticSearchConfig{Endpoint: "http://fake-endpoint:7777"}
		idx = resources.IndicesConfig{}
		config = resources.Config{Authorization: "fake-autho",
			IdxConfig: idx,
			ESConfig:  es,
			CCConfig:  cc,
		}
		httpmock.RegisterResponder("GET", utils.FormatURL(es.Endpoint),
			httpmock.NewStringResponder(200, `[{}]`))

		ctl, err = controller.NewController(logger, config)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		httpmock.Reset()
	})
	Context(".GetData", func() {

		It("should fail when http execute fails", func() {

			request := resources.GetDataContractTypesRequest{
				Authorization: map[string]string{"authorization": "fake-auth"},
				QueryPath:     "fake-path",
				Container:     nil,
			}

			response := ctl.GetDataContractTypes(request)
			Expect(response.Error).To(Equal("error-getting-data-from-chaincode"))

		})

		It("should fail when http execute does not return 200 ok", func() {

			request := resources.GetDataContractTypesRequest{
				Authorization: map[string]string{"authorization": "fake-auth"},
				QueryPath:     "fake-path",
				Container:     nil,
			}

			httpmock.RegisterResponder("GET", utils.FormatURL(cc.Endpoint, request.QueryPath),
				httpmock.NewStringResponder(404, `[{}]`))

			response := ctl.GetDataContractTypes(request)
			Expect(response.Error).To(Equal("error-finding-data-response-code-404"))

		})

		It("should fail when chaincode returns a non json response", func() {

			request := resources.GetDataContractTypesRequest{
				Authorization: map[string]string{"authorization": "fake-auth"},
				QueryPath:     "fake-path",
				Container:     nil,
			}

			httpmock.RegisterResponder("GET", utils.FormatURL(cc.Endpoint, request.QueryPath),
				httpmock.NewStringResponder(200, `fake-non-json-response`))

			response := ctl.GetDataContractTypes(request)
			Expect(response.Error).To(Equal("failed-to-unmarshal-response"))

		})

		It("should fail when chaincode returns a valid json response and the response has an error message", func() {

			request := resources.GetDataContractTypesRequest{
				Authorization: map[string]string{"authorization": "fake-auth"},
				QueryPath:     "fake-path",
				Container:     nil,
			}

			jsonResp, err := utils.ReadFile(logger, "../assets/getDataContractTypesResponseWithError.json")
			Expect(err).NotTo(HaveOccurred())
			httpmock.RegisterResponder("GET", utils.FormatURL(cc.Endpoint, request.QueryPath), httpmock.NewStringResponder(200, string(jsonResp)))

			response := ctl.GetDataContractTypes(request)
			Expect(response.Error).To(Equal("error-message"))

		})

		It("should succeed when chaincode returns a valid json response", func() {

			request := resources.GetDataContractTypesRequest{
				Authorization: map[string]string{"authorization": "fake-auth"},
				QueryPath:     "fake-path",
				Container:     nil,
			}

			jsonResp, err := utils.ReadFile(logger, "../assets/getDataContractTypesResponse.json")
			Expect(err).NotTo(HaveOccurred())
			httpmock.RegisterResponder("GET", utils.FormatURL(cc.Endpoint, request.QueryPath), httpmock.NewStringResponder(200, string(jsonResp)))

			response := ctl.GetDataContractTypes(request)
			Expect(response.Error).To(Equal(""))
			Expect(response.ResponseMetadata[0].ResponseMetadata.Bookmark).To(Equal("g1AAAACieJw9zL0KwjAUQOHYyck3yQ0ILi6Ci-AjuJRLEtq0-SP3xlCf3nTQ9XD4vBBimAcjTsRYeLXbaJJ2JvR-m5kzXZXSaYqOk3QmQyXZLLE8Awb8pIiNQKeg_n2k-nhu-nVpO3v8sb6Dh12V_QY_WaC3B4PcmbJazh617VLkgprhXslFS7QsX6qLNpA"))
			Expect(response.Response[0].ID).To(Equal("5253d0d6-bd02-4802-a3c9-f477857f44ff"))
		})

	})

})
