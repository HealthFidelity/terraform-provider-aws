package cloudformation

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_cloudformation_stack_set_instance", &resource.Sweeper{
		Name: "aws_cloudformation_stack_set_instance",
		F:    sweepStackSetInstances,
	})

	resource.AddTestSweepers("aws_cloudformation_stack_set", &resource.Sweeper{
		Name: "aws_cloudformation_stack_set",
		Dependencies: []string{
			"aws_cloudformation_stack_set_instance",
		},
		F: sweepStackSets,
	})

	resource.AddTestSweepers("aws_cloudformation_stack", &resource.Sweeper{
		Name: "aws_cloudformation_stack",
		Dependencies: []string{
			"aws_cloudformation_stack_set_instance",
		},
		F: sweepStacks,
	})
}

func sweepStackSetInstances(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)

	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}

	conn := client.(*conns.AWSClient).CloudFormationConn
	stackSets, err := ListStackSets(conn)

	if sweep.SkipSweepError(err) || tfawserr.ErrMessageContains(err, "ValidationError", "AWS CloudFormation StackSets is not supported") {
		log.Printf("[WARN] Skipping CloudFormation StackSet Instance sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing CloudFormation StackSets: %w", err)
	}

	var sweeperErrs *multierror.Error

	for _, stackSet := range stackSets {
		stackSetName := aws.StringValue(stackSet.StackSetName)
		instances, err := ListStackSetInstances(conn, stackSetName)

		if err != nil {
			sweeperErr := fmt.Errorf("error listing CloudFormation StackSet (%s) Instances: %w", stackSetName, err)
			log.Printf("[ERROR] %s", sweeperErr)
			sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
			continue
		}

		for _, instance := range instances {
			accountID := aws.StringValue(instance.Account)
			region := aws.StringValue(instance.Region)
			id := fmt.Sprintf("%s / %s / %s", stackSetName, accountID, region)

			input := &cloudformation.DeleteStackInstancesInput{
				Accounts:     aws.StringSlice([]string{accountID}),
				OperationId:  aws.String(resource.UniqueId()),
				Regions:      aws.StringSlice([]string{region}),
				RetainStacks: aws.Bool(false),
				StackSetName: aws.String(stackSetName),
			}

			log.Printf("[INFO] Deleting CloudFormation StackSet Instance: %s", id)
			output, err := conn.DeleteStackInstances(input)

			if tfawserr.ErrMessageContains(err, cloudformation.ErrCodeStackInstanceNotFoundException, "") {
				continue
			}

			if tfawserr.ErrMessageContains(err, cloudformation.ErrCodeStackSetNotFoundException, "") {
				continue
			}

			if err != nil {
				sweeperErr := fmt.Errorf("error deleting CloudFormation StackSet Instance (%s): %w", id, err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
				continue
			}

			if err := WaitStackSetOperationSucceeded(conn, stackSetName, aws.StringValue(output.OperationId), StackSetInstanceDeletedDefaultTimeout); err != nil {
				sweeperErr := fmt.Errorf("error waiting for CloudFormation StackSet Instance (%s) deletion: %w", id, err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
				continue
			}
		}

	}

	return sweeperErrs.ErrorOrNil()
}

func sweepStackSets(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)

	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}

	conn := client.(*conns.AWSClient).CloudFormationConn
	stackSets, err := ListStackSets(conn)

	if sweep.SkipSweepError(err) || tfawserr.ErrMessageContains(err, "ValidationError", "AWS CloudFormation StackSets is not supported") {
		log.Printf("[WARN] Skipping CloudFormation StackSet sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing CloudFormation StackSets: %w", err)
	}

	var sweeperErrs *multierror.Error

	for _, stackSet := range stackSets {
		input := &cloudformation.DeleteStackSetInput{
			StackSetName: stackSet.StackSetName,
		}
		name := aws.StringValue(stackSet.StackSetName)

		log.Printf("[INFO] Deleting CloudFormation StackSet: %s", name)
		_, err := conn.DeleteStackSet(input)

		if tfawserr.ErrMessageContains(err, cloudformation.ErrCodeStackSetNotFoundException, "") {
			continue
		}

		if err != nil {
			sweeperErr := fmt.Errorf("error deleting CloudFormation StackSet (%s): %w", name, err)
			log.Printf("[ERROR] %s", sweeperErr)
			sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
			continue
		}
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepStacks(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)

	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.(*conns.AWSClient).CloudFormationConn
	input := &cloudformation.ListStacksInput{
		StackStatusFilter: aws.StringSlice([]string{
			cloudformation.StackStatusCreateComplete,
			cloudformation.StackStatusImportComplete,
			cloudformation.StackStatusRollbackComplete,
			cloudformation.StackStatusUpdateComplete,
		}),
	}
	var sweeperErrs *multierror.Error

	err = conn.ListStacksPages(input, func(page *cloudformation.ListStacksOutput, lastPage bool) bool {
		for _, stack := range page.StackSummaries {
			input := &cloudformation.DeleteStackInput{
				StackName: stack.StackName,
			}
			name := aws.StringValue(stack.StackName)

			log.Printf("[INFO] Deleting CloudFormation Stack: %s", name)
			_, err := conn.DeleteStack(input)

			if err != nil {
				sweeperErr := fmt.Errorf("error deleting CloudFormation Stack (%s): %w", name, err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
				continue
			}
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping CloudFormation Stack sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing CloudFormation Stacks: %s", err)
	}

	return sweeperErrs.ErrorOrNil()
}
