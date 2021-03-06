package config_test

import (
	. "github.com/Wikia/konfigurator/config"

	"strings"

	"github.com/Wikia/konfigurator/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _ = Describe("Config", func() {
	Context("with sample variable configuration", func() {
		conf := map[string]string{
			"foo":       "vault(key1)",
			"bar":       "consul(key2)",
			"simple":    "simple(key3)",
			"config":    "simple(key6)->config",
			"inline":    "simple(key7)->config",
			"reference": "reference(node.name)->config",
			"simple2":   "simple(key4)->config",
			"simple3":   "simple(key5)->secret",
			"layered":   "layered_consul(key6#sample_app@development)",
		}

		expectedDefinitions := []interface{}{
			model.VariableDef{Name: "foo", Value: "key1", Destination: model.SECRET, Source: model.VAULT, Context: map[string]string{}},
			model.VariableDef{Name: "bar", Value: "key2", Destination: model.CONFIGMAP, Source: model.CONSUL, Context: map[string]string{}},
			model.VariableDef{Name: "simple", Value: "key3", Destination: model.CONFIGMAP, Source: model.SIMPLE, Context: map[string]string{}},
			model.VariableDef{Name: "config", Value: "key6", Destination: model.CONFIGMAP, Source: model.SIMPLE, Context: map[string]string{}},
			model.VariableDef{Name: "inline", Value: "key7", Destination: model.CONFIGMAP, Source: model.SIMPLE, Context: map[string]string{}},
			model.VariableDef{Name: "reference", Value: "node.name", Destination: model.CONFIGMAP, Source: model.REFERENCE, Context: map[string]string{}},
			model.VariableDef{Name: "simple2", Value: "key4", Destination: model.CONFIGMAP, Source: model.SIMPLE, Context: map[string]string{}},
			model.VariableDef{Name: "simple3", Value: "key5", Destination: model.SECRET, Source: model.SIMPLE, Context: map[string]string{}},
			model.VariableDef{Name: "layered", Value: "key6", Destination: model.CONFIGMAP, Source: model.LAYERED_CONSUL, Context: map[string]string{"appname": "sample_app", "environment": "development"}},
		}

		It("should properly parse variable definitions", func() {
			ret, err := ParseVariableDefinitions(conf)

			Expect(err).NotTo(HaveOccurred())
			Expect(ret).Should(ConsistOf(expectedDefinitions...))
		})
	})

	Context("with invalid variable source", func() {
		conf := map[string]string{
			"foo": "vaulter(key1)",
		}

		It("should throw error", func() {
			_, err := ParseVariableDefinitions(conf)

			Expect(err).To(HaveOccurred())
		})
	})

	Context("with invalid variable destination", func() {
		conf := map[string]string{
			"foo": "vault(key1)->secrets",
		}

		It("should throw error", func() {
			_, err := ParseVariableDefinitions(conf)

			Expect(err).To(HaveOccurred())
		})
	})

	Context("with invalid layered consul value", func() {
		conf := map[string]string{
			"foo": "layered_consul(key1#ssss)->secrets",
		}

		It("should throw error", func() {
			_, err := ParseVariableDefinitions(conf)

			Expect(err).To(HaveOccurred())
		})
	})

	Context("with sample configuration", func() {
		conf := `
LogLevel: debug
Consul:
  Address: consul:8500
  Datacenter: dev
  TlsSkipVerify: true
  Token: 123foo
Vault:
  Address: https://vault:8200
  TlsSkipVerify: true
Application:
  name: app
  namespace: dev
  definitions:
    foo: simple(bar)`

		var testCmd = &cobra.Command{}

		It("should read config without error", func() {
			viper.SetConfigType("yaml")
			err := viper.ReadConfig(strings.NewReader(conf))
			Expect(err).NotTo(HaveOccurred())

			err = Setup(testCmd)
			Expect(err).NotTo(HaveOccurred())
			Expect(viper.GetString("vault.tokenpath")).NotTo(BeEmpty())

			config := Get()
			Expect(config).NotTo(BeNil())
			err = viper.Unmarshal(&config)
			Expect(err).NotTo(HaveOccurred())

			Expect(config).NotTo(BeNil())
			Expect(config.LogLevel).To(Equal("debug"))
			Expect(config.Consul.Address).To(Equal("consul:8500"))
			Expect(config.Consul.Datacenter).To(Equal("dev"))
			Expect(config.Consul.TLSSkipVerify).To(Equal(true))
			Expect(config.Consul.Token).To(Equal("123foo"))
			Expect(config.Vault.Address).To(Equal("https://vault:8200"))
			Expect(config.Vault.TLSSkipVerify).To(Equal(true))
			Expect(config.Application.Name).To(Equal("app"))
			Expect(config.Application.Namespace).To(Equal("dev"))
			Expect(config.Application.Definitions).NotTo(BeEmpty())
			Expect(config.Application.Definitions["foo"]).To(Equal("simple(bar)"))
		})
	})
})
