// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	Factories *factories.FactoriesClient

	// TODO: convert to using hashicorp/go-azure-sdk
	DataFlowClient                *datafactory.DataFlowsClient
	DatasetClient                 *datafactory.DatasetsClient
	IntegrationRuntimesClient     *datafactory.IntegrationRuntimesClient
	LinkedServiceClient           *datafactory.LinkedServicesClient
	ManagedPrivateEndpointsClient *datafactory.ManagedPrivateEndpointsClient
	ManagedVirtualNetworksClient  *datafactory.ManagedVirtualNetworksClient
	PipelinesClient               *datafactory.PipelinesClient
	TriggersClient                *datafactory.TriggersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	dataFlowClient := datafactory.NewDataFlowsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&dataFlowClient.Client, o.ResourceManagerAuthorizer)

	DatasetClient := datafactory.NewDatasetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatasetClient.Client, o.ResourceManagerAuthorizer)

	factoriesClient, err := factories.NewFactoriesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Factories client: %+v", err)
	}
	o.Configure(factoriesClient.Client, o.Authorizers.ResourceManager)

	IntegrationRuntimesClient := datafactory.NewIntegrationRuntimesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&IntegrationRuntimesClient.Client, o.ResourceManagerAuthorizer)

	LinkedServiceClient := datafactory.NewLinkedServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LinkedServiceClient.Client, o.ResourceManagerAuthorizer)

	ManagedPrivateEndpointsClient := datafactory.NewManagedPrivateEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ManagedPrivateEndpointsClient.Client, o.ResourceManagerAuthorizer)

	ManagedVirtualNetworksClient := datafactory.NewManagedVirtualNetworksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ManagedVirtualNetworksClient.Client, o.ResourceManagerAuthorizer)

	PipelinesClient := datafactory.NewPipelinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PipelinesClient.Client, o.ResourceManagerAuthorizer)

	TriggersClient := datafactory.NewTriggersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&TriggersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DataFlowClient:                &dataFlowClient,
		DatasetClient:                 &DatasetClient,
		Factories:                     factoriesClient,
		IntegrationRuntimesClient:     &IntegrationRuntimesClient,
		LinkedServiceClient:           &LinkedServiceClient,
		ManagedPrivateEndpointsClient: &ManagedPrivateEndpointsClient,
		ManagedVirtualNetworksClient:  &ManagedVirtualNetworksClient,
		PipelinesClient:               &PipelinesClient,
		TriggersClient:                &TriggersClient,
	}, nil
}
