package servicecatalog_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/servicecatalog"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func TestAccAWSServiceCatalogPortfolioDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix("tf-acc-test")
	dataSourceName := "data.aws_servicecatalog_portfolio.test"
	resourceName := "aws_servicecatalog_portfolio.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, servicecatalog.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckServiceCatlaogPortfolioDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAwsServiceCatalogPortfolioDataSourceConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "arn", dataSourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "created_time", dataSourceName, "created_time"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "provider_name", dataSourceName, "provider_name"),
					resource.TestCheckResourceAttrPair(resourceName, "tags.%", dataSourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(resourceName, "tags.Chicane", dataSourceName, "tags.Chicane"),
				),
			},
		},
	})
}

func testAccCheckAwsServiceCatalogPortfolioDataSourceConfigBasic(rName string) string {
	return acctest.ConfigCompose(testAccCheckAwsServiceCatalogPortfolioResourceConfigTags1(rName, "Chicane", "Nick"), `
data "aws_servicecatalog_portfolio" "test" {
  id = aws_servicecatalog_portfolio.test.id
}
`)
}
