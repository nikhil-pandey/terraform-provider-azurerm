package componentsapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OperationId{})
}

var _ resourceids.ResourceId = &OperationId{}

// OperationId is a struct representing the Resource ID for a Operation
type OperationId struct {
	SubscriptionId    string
	ResourceGroupName string
	ComponentName     string
	PurgeId           string
}

// NewOperationID returns a new OperationId struct
func NewOperationID(subscriptionId string, resourceGroupName string, componentName string, purgeId string) OperationId {
	return OperationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ComponentName:     componentName,
		PurgeId:           purgeId,
	}
}

// ParseOperationID parses 'input' into a OperationId
func ParseOperationID(input string) (*OperationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OperationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OperationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOperationIDInsensitively parses 'input' case-insensitively into a OperationId
// note: this method should only be used for API response data and not user input
func ParseOperationIDInsensitively(input string) (*OperationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OperationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OperationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OperationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ComponentName, ok = input.Parsed["componentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "componentName", input)
	}

	if id.PurgeId, ok = input.Parsed["purgeId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "purgeId", input)
	}

	return nil
}

// ValidateOperationID checks that 'input' can be parsed as a Operation ID
func ValidateOperationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOperationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Operation ID
func (id OperationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/components/%s/operations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ComponentName, id.PurgeId)
}

// Segments returns a slice of Resource ID Segments which comprise this Operation ID
func (id OperationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticComponents", "components", "components"),
		resourceids.UserSpecifiedSegment("componentName", "componentValue"),
		resourceids.StaticSegment("staticOperations", "operations", "operations"),
		resourceids.UserSpecifiedSegment("purgeId", "purgeIdValue"),
	}
}

// String returns a human-readable description of this Operation ID
func (id OperationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Component Name: %q", id.ComponentName),
		fmt.Sprintf("Purge: %q", id.PurgeId),
	}
	return fmt.Sprintf("Operation (%s)", strings.Join(components, "\n"))
}
