package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func TestAccAwsServiceQuotasServiceDataSource_ServiceName(t *testing.T) {
	dataSourceName := "data.aws_servicequotas_service.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(servicequotas.EndpointsID, t) },
		ErrorCheck: acctest.ErrorCheck(t, servicequotas.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsServiceQuotasServiceDataSourceConfigServiceName("Amazon Virtual Private Cloud (Amazon VPC)"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "service_code", "vpc"),
				),
			},
		},
	})
}

func testAccAwsServiceQuotasServiceDataSourceConfigServiceName(serviceName string) string {
	return fmt.Sprintf(`
data "aws_servicequotas_service" "test" {
  service_name = %[1]q
}
`, serviceName)
}
