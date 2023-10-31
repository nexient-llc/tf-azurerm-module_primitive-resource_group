package tests

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TerraTestSuiteResourceGroup struct {
	suite.Suite
	TerraformOptions *terraform.Options
}

type ResourceGroup struct {
	Location string `mapstructure:"location"`
}

// setup to do before any test runs
func (suite *TerraTestSuiteResourceGroup) SetupSuite() {
	inputResourceGroup := ResourceGroup{
		Location: "West US",
	}

	var resourceGroupOptions map[string]interface{}
	err := mapstructure.Decode(inputResourceGroup, &resourceGroupOptions)
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.TerraformOptions = terraform.WithDefaultRetryableErrors(suite.T(), &terraform.Options{
		TerraformDir: "../..",
		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"resource_group":      resourceGroupOptions,
			"resource_group_name": "iac-dev-000-rg-001",
			"tags": map[string]interface{}{
				"provisioner": "terraform",
			},
		},
	})
	terraform.InitAndApplyAndIdempotent(suite.T(), suite.TerraformOptions)
}

func (suite *TerraTestSuiteResourceGroup) TestTerraformAzureRG() {
	_ = files.CopyFile("../../provider.tf", "../test-provider.tf")

	// Run `terraform output` to get the values of output variables
	var outputResourceGroup ResourceGroup
	terraform.OutputStruct(suite.T(), suite.TerraformOptions, "resource_group", &outputResourceGroup)
	expectedName := "iac-dev-000-rg-001"
	location := outputResourceGroup.Location

	// NOTE: "subscriptionID" is overridden by the environment variable "ARM_SUBSCRIPTION_ID". <>
	subscriptionID := ""
	// Assert
	assert.True(suite.T(), azure.ResourceGroupExists(suite.T(), expectedName, subscriptionID))
	actualRG := azure.GetAResourceGroup(suite.T(), expectedName, subscriptionID)
	assert.Equal(suite.T(), expectedName, *actualRG.Name)
	assert.Equal(suite.T(), location, *actualRG.Location)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRunSuite(t *testing.T) {
	suite.Run(t, new(TerraTestSuiteResourceGroup))
}

// TearDownAllSuite has a TearDownSuite method, which will run after all the tests in the suite have been run.
func (suite *TerraTestSuiteResourceGroup) TearDownSuite() {
	terraform.Destroy(suite.T(), suite.TerraformOptions)
	os.Remove("../test-provider.tf")
}
