package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/awslabs/goformation"
	"github.com/rodaine/hclencoder"
)

type Resource struct {
	Kind      string `hcl:",key"`
	Name      string `hcl:",key"`
	CidrBlock string `hcl:"cidr_block" hcle:"omitempty"`
}

type Config struct {
	Resources []Resource `hcl:"resource"`
}

var filename = "main.tf"

func main() {
	template, err := goformation.Open("template.json")
	if err != nil {
		log.Fatalf("There was an error processing the template: %s", err)
	}

	resources := make([]Resource, 0)

	vpcs := template.GetAllAWSEC2VPCResources()
	for name, vpc := range vpcs {
		resources = append(resources, Resource{
			Kind:      "aws_vpc",
			Name:      strings.ToLower(name),
			CidrBlock: vpc.CidrBlock,
		})
	}

	input := &Config{Resources: resources}

	hcl, err := hclencoder.Encode(input)
	if err != nil {
		log.Fatal("unable to encode", err)
	}

	if err := ioutil.WriteFile(filename, hcl, 0644); err != nil {
		log.Fatal("unable to write file", err)
	}

	log.Printf("Write file %s\n", filename)
}
