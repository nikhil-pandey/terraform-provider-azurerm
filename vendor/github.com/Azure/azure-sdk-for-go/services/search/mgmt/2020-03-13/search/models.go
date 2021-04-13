package search

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"encoding/json"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// The package's fully qualified name.
const fqdn = "github.com/Azure/azure-sdk-for-go/services/search/mgmt/2020-03-13/search"

// AdminKeyResult response containing the primary and secondary admin API keys for a given Azure Cognitive
// Search service.
type AdminKeyResult struct {
	autorest.Response `json:"-"`
	// PrimaryKey - READ-ONLY; The primary admin API key of the Search service.
	PrimaryKey *string `json:"primaryKey,omitempty"`
	// SecondaryKey - READ-ONLY; The secondary admin API key of the Search service.
	SecondaryKey *string `json:"secondaryKey,omitempty"`
}

// CheckNameAvailabilityInput input of check name availability API.
type CheckNameAvailabilityInput struct {
	// Name - The Search service name to validate. Search service names must only contain lowercase letters, digits or dashes, cannot use dash as the first two or last one characters, cannot contain consecutive dashes, and must be between 2 and 60 characters in length.
	Name *string `json:"name,omitempty"`
	// Type - The type of the resource whose name is to be validated. This value must always be 'searchServices'.
	Type *string `json:"type,omitempty"`
}

// CheckNameAvailabilityOutput output of check name availability API.
type CheckNameAvailabilityOutput struct {
	autorest.Response `json:"-"`
	// IsNameAvailable - READ-ONLY; A value indicating whether the name is available.
	IsNameAvailable *bool `json:"nameAvailable,omitempty"`
	// Reason - READ-ONLY; The reason why the name is not available. 'Invalid' indicates the name provided does not match the naming requirements (incorrect length, unsupported characters, etc.). 'AlreadyExists' indicates that the name is already in use and is therefore unavailable. Possible values include: 'Invalid', 'AlreadyExists'
	Reason UnavailableNameReason `json:"reason,omitempty"`
	// Message - READ-ONLY; A message that explains why the name is invalid and provides resource naming requirements. Available only if 'Invalid' is returned in the 'reason' property.
	Message *string `json:"message,omitempty"`
}

// CloudError contains information about an API error.
type CloudError struct {
	// Error - Describes a particular API error with an error code and a message.
	Error *CloudErrorBody `json:"error,omitempty"`
}

// CloudErrorBody describes a particular API error with an error code and a message.
type CloudErrorBody struct {
	// Code - An error code that describes the error condition more precisely than an HTTP status code. Can be used to programmatically handle specific error cases.
	Code *string `json:"code,omitempty"`
	// Message - A message that describes the error in detail and provides debugging information.
	Message *string `json:"message,omitempty"`
	// Target - The target of the particular error (for example, the name of the property in error).
	Target *string `json:"target,omitempty"`
	// Details - Contains nested errors that are related to this error.
	Details *[]CloudErrorBody `json:"details,omitempty"`
}

// Identity identity for the resource.
type Identity struct {
	// PrincipalID - READ-ONLY; The principal ID of resource identity.
	PrincipalID *string `json:"principalId,omitempty"`
	// TenantID - READ-ONLY; The tenant ID of resource.
	TenantID *string `json:"tenantId,omitempty"`
	// Type - The identity type. Possible values include: 'None', 'SystemAssigned'
	Type IdentityType `json:"type,omitempty"`
}

// MarshalJSON is the custom marshaler for Identity.
func (i Identity) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if i.Type != "" {
		objectMap["type"] = i.Type
	}
	return json.Marshal(objectMap)
}

// IPRule the IP restriction rule of the Azure Cognitive Search service.
type IPRule struct {
	// Value - Value corresponding to a single IPv4 address (eg., 123.1.2.3) or an IP range in CIDR format (eg., 123.1.2.3/24) to be allowed.
	Value *string `json:"value,omitempty"`
}

// ListQueryKeysResult response containing the query API keys for a given Azure Cognitive Search service.
type ListQueryKeysResult struct {
	autorest.Response `json:"-"`
	// Value - READ-ONLY; The query keys for the Azure Cognitive Search service.
	Value *[]QueryKey `json:"value,omitempty"`
	// NextLink - READ-ONLY; Request URL that can be used to query next page of query keys. Returned when the total number of requested query keys exceed maximum page size.
	NextLink *string `json:"nextLink,omitempty"`
}

// ListQueryKeysResultIterator provides access to a complete listing of QueryKey values.
type ListQueryKeysResultIterator struct {
	i    int
	page ListQueryKeysResultPage
}

// NextWithContext advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
func (iter *ListQueryKeysResultIterator) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ListQueryKeysResultIterator.NextWithContext")
		defer func() {
			sc := -1
			if iter.Response().Response.Response != nil {
				sc = iter.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	iter.i++
	if iter.i < len(iter.page.Values()) {
		return nil
	}
	err = iter.page.NextWithContext(ctx)
	if err != nil {
		iter.i--
		return err
	}
	iter.i = 0
	return nil
}

// Next advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
// Deprecated: Use NextWithContext() instead.
func (iter *ListQueryKeysResultIterator) Next() error {
	return iter.NextWithContext(context.Background())
}

// NotDone returns true if the enumeration should be started or is not yet complete.
func (iter ListQueryKeysResultIterator) NotDone() bool {
	return iter.page.NotDone() && iter.i < len(iter.page.Values())
}

// Response returns the raw server response from the last page request.
func (iter ListQueryKeysResultIterator) Response() ListQueryKeysResult {
	return iter.page.Response()
}

// Value returns the current value or a zero-initialized value if the
// iterator has advanced beyond the end of the collection.
func (iter ListQueryKeysResultIterator) Value() QueryKey {
	if !iter.page.NotDone() {
		return QueryKey{}
	}
	return iter.page.Values()[iter.i]
}

// Creates a new instance of the ListQueryKeysResultIterator type.
func NewListQueryKeysResultIterator(page ListQueryKeysResultPage) ListQueryKeysResultIterator {
	return ListQueryKeysResultIterator{page: page}
}

// IsEmpty returns true if the ListResult contains no values.
func (lqkr ListQueryKeysResult) IsEmpty() bool {
	return lqkr.Value == nil || len(*lqkr.Value) == 0
}

// hasNextLink returns true if the NextLink is not empty.
func (lqkr ListQueryKeysResult) hasNextLink() bool {
	return lqkr.NextLink != nil && len(*lqkr.NextLink) != 0
}

// listQueryKeysResultPreparer prepares a request to retrieve the next set of results.
// It returns nil if no more results exist.
func (lqkr ListQueryKeysResult) listQueryKeysResultPreparer(ctx context.Context) (*http.Request, error) {
	if !lqkr.hasNextLink() {
		return nil, nil
	}
	return autorest.Prepare((&http.Request{}).WithContext(ctx),
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(lqkr.NextLink)))
}

// ListQueryKeysResultPage contains a page of QueryKey values.
type ListQueryKeysResultPage struct {
	fn   func(context.Context, ListQueryKeysResult) (ListQueryKeysResult, error)
	lqkr ListQueryKeysResult
}

// NextWithContext advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *ListQueryKeysResultPage) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ListQueryKeysResultPage.NextWithContext")
		defer func() {
			sc := -1
			if page.Response().Response.Response != nil {
				sc = page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	for {
		next, err := page.fn(ctx, page.lqkr)
		if err != nil {
			return err
		}
		page.lqkr = next
		if !next.hasNextLink() || !next.IsEmpty() {
			break
		}
	}
	return nil
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
// Deprecated: Use NextWithContext() instead.
func (page *ListQueryKeysResultPage) Next() error {
	return page.NextWithContext(context.Background())
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page ListQueryKeysResultPage) NotDone() bool {
	return !page.lqkr.IsEmpty()
}

// Response returns the raw server response from the last page request.
func (page ListQueryKeysResultPage) Response() ListQueryKeysResult {
	return page.lqkr
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page ListQueryKeysResultPage) Values() []QueryKey {
	if page.lqkr.IsEmpty() {
		return nil
	}
	return *page.lqkr.Value
}

// Creates a new instance of the ListQueryKeysResultPage type.
func NewListQueryKeysResultPage(cur ListQueryKeysResult, getNextPage func(context.Context, ListQueryKeysResult) (ListQueryKeysResult, error)) ListQueryKeysResultPage {
	return ListQueryKeysResultPage{
		fn:   getNextPage,
		lqkr: cur,
	}
}

// NetworkRuleSet network specific rules that determine how the Azure Cognitive Search service may be
// reached.
type NetworkRuleSet struct {
	// IPRules - A list of IP restriction rules that defines the inbound network(s) with allowing access to the search service endpoint. At the meantime, all other public IP networks are blocked by the firewall. These restriction rules are applied only when the 'publicNetworkAccess' of the search service is 'enabled'; otherwise, traffic over public interface is not allowed even with any public IP rules, and private endpoint connections would be the exclusive access method.
	IPRules *[]IPRule `json:"ipRules,omitempty"`
}

// Operation describes a REST API operation.
type Operation struct {
	// Name - READ-ONLY; The name of the operation. This name is of the form {provider}/{resource}/{operation}.
	Name *string `json:"name,omitempty"`
	// Display - READ-ONLY; The object that describes the operation.
	Display *OperationDisplay `json:"display,omitempty"`
}

// OperationDisplay the object that describes the operation.
type OperationDisplay struct {
	// Provider - READ-ONLY; The friendly name of the resource provider.
	Provider *string `json:"provider,omitempty"`
	// Operation - READ-ONLY; The operation type: read, write, delete, listKeys/action, etc.
	Operation *string `json:"operation,omitempty"`
	// Resource - READ-ONLY; The resource type on which the operation is performed.
	Resource *string `json:"resource,omitempty"`
	// Description - READ-ONLY; The friendly name of the operation.
	Description *string `json:"description,omitempty"`
}

// OperationListResult the result of the request to list REST API operations. It contains a list of
// operations and a URL  to get the next set of results.
type OperationListResult struct {
	autorest.Response `json:"-"`
	// Value - READ-ONLY; The list of operations supported by the resource provider.
	Value *[]Operation `json:"value,omitempty"`
	// NextLink - READ-ONLY; The URL to get the next set of operation list results, if any.
	NextLink *string `json:"nextLink,omitempty"`
}

// PrivateEndpointConnection describes an existing Private Endpoint connection to the Azure Cognitive
// Search service.
type PrivateEndpointConnection struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The ID of the private endpoint connection. This can be used with the Azure Resource Manager to link resources together.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the private endpoint connection.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The resource type.
	Type *string `json:"type,omitempty"`
	// Properties - Describes the properties of an existing Private Endpoint connection to the Azure Cognitive Search service.
	Properties *PrivateEndpointConnectionProperties `json:"properties,omitempty"`
}

// MarshalJSON is the custom marshaler for PrivateEndpointConnection.
func (pec PrivateEndpointConnection) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if pec.Properties != nil {
		objectMap["properties"] = pec.Properties
	}
	return json.Marshal(objectMap)
}

// PrivateEndpointConnectionListResult response containing a list of Private Endpoint connections.
type PrivateEndpointConnectionListResult struct {
	autorest.Response `json:"-"`
	// Value - READ-ONLY; The list of Private Endpoint connections.
	Value *[]PrivateEndpointConnection `json:"value,omitempty"`
	// NextLink - READ-ONLY; Request URL that can be used to query next page of private endpoint connections. Returned when the total number of requested private endpoint connections exceed maximum page size.
	NextLink *string `json:"nextLink,omitempty"`
}

// PrivateEndpointConnectionListResultIterator provides access to a complete listing of
// PrivateEndpointConnection values.
type PrivateEndpointConnectionListResultIterator struct {
	i    int
	page PrivateEndpointConnectionListResultPage
}

// NextWithContext advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
func (iter *PrivateEndpointConnectionListResultIterator) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateEndpointConnectionListResultIterator.NextWithContext")
		defer func() {
			sc := -1
			if iter.Response().Response.Response != nil {
				sc = iter.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	iter.i++
	if iter.i < len(iter.page.Values()) {
		return nil
	}
	err = iter.page.NextWithContext(ctx)
	if err != nil {
		iter.i--
		return err
	}
	iter.i = 0
	return nil
}

// Next advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
// Deprecated: Use NextWithContext() instead.
func (iter *PrivateEndpointConnectionListResultIterator) Next() error {
	return iter.NextWithContext(context.Background())
}

// NotDone returns true if the enumeration should be started or is not yet complete.
func (iter PrivateEndpointConnectionListResultIterator) NotDone() bool {
	return iter.page.NotDone() && iter.i < len(iter.page.Values())
}

// Response returns the raw server response from the last page request.
func (iter PrivateEndpointConnectionListResultIterator) Response() PrivateEndpointConnectionListResult {
	return iter.page.Response()
}

// Value returns the current value or a zero-initialized value if the
// iterator has advanced beyond the end of the collection.
func (iter PrivateEndpointConnectionListResultIterator) Value() PrivateEndpointConnection {
	if !iter.page.NotDone() {
		return PrivateEndpointConnection{}
	}
	return iter.page.Values()[iter.i]
}

// Creates a new instance of the PrivateEndpointConnectionListResultIterator type.
func NewPrivateEndpointConnectionListResultIterator(page PrivateEndpointConnectionListResultPage) PrivateEndpointConnectionListResultIterator {
	return PrivateEndpointConnectionListResultIterator{page: page}
}

// IsEmpty returns true if the ListResult contains no values.
func (peclr PrivateEndpointConnectionListResult) IsEmpty() bool {
	return peclr.Value == nil || len(*peclr.Value) == 0
}

// hasNextLink returns true if the NextLink is not empty.
func (peclr PrivateEndpointConnectionListResult) hasNextLink() bool {
	return peclr.NextLink != nil && len(*peclr.NextLink) != 0
}

// privateEndpointConnectionListResultPreparer prepares a request to retrieve the next set of results.
// It returns nil if no more results exist.
func (peclr PrivateEndpointConnectionListResult) privateEndpointConnectionListResultPreparer(ctx context.Context) (*http.Request, error) {
	if !peclr.hasNextLink() {
		return nil, nil
	}
	return autorest.Prepare((&http.Request{}).WithContext(ctx),
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(peclr.NextLink)))
}

// PrivateEndpointConnectionListResultPage contains a page of PrivateEndpointConnection values.
type PrivateEndpointConnectionListResultPage struct {
	fn    func(context.Context, PrivateEndpointConnectionListResult) (PrivateEndpointConnectionListResult, error)
	peclr PrivateEndpointConnectionListResult
}

// NextWithContext advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *PrivateEndpointConnectionListResultPage) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/PrivateEndpointConnectionListResultPage.NextWithContext")
		defer func() {
			sc := -1
			if page.Response().Response.Response != nil {
				sc = page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	for {
		next, err := page.fn(ctx, page.peclr)
		if err != nil {
			return err
		}
		page.peclr = next
		if !next.hasNextLink() || !next.IsEmpty() {
			break
		}
	}
	return nil
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
// Deprecated: Use NextWithContext() instead.
func (page *PrivateEndpointConnectionListResultPage) Next() error {
	return page.NextWithContext(context.Background())
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page PrivateEndpointConnectionListResultPage) NotDone() bool {
	return !page.peclr.IsEmpty()
}

// Response returns the raw server response from the last page request.
func (page PrivateEndpointConnectionListResultPage) Response() PrivateEndpointConnectionListResult {
	return page.peclr
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page PrivateEndpointConnectionListResultPage) Values() []PrivateEndpointConnection {
	if page.peclr.IsEmpty() {
		return nil
	}
	return *page.peclr.Value
}

// Creates a new instance of the PrivateEndpointConnectionListResultPage type.
func NewPrivateEndpointConnectionListResultPage(cur PrivateEndpointConnectionListResult, getNextPage func(context.Context, PrivateEndpointConnectionListResult) (PrivateEndpointConnectionListResult, error)) PrivateEndpointConnectionListResultPage {
	return PrivateEndpointConnectionListResultPage{
		fn:    getNextPage,
		peclr: cur,
	}
}

// PrivateEndpointConnectionProperties describes the properties of an existing Private Endpoint connection
// to the Azure Cognitive Search service.
type PrivateEndpointConnectionProperties struct {
	// PrivateEndpoint - The private endpoint resource from Microsoft.Network provider.
	PrivateEndpoint *PrivateEndpointConnectionPropertiesPrivateEndpoint `json:"privateEndpoint,omitempty"`
	// PrivateLinkServiceConnectionState - Describes the current state of an existing Private Link Service connection to the Azure Private Endpoint.
	PrivateLinkServiceConnectionState *PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState `json:"privateLinkServiceConnectionState,omitempty"`
}

// PrivateEndpointConnectionPropertiesPrivateEndpoint the private endpoint resource from Microsoft.Network
// provider.
type PrivateEndpointConnectionPropertiesPrivateEndpoint struct {
	// ID - The resource id of the private endpoint resource from Microsoft.Network provider.
	ID *string `json:"id,omitempty"`
}

// PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState describes the current state of an
// existing Private Link Service connection to the Azure Private Endpoint.
type PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState struct {
	// Status - Status of the the private link service connection. Can be Pending, Approved, Rejected, or Disconnected. Possible values include: 'Pending', 'Approved', 'Rejected', 'Disconnected'
	Status PrivateLinkServiceConnectionStatus `json:"status,omitempty"`
	// Description - The description for the private link service connection state.
	Description *string `json:"description,omitempty"`
	// ActionsRequired - A description of any extra actions that may be required.
	ActionsRequired *string `json:"actionsRequired,omitempty"`
}

// PrivateLinkResource describes a supported private link resource for the Azure Cognitive Search service.
type PrivateLinkResource struct {
	// ID - READ-ONLY; The ID of the private link resource.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the private link resource.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The resource type.
	Type *string `json:"type,omitempty"`
	// Properties - READ-ONLY; Describes the properties of a supported private link resource for the Azure Cognitive Search service.
	Properties *PrivateLinkResourceProperties `json:"properties,omitempty"`
}

// PrivateLinkResourceProperties describes the properties of a supported private link resource for the
// Azure Cognitive Search service.
type PrivateLinkResourceProperties struct {
	// GroupID - READ-ONLY; The group ID of the private link resource.
	GroupID *string `json:"groupId,omitempty"`
	// RequiredMembers - READ-ONLY; The list of required members of the private link resource.
	RequiredMembers *[]string `json:"requiredMembers,omitempty"`
	// RequiredZoneNames - READ-ONLY; The list of required DNS zone names of the private link resource.
	RequiredZoneNames *[]string `json:"requiredZoneNames,omitempty"`
}

// PrivateLinkResourcesResult response containing a list of supported Private Link Resources.
type PrivateLinkResourcesResult struct {
	autorest.Response `json:"-"`
	// Value - READ-ONLY; The list of supported Private Link Resources.
	Value *[]PrivateLinkResource `json:"value,omitempty"`
}

// QueryKey describes an API key for a given Azure Cognitive Search service that has permissions for query
// operations only.
type QueryKey struct {
	autorest.Response `json:"-"`
	// Name - READ-ONLY; The name of the query API key; may be empty.
	Name *string `json:"name,omitempty"`
	// Key - READ-ONLY; The value of the query API key.
	Key *string `json:"key,omitempty"`
}

// Resource base type for all Azure resources.
type Resource struct {
	// ID - READ-ONLY; The ID of the resource. This can be used with the Azure Resource Manager to link resources together.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the resource.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The resource type.
	Type *string `json:"type,omitempty"`
	// Location - The geographic location of the resource. This must be one of the supported and registered Azure Geo Regions (for example, West US, East US, Southeast Asia, and so forth). This property is required when creating a new resource.
	Location *string `json:"location,omitempty"`
	// Tags - Tags to help categorize the resource in the Azure portal.
	Tags map[string]*string `json:"tags"`
	// Identity - The identity of the resource.
	Identity *Identity `json:"identity,omitempty"`
}

// MarshalJSON is the custom marshaler for Resource.
func (r Resource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if r.Location != nil {
		objectMap["location"] = r.Location
	}
	if r.Tags != nil {
		objectMap["tags"] = r.Tags
	}
	if r.Identity != nil {
		objectMap["identity"] = r.Identity
	}
	return json.Marshal(objectMap)
}

// Service describes an Azure Cognitive Search service and its current state.
type Service struct {
	autorest.Response `json:"-"`
	// ServiceProperties - Properties of the Search service.
	*ServiceProperties `json:"properties,omitempty"`
	// Sku - The SKU of the Search Service, which determines price tier and capacity limits. This property is required when creating a new Search Service.
	Sku *Sku `json:"sku,omitempty"`
	// ID - READ-ONLY; The ID of the resource. This can be used with the Azure Resource Manager to link resources together.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The name of the resource.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The resource type.
	Type *string `json:"type,omitempty"`
	// Location - The geographic location of the resource. This must be one of the supported and registered Azure Geo Regions (for example, West US, East US, Southeast Asia, and so forth). This property is required when creating a new resource.
	Location *string `json:"location,omitempty"`
	// Tags - Tags to help categorize the resource in the Azure portal.
	Tags map[string]*string `json:"tags"`
	// Identity - The identity of the resource.
	Identity *Identity `json:"identity,omitempty"`
}

// MarshalJSON is the custom marshaler for Service.
func (s Service) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if s.ServiceProperties != nil {
		objectMap["properties"] = s.ServiceProperties
	}
	if s.Sku != nil {
		objectMap["sku"] = s.Sku
	}
	if s.Location != nil {
		objectMap["location"] = s.Location
	}
	if s.Tags != nil {
		objectMap["tags"] = s.Tags
	}
	if s.Identity != nil {
		objectMap["identity"] = s.Identity
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON is the custom unmarshaler for Service struct.
func (s *Service) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "properties":
			if v != nil {
				var serviceProperties ServiceProperties
				err = json.Unmarshal(*v, &serviceProperties)
				if err != nil {
					return err
				}
				s.ServiceProperties = &serviceProperties
			}
		case "sku":
			if v != nil {
				var sku Sku
				err = json.Unmarshal(*v, &sku)
				if err != nil {
					return err
				}
				s.Sku = &sku
			}
		case "id":
			if v != nil {
				var ID string
				err = json.Unmarshal(*v, &ID)
				if err != nil {
					return err
				}
				s.ID = &ID
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				s.Name = &name
			}
		case "type":
			if v != nil {
				var typeVar string
				err = json.Unmarshal(*v, &typeVar)
				if err != nil {
					return err
				}
				s.Type = &typeVar
			}
		case "location":
			if v != nil {
				var location string
				err = json.Unmarshal(*v, &location)
				if err != nil {
					return err
				}
				s.Location = &location
			}
		case "tags":
			if v != nil {
				var tags map[string]*string
				err = json.Unmarshal(*v, &tags)
				if err != nil {
					return err
				}
				s.Tags = tags
			}
		case "identity":
			if v != nil {
				var identity Identity
				err = json.Unmarshal(*v, &identity)
				if err != nil {
					return err
				}
				s.Identity = &identity
			}
		}
	}

	return nil
}

// ServiceListResult response containing a list of Azure Cognitive Search services.
type ServiceListResult struct {
	autorest.Response `json:"-"`
	// Value - READ-ONLY; The list of Search services.
	Value *[]Service `json:"value,omitempty"`
	// NextLink - READ-ONLY; Request URL that can be used to query next page of search services. Returned when the total number of requested search services exceed maximum page size.
	NextLink *string `json:"nextLink,omitempty"`
}

// ServiceListResultIterator provides access to a complete listing of Service values.
type ServiceListResultIterator struct {
	i    int
	page ServiceListResultPage
}

// NextWithContext advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
func (iter *ServiceListResultIterator) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServiceListResultIterator.NextWithContext")
		defer func() {
			sc := -1
			if iter.Response().Response.Response != nil {
				sc = iter.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	iter.i++
	if iter.i < len(iter.page.Values()) {
		return nil
	}
	err = iter.page.NextWithContext(ctx)
	if err != nil {
		iter.i--
		return err
	}
	iter.i = 0
	return nil
}

// Next advances to the next value.  If there was an error making
// the request the iterator does not advance and the error is returned.
// Deprecated: Use NextWithContext() instead.
func (iter *ServiceListResultIterator) Next() error {
	return iter.NextWithContext(context.Background())
}

// NotDone returns true if the enumeration should be started or is not yet complete.
func (iter ServiceListResultIterator) NotDone() bool {
	return iter.page.NotDone() && iter.i < len(iter.page.Values())
}

// Response returns the raw server response from the last page request.
func (iter ServiceListResultIterator) Response() ServiceListResult {
	return iter.page.Response()
}

// Value returns the current value or a zero-initialized value if the
// iterator has advanced beyond the end of the collection.
func (iter ServiceListResultIterator) Value() Service {
	if !iter.page.NotDone() {
		return Service{}
	}
	return iter.page.Values()[iter.i]
}

// Creates a new instance of the ServiceListResultIterator type.
func NewServiceListResultIterator(page ServiceListResultPage) ServiceListResultIterator {
	return ServiceListResultIterator{page: page}
}

// IsEmpty returns true if the ListResult contains no values.
func (slr ServiceListResult) IsEmpty() bool {
	return slr.Value == nil || len(*slr.Value) == 0
}

// hasNextLink returns true if the NextLink is not empty.
func (slr ServiceListResult) hasNextLink() bool {
	return slr.NextLink != nil && len(*slr.NextLink) != 0
}

// serviceListResultPreparer prepares a request to retrieve the next set of results.
// It returns nil if no more results exist.
func (slr ServiceListResult) serviceListResultPreparer(ctx context.Context) (*http.Request, error) {
	if !slr.hasNextLink() {
		return nil, nil
	}
	return autorest.Prepare((&http.Request{}).WithContext(ctx),
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(to.String(slr.NextLink)))
}

// ServiceListResultPage contains a page of Service values.
type ServiceListResultPage struct {
	fn  func(context.Context, ServiceListResult) (ServiceListResult, error)
	slr ServiceListResult
}

// NextWithContext advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
func (page *ServiceListResultPage) NextWithContext(ctx context.Context) (err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServiceListResultPage.NextWithContext")
		defer func() {
			sc := -1
			if page.Response().Response.Response != nil {
				sc = page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	for {
		next, err := page.fn(ctx, page.slr)
		if err != nil {
			return err
		}
		page.slr = next
		if !next.hasNextLink() || !next.IsEmpty() {
			break
		}
	}
	return nil
}

// Next advances to the next page of values.  If there was an error making
// the request the page does not advance and the error is returned.
// Deprecated: Use NextWithContext() instead.
func (page *ServiceListResultPage) Next() error {
	return page.NextWithContext(context.Background())
}

// NotDone returns true if the page enumeration should be started or is not yet complete.
func (page ServiceListResultPage) NotDone() bool {
	return !page.slr.IsEmpty()
}

// Response returns the raw server response from the last page request.
func (page ServiceListResultPage) Response() ServiceListResult {
	return page.slr
}

// Values returns the slice of values for the current page or nil if there are no values.
func (page ServiceListResultPage) Values() []Service {
	if page.slr.IsEmpty() {
		return nil
	}
	return *page.slr.Value
}

// Creates a new instance of the ServiceListResultPage type.
func NewServiceListResultPage(cur ServiceListResult, getNextPage func(context.Context, ServiceListResult) (ServiceListResult, error)) ServiceListResultPage {
	return ServiceListResultPage{
		fn:  getNextPage,
		slr: cur,
	}
}

// ServiceProperties properties of the Search service.
type ServiceProperties struct {
	// ReplicaCount - The number of replicas in the Search service. If specified, it must be a value between 1 and 12 inclusive for standard SKUs or between 1 and 3 inclusive for basic SKU.
	ReplicaCount *int32 `json:"replicaCount,omitempty"`
	// PartitionCount - The number of partitions in the Search service; if specified, it can be 1, 2, 3, 4, 6, or 12. Values greater than 1 are only valid for standard SKUs. For 'standard3' services with hostingMode set to 'highDensity', the allowed values are between 1 and 3.
	PartitionCount *int32 `json:"partitionCount,omitempty"`
	// HostingMode - Applicable only for the standard3 SKU. You can set this property to enable up to 3 high density partitions that allow up to 1000 indexes, which is much higher than the maximum indexes allowed for any other SKU. For the standard3 SKU, the value is either 'default' or 'highDensity'. For all other SKUs, this value must be 'default'. Possible values include: 'Default', 'HighDensity'
	HostingMode HostingMode `json:"hostingMode,omitempty"`
	// PublicNetworkAccess - This value can be set to 'enabled' to avoid breaking changes on existing customer resources and templates. If set to 'disabled', traffic over public interface is not allowed, and private endpoint connections would be the exclusive access method. Possible values include: 'Enabled', 'Disabled'
	PublicNetworkAccess PublicNetworkAccess `json:"publicNetworkAccess,omitempty"`
	// Status - READ-ONLY; The status of the Search service. Possible values include: 'running': The Search service is running and no provisioning operations are underway. 'provisioning': The Search service is being provisioned or scaled up or down. 'deleting': The Search service is being deleted. 'degraded': The Search service is degraded. This can occur when the underlying search units are not healthy. The Search service is most likely operational, but performance might be slow and some requests might be dropped. 'disabled': The Search service is disabled. In this state, the service will reject all API requests. 'error': The Search service is in an error state. If your service is in the degraded, disabled, or error states, it means the Azure Cognitive Search team is actively investigating the underlying issue. Dedicated services in these states are still chargeable based on the number of search units provisioned. Possible values include: 'ServiceStatusRunning', 'ServiceStatusProvisioning', 'ServiceStatusDeleting', 'ServiceStatusDegraded', 'ServiceStatusDisabled', 'ServiceStatusError'
	Status ServiceStatus `json:"status,omitempty"`
	// StatusDetails - READ-ONLY; The details of the Search service status.
	StatusDetails *string `json:"statusDetails,omitempty"`
	// ProvisioningState - READ-ONLY; The state of the last provisioning operation performed on the Search service. Provisioning is an intermediate state that occurs while service capacity is being established. After capacity is set up, provisioningState changes to either 'succeeded' or 'failed'. Client applications can poll provisioning status (the recommended polling interval is from 30 seconds to one minute) by using the Get Search Service operation to see when an operation is completed. If you are using the free service, this value tends to come back as 'succeeded' directly in the call to Create Search service. This is because the free service uses capacity that is already set up. Possible values include: 'Succeeded', 'Provisioning', 'Failed'
	ProvisioningState ProvisioningState `json:"provisioningState,omitempty"`
	// NetworkRuleSet - Network specific rules that determine how the Azure Cognitive Search service may be reached.
	NetworkRuleSet *NetworkRuleSet `json:"networkRuleSet,omitempty"`
	// PrivateEndpointConnections - READ-ONLY; The list of private endpoint connections to the Azure Cognitive Search service.
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
}

// MarshalJSON is the custom marshaler for ServiceProperties.
func (sp ServiceProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if sp.ReplicaCount != nil {
		objectMap["replicaCount"] = sp.ReplicaCount
	}
	if sp.PartitionCount != nil {
		objectMap["partitionCount"] = sp.PartitionCount
	}
	if sp.HostingMode != "" {
		objectMap["hostingMode"] = sp.HostingMode
	}
	if sp.PublicNetworkAccess != "" {
		objectMap["publicNetworkAccess"] = sp.PublicNetworkAccess
	}
	if sp.NetworkRuleSet != nil {
		objectMap["networkRuleSet"] = sp.NetworkRuleSet
	}
	return json.Marshal(objectMap)
}

// ServicesCreateOrUpdateFuture an abstraction for monitoring and retrieving the results of a long-running
// operation.
type ServicesCreateOrUpdateFuture struct {
	azure.FutureAPI
	// Result returns the result of the asynchronous operation.
	// If the operation has not completed it will return an error.
	Result func(ServicesClient) (Service, error)
}

// UnmarshalJSON is the custom unmarshaller for CreateFuture.
func (future *ServicesCreateOrUpdateFuture) UnmarshalJSON(body []byte) error {
	var azFuture azure.Future
	if err := json.Unmarshal(body, &azFuture); err != nil {
		return err
	}
	future.FutureAPI = &azFuture
	future.Result = future.result
	return nil
}

// result is the default implementation for ServicesCreateOrUpdateFuture.Result.
func (future *ServicesCreateOrUpdateFuture) result(client ServicesClient) (s Service, err error) {
	var done bool
	done, err = future.DoneWithContext(context.Background(), client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "search.ServicesCreateOrUpdateFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = azure.NewAsyncOpIncompleteError("search.ServicesCreateOrUpdateFuture")
		return
	}
	sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	if s.Response.Response, err = future.GetResult(sender); err == nil && s.Response.Response.StatusCode != http.StatusNoContent {
		s, err = client.CreateOrUpdateResponder(s.Response.Response)
		if err != nil {
			err = autorest.NewErrorWithError(err, "search.ServicesCreateOrUpdateFuture", "Result", s.Response.Response, "Failure responding to request")
		}
	}
	return
}

// Sku defines the SKU of an Azure Cognitive Search Service, which determines price tier and capacity
// limits.
type Sku struct {
	// Name - The SKU of the Search service. Valid values include: 'free': Shared service. 'basic': Dedicated service with up to 3 replicas. 'standard': Dedicated service with up to 12 partitions and 12 replicas. 'standard2': Similar to standard, but with more capacity per search unit. 'standard3': The largest Standard offering with up to 12 partitions and 12 replicas (or up to 3 partitions with more indexes if you also set the hostingMode property to 'highDensity'). 'storage_optimized_l1': Supports 1TB per partition, up to 12 partitions. 'storage_optimized_l2': Supports 2TB per partition, up to 12 partitions.'. Possible values include: 'Free', 'Basic', 'Standard', 'Standard2', 'Standard3', 'StorageOptimizedL1', 'StorageOptimizedL2'
	Name SkuName `json:"name,omitempty"`
}
