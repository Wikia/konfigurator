package config_test

import (
	. "github.com/Wikia/konfigurator/config"

	"github.com/Wikia/konfigurator/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Context("with sample configuration", func() {
		conf := map[string]string{
			"foo":       "vault(key1)",
			"bar":       "consul(key2)",
			"simple":    "simple(key3)",
			"reference": "simple(node.name)->reference",
			"simple2":   "simple(key4)->config",
			"simple3":   "simple(key5)->secret",
		}

		expectedDefinitions := []interface{}{
			model.VariableDef{Name: "foo", Value: "key1", Type: model.SECRET, Source: model.VAULT},
			model.VariableDef{Name: "bar", Value: "key2", Type: model.CONFIGMAP, Source: model.CONSUL},
			model.VariableDef{Name: "simple", Value: "key3", Type: model.CONFIGMAP, Source: model.SIMPLE},
			model.VariableDef{Name: "reference", Value: "node.name", Type: model.REFERENCE, Source: model.SIMPLE},
			model.VariableDef{Name: "simple2", Value: "key4", Type: model.CONFIGMAP, Source: model.SIMPLE},
			model.VariableDef{Name: "simple3", Value: "key5", Type: model.SECRET, Source: model.SIMPLE},
		}

		It("should properly parse variable definitions", func() {
			ret, err := ParseVariableDefinitions(conf)

			Expect(err).NotTo(HaveOccurred())
			Expect(ret).Should(ConsistOf(expectedDefinitions...))
		})
	})
})
