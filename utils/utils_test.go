//
// Copyright (c) 2019 LG Electronics Inc.
// SPDX-License-Identifier: Apache-2.0
//

package utils_test

import (
	"log"
	"os"

	"github.com/lgsvl/data-marketplace-esearchFeeder/resources"
	"github.com/lgsvl/data-marketplace-esearchFeeder/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cli-controller", func() {
	var logger *log.Logger
	BeforeEach(func() {
		logger = log.New(os.Stdout, "esfeeder-utils-test: ", log.Lshortfile|log.LstdFlags)

	})
	AfterEach(func() {
		os.Unsetenv("INDEX_TIMES")
		os.Unsetenv("INDEX_PATHS")
		Expect(os.Getenv("INDEX_TIMES")).To(Equal(""))
		Expect(os.Getenv("INDEX_PATHS")).To(Equal(""))
	})

	Context(".LoadConfig", func() {
		It("should fail when index times is not a json object", func() {

			_, err := utils.LoadConfig(logger)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("error-unmarshalling-jsonIndexTimes--error-unexpected end of JSON input"))
		})

		It("should fail when index paths is not a json object", func() {
			os.Setenv("INDEX_TIMES", "{\"DataContractType\":10}")
			_, err := utils.LoadConfig(logger)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("error-unmarshalling-jsonIndexPaths--error-unexpected end of JSON input"))
		})

		It("should succeed when config is valid", func() {
			os.Setenv("INDEX_TIMES", "{\"DataContractType\":10}")
			os.Setenv("INDEX_PATHS", "{\"DataContractType\":\"fake-url\"}")

			config, err := utils.LoadConfig(logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(config.CCConfig.Endpoint).To(Equal(resources.DefaultChaincodeRESTEndpoint))
			Expect(config.ESConfig.Endpoint).To(Equal(resources.DefaultElasticSearchEndpoint))
			Expect(config.IdxConfig.IndicesToCCQueryMap).ToNot(BeNil())
			Expect(config.IdxConfig.IndicesToTimeMap).ToNot(BeNil())

		})
	})
})
