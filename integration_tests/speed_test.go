package integration_tests_test

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/myshkin5/jsonstruct"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type WikipediaSample struct {
	FirstName    string        `json:"firstName"`
	LastName     string        `json:"lastName"`
	IsAlive      bool          `json:"isAlive"`
	Age          int           `json:"age"`
	Address      Address       `json:"address"`
	PhoneNumbers []PhoneNumber `json:"phoneNumbers"`
	Children     []string      `json:"children"`
	Spouse       *string       `json:"spouse"`
}

type Address struct {
	StreetAddress string `json:"streetAddress"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postalCode"`
}

type PhoneNumber struct {
	NumberType string `json:"type"`
	Number     string `json:"number"`
}

var _ = Describe("JSONStruct", func() {
	var (
		buffer []byte
	)

	BeforeEach(func() {
		buffer = loadFromFile("./wikipedia_sample.json")
	})

	Measure("just a read", func(b Benchmarker) {
		// Prime any caches
		parseToWikipediaSample(buffer)
		parseToJSONStruct(buffer)
		parseToWikipediaSample(buffer)
		parseToJSONStruct(buffer)

		var jsonRunDuration, wikipediaRunDuration time.Duration
		for i := 0; i < 100; i++ {
			start := time.Now()
			if i%2 == 0 {
				parseToWikipediaSample(buffer)
				wikipediaRunDuration += time.Now().Sub(start)
			} else {
				parseAndReadJSONStruct(buffer)
				jsonRunDuration += time.Now().Sub(start)
			}
		}

		Expect(jsonRunDuration / 2).To(BeNumerically("<", wikipediaRunDuration))

		b.RecordValue("Wikipedia Sample Struct Run Duration (ns)", float64(wikipediaRunDuration.Nanoseconds()))
		b.RecordValue("JSON Struct Run Duration (ns)", float64(jsonRunDuration.Nanoseconds()))
	}, 10)

	Measure("read followed by five value gets", func(b Benchmarker) {
		// Prime any caches
		parseAndReadWikipediaSample(buffer)
		parseAndReadJSONStruct(buffer)
		parseAndReadWikipediaSample(buffer)
		parseAndReadJSONStruct(buffer)

		var jsonRunDuration, wikipediaRunDuration time.Duration
		for i := 0; i < 100; i++ {
			start := time.Now()
			if i%2 == 0 {
				parseAndReadWikipediaSample(buffer)
				wikipediaRunDuration += time.Now().Sub(start)
			} else {
				parseAndReadJSONStruct(buffer)
				jsonRunDuration += time.Now().Sub(start)
			}
		}

		Expect(jsonRunDuration / 2).To(BeNumerically("<", wikipediaRunDuration))

		b.RecordValue("Wikipedia Sample Struct Run Duration (ns)", float64(wikipediaRunDuration.Nanoseconds()))
		b.RecordValue("JSON Struct Run Duration (ns)", float64(jsonRunDuration.Nanoseconds()))
	}, 10)
})

func loadFromFile(filename string) []byte {
	buffer, err := ioutil.ReadFile(filename)
	Expect(err).NotTo(HaveOccurred())
	return buffer
}

func parseToJSONStruct(buffer []byte) jsonstruct.JSONStruct {
	var jsonStruct jsonstruct.JSONStruct
	err := json.Unmarshal(buffer, &jsonStruct)
	Expect(err).NotTo(HaveOccurred())
	return jsonStruct
}

func parseToWikipediaSample(buffer []byte) WikipediaSample {
	var sample WikipediaSample
	err := json.Unmarshal(buffer, &sample)
	Expect(err).NotTo(HaveOccurred())
	return sample
}

func parseAndReadJSONStruct(buffer []byte) {
	jsonStruct := parseToJSONStruct(buffer)

	firstName, ok := jsonStruct.String("firstName")
	Expect(ok).To(BeTrue())
	Expect(firstName).To(Equal("John"))
	lastName, ok := jsonStruct.String("lastName")
	Expect(ok).To(BeTrue())
	Expect(lastName).To(Equal("Smith"))
	age, ok := jsonStruct.Int("age")
	Expect(ok).To(BeTrue())
	Expect(age).To(Equal(25))
	streetAddress, ok := jsonStruct.String("address.streetAddress")
	Expect(ok).To(BeTrue())
	Expect(streetAddress).To(Equal("21 2nd Street"))
	postalCode, ok := jsonStruct.String("address.postalCode")
	Expect(ok).To(BeTrue())
	Expect(postalCode).To(Equal("10021-3100"))
}

func parseAndReadWikipediaSample(buffer []byte) {
	sample := parseToWikipediaSample(buffer)

	Expect(sample.FirstName).To(Equal("John"))
	Expect(sample.LastName).To(Equal("Smith"))
	Expect(sample.Age).To(Equal(25))
	Expect(sample.Address.StreetAddress).To(Equal("21 2nd Street"))
	Expect(sample.Address.PostalCode).To(Equal("10021-3100"))
}
