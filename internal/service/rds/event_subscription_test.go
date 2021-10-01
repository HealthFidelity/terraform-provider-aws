package rds_test

import (
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/go-multierror"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfrds "github.com/hashicorp/terraform-provider-aws/internal/service/rds"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)





func TestAccRDSEventSubscription_basicUpdate(t *testing.T) {
	var v rds.EventSubscription
	rInt := sdkacctest.RandInt()
	resourceName := "aws_db_event_subscription.test"
	subscriptionName := fmt.Sprintf("tf-acc-test-rds-event-subs-%d", rInt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, rds.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventSubscriptionConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(resourceName, &v),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "rds", regexp.MustCompile(fmt.Sprintf("es:%s$", subscriptionName))),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_type", "db-instance"),
					resource.TestCheckResourceAttr(resourceName, "name", subscriptionName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     subscriptionName,
			},
			{
				Config: testAccEventSubscriptionUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "source_type", "db-parameter-group"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "new-name"),
				),
			},
		},
	})
}

func TestAccRDSEventSubscription_disappears(t *testing.T) {
	var eventSubscription rds.EventSubscription
	rInt := sdkacctest.RandInt()
	resourceName := "aws_db_event_subscription.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, rds.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventSubscriptionConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(resourceName, &eventSubscription),
					testAccCheckEventSubscriptionDisappears(&eventSubscription),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRDSEventSubscription_withPrefix(t *testing.T) {
	var v rds.EventSubscription
	rInt := sdkacctest.RandInt()
	startsWithPrefix := regexp.MustCompile("^tf-acc-test-rds-event-subs-")
	resourceName := "aws_db_event_subscription.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, rds.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventSubscriptionWithPrefixConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_type", "db-instance"),
					resource.TestMatchResourceAttr(resourceName, "name", startsWithPrefix),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", "name"),
				),
			},
		},
	})
}

func TestAccRDSEventSubscription_withSourceIDs(t *testing.T) {
	var v rds.EventSubscription
	rInt := sdkacctest.RandInt()
	resourceName := "aws_db_event_subscription.test"
	subscriptionName := fmt.Sprintf("tf-acc-test-rds-event-subs-with-ids-%d", rInt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, rds.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventSubscriptionWithSourceIDsConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_type", "db-parameter-group"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-acc-test-rds-event-subs-with-ids-%d", rInt)),
					resource.TestCheckResourceAttr(resourceName, "source_ids.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     subscriptionName,
			},
			{
				Config: testAccEventSubscriptionUpdateSourceIDsConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_type", "db-parameter-group"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-acc-test-rds-event-subs-with-ids-%d", rInt)),
					resource.TestCheckResourceAttr(resourceName, "source_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccRDSEventSubscription_categoryUpdate(t *testing.T) {
	var v rds.EventSubscription
	rInt := sdkacctest.RandInt()
	resourceName := "aws_db_event_subscription.test"
	subscriptionName := fmt.Sprintf("tf-acc-test-rds-event-subs-%d", rInt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, rds.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckEventSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventSubscriptionConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_type", "db-instance"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-acc-test-rds-event-subs-%d", rInt)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     subscriptionName,
			},
			{
				Config: testAccEventSubscriptionUpdateCategoriesConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_type", "db-instance"),
				),
			},
		},
	})
}

func testAccCheckEventSubscriptionExists(n string, v *rds.EventSubscription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No RDS Event Subscription is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).RDSConn

		eventSubscription, err := tfrds.EventSubscriptionRetrieve(rs.Primary.ID, conn)

		if err != nil {
			return err
		}

		if eventSubscription == nil {
			return fmt.Errorf("RDS Event Subscription not found")
		}

		*v = *eventSubscription

		return nil
	}
}

func testAccCheckEventSubscriptionDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).RDSConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_db_event_subscription" {
			continue
		}

		eventSubscription, err := tfrds.EventSubscriptionRetrieve(rs.Primary.ID, conn)

		if tfawserr.ErrMessageContains(err, rds.ErrCodeSubscriptionNotFoundFault, "") {
			continue
		}

		if err != nil {
			return err
		}

		if eventSubscription != nil {
			return fmt.Errorf("RDS Event Subscription (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckEventSubscriptionDisappears(eventSubscription *rds.EventSubscription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).RDSConn

		input := &rds.DeleteEventSubscriptionInput{
			SubscriptionName: eventSubscription.CustSubscriptionId,
		}

		_, err := conn.DeleteEventSubscription(input)

		if err != nil {
			return err
		}

		return tfrds.WaitForEventSubscriptionDeletion(conn, aws.StringValue(eventSubscription.CustSubscriptionId), 10*time.Minute)
	}
}

func testAccEventSubscriptionConfig(rInt int) string {
	return fmt.Sprintf(`
resource "aws_sns_topic" "aws_sns_topic" {
  name = "tf-acc-test-rds-event-subs-sns-topic-%[1]d"
}

resource "aws_db_event_subscription" "test" {
  name        = "tf-acc-test-rds-event-subs-%[1]d"
  sns_topic   = aws_sns_topic.aws_sns_topic.arn
  source_type = "db-instance"

  event_categories = [
    "availability",
    "backup",
    "creation",
    "deletion",
    "maintenance",
  ]

  tags = {
    Name = "name"
  }
}
`, rInt)
}

func testAccEventSubscriptionWithPrefixConfig(rInt int) string {
	return fmt.Sprintf(`
resource "aws_sns_topic" "aws_sns_topic" {
  name = "tf-acc-test-rds-event-subs-sns-topic-%d"
}

resource "aws_db_event_subscription" "test" {
  name_prefix = "tf-acc-test-rds-event-subs-"
  sns_topic   = aws_sns_topic.aws_sns_topic.arn
  source_type = "db-instance"

  event_categories = [
    "availability",
    "backup",
    "creation",
    "deletion",
    "maintenance",
  ]

  tags = {
    Name = "name"
  }
}
`, rInt)
}

func testAccEventSubscriptionUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "aws_sns_topic" "aws_sns_topic" {
  name = "tf-acc-test-rds-event-subs-sns-topic-%[1]d"
}

resource "aws_db_event_subscription" "test" {
  name        = "tf-acc-test-rds-event-subs-%[1]d"
  sns_topic   = aws_sns_topic.aws_sns_topic.arn
  enabled     = false
  source_type = "db-parameter-group"

  event_categories = [
    "configuration change",
  ]

  tags = {
    Name = "new-name"
  }
}
`, rInt)
}

func testAccEventSubscriptionWithSourceIDsConfig(rInt int) string {
	return fmt.Sprintf(`
resource "aws_sns_topic" "aws_sns_topic" {
  name = "tf-acc-test-rds-event-subs-sns-topic-%[1]d"
}

resource "aws_db_parameter_group" "test" {
  name        = "db-parameter-group-event-%[1]d"
  family      = "mysql5.6"
  description = "Test parameter group for terraform"
}

resource "aws_db_event_subscription" "test" {
  name        = "tf-acc-test-rds-event-subs-with-ids-%[1]d"
  sns_topic   = aws_sns_topic.aws_sns_topic.arn
  source_type = "db-parameter-group"
  source_ids  = [aws_db_parameter_group.test.id]

  event_categories = [
    "configuration change",
  ]

  tags = {
    Name = "name"
  }
}
`, rInt)
}

func testAccEventSubscriptionUpdateSourceIDsConfig(rInt int) string {
	return fmt.Sprintf(`
resource "aws_sns_topic" "aws_sns_topic" {
  name = "tf-acc-test-rds-event-subs-sns-topic-%[1]d"
}

resource "aws_db_parameter_group" "test" {
  name        = "db-parameter-group-event-%[1]d"
  family      = "mysql5.6"
  description = "Test parameter group for terraform"
}

resource "aws_db_parameter_group" "test2" {
  name        = "db-parameter-group-event-2-%[1]d"
  family      = "mysql5.6"
  description = "Test parameter group for terraform"
}

resource "aws_db_event_subscription" "test" {
  name        = "tf-acc-test-rds-event-subs-with-ids-%[1]d"
  sns_topic   = aws_sns_topic.aws_sns_topic.arn
  source_type = "db-parameter-group"
  source_ids  = [aws_db_parameter_group.test.id, aws_db_parameter_group.test2.id]

  event_categories = [
    "configuration change",
  ]

  tags = {
    Name = "name"
  }
}
`, rInt)
}

func testAccEventSubscriptionUpdateCategoriesConfig(rInt int) string {
	return fmt.Sprintf(`
resource "aws_sns_topic" "aws_sns_topic" {
  name = "tf-acc-test-rds-event-subs-sns-topic-%[1]d"
}

resource "aws_db_event_subscription" "test" {
  name        = "tf-acc-test-rds-event-subs-%[1]d"
  sns_topic   = aws_sns_topic.aws_sns_topic.arn
  source_type = "db-instance"

  event_categories = [
    "availability",
  ]

  tags = {
    Name = "name"
  }
}
`, rInt)
}
