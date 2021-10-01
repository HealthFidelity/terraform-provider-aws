// Code generated by internal/generate/tags/main.go; DO NOT EDIT.

package cloudhsmv2

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudhsmv2"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// ListTags lists cloudhsmv2 service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(conn *cloudhsmv2.CloudHSMV2, identifier string) (tftags.KeyValueTags, error) {
	input := &cloudhsmv2.ListTagsInput{
		ResourceId: aws.String(identifier),
	}

	output, err := conn.ListTags(input)

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.TagList), nil
}

// []*SERVICE.Tag handling

// Tags returns cloudhsmv2 service tags.
func Tags(tags tftags.KeyValueTags) []*cloudhsmv2.Tag {
	result := make([]*cloudhsmv2.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &cloudhsmv2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from cloudhsmv2 service tags.
func KeyValueTags(tags []*cloudhsmv2.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}

// UpdateTags updates cloudhsmv2 service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(conn *cloudhsmv2.CloudHSMV2, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &cloudhsmv2.UntagResourceInput{
			ResourceId: aws.String(identifier),
			TagKeyList: aws.StringSlice(removedTags.IgnoreAws().Keys()),
		}

		_, err := conn.UntagResource(input)

		if err != nil {
			return fmt.Errorf("error untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &cloudhsmv2.TagResourceInput{
			ResourceId: aws.String(identifier),
			TagList:    Tags(updatedTags.IgnoreAws()),
		}

		_, err := conn.TagResource(input)

		if err != nil {
			return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}
