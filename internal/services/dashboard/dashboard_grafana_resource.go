// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dashboard

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2023-09-01/grafanaresource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DashboardGrafanaModel struct {
	Name                              string                                            `tfschema:"name"`
	ResourceGroupName                 string                                            `tfschema:"resource_group_name"`
	ApiKeyEnabled                     bool                                              `tfschema:"api_key_enabled"`
	AutoGeneratedDomainNameLabelScope grafanaresource.AutoGeneratedDomainNameLabelScope `tfschema:"auto_generated_domain_name_label_scope"`
	SMTP                              []SMTPConfigurationModel                          `tfschema:"smtp"`
	DeterministicOutboundIPEnabled    bool                                              `tfschema:"deterministic_outbound_ip_enabled"`
	AzureMonitorWorkspaceIntegrations []AzureMonitorWorkspaceIntegrationModel           `tfschema:"azure_monitor_workspace_integrations"`
	Location                          string                                            `tfschema:"location"`
	PublicNetworkAccessEnabled        bool                                              `tfschema:"public_network_access_enabled"`
	Sku                               string                                            `tfschema:"sku"`
	Tags                              map[string]string                                 `tfschema:"tags"`
	ZoneRedundancyEnabled             bool                                              `tfschema:"zone_redundancy_enabled"`
	Endpoint                          string                                            `tfschema:"endpoint"`
	GrafanaVersion                    string                                            `tfschema:"grafana_version"`
	GrafanaMajorVersion               string                                            `tfschema:"grafana_major_version"`
	OutboundIPs                       []string                                          `tfschema:"outbound_ip"`
}

type AzureMonitorWorkspaceIntegrationModel struct {
	ResourceId string `tfschema:"resource_id"`
}

type SMTPConfigurationModel struct {
	SMTPEnabled    bool   `tfschema:"enabled"`
	FromAddress    string `tfschema:"from_address"`
	FromName       string `tfschema:"from_name"`
	Host           string `tfschema:"host"`
	Password       string `tfschema:"password"`
	SkipVerify     bool   `tfschema:"verification_skip_enabled"`
	User           string `tfschema:"user"`
	StartTLSPolicy string `tfschema:"start_tls_policy"`
}

type DashboardGrafanaResource struct{}

var _ sdk.ResourceWithUpdate = DashboardGrafanaResource{}

func (r DashboardGrafanaResource) ResourceType() string {
	return "azurerm_dashboard_grafana"
}

func (r DashboardGrafanaResource) ModelObject() interface{} {
	return &DashboardGrafanaModel{}
}

func (r DashboardGrafanaResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return grafanaresource.ValidateGrafanaID
}

func (r DashboardGrafanaResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z][-a-zA-Z\d]{0,21}[a-zA-Z\d]$`),
				`The name length must be from 2 to 23 characters. The name can only contain letters, numbers and dashes, and it must begin with a letter and end with a letter or digit.`,
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"api_key_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"auto_generated_domain_name_label_scope": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(grafanaresource.AutoGeneratedDomainNameLabelScopeTenantReuse),
			ValidateFunc: validation.StringInSlice([]string{
				string(grafanaresource.AutoGeneratedDomainNameLabelScopeTenantReuse),
			}, false),
		},

		"smtp": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"from_address": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"from_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "Azure Managed Grafana Notification",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"host": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"password": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"verification_skip_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"user": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"start_tls_policy": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(grafanaresource.StartTLSPolicyOpportunisticStartTLS),
							string(grafanaresource.StartTLSPolicyMandatoryStartTLS),
							string(grafanaresource.StartTLSPolicyNoStartTLS),
						}, false),
					},
				},
			},
		},

		"deterministic_outbound_ip_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"azure_monitor_workspace_integrations": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"resource_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptionalForceNew(),

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"grafana_major_version": {
			Type: pluginsdk.TypeString,
			// TODO: make this field Required (with no default) in 4.0
			Optional: true,
			ForceNew: true,
			Default:  "9",
			ValidateFunc: validation.StringInSlice([]string{
				"9", "10",
			}, false),
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "Standard",
			ValidateFunc: validation.StringInSlice([]string{
				"Standard",
				"Essential",
			}, false),
		},

		"tags": commonschema.Tags(),

		"zone_redundancy_enabled": {
			Type:     pluginsdk.TypeBool,
			ForceNew: true,
			Optional: true,
			Default:  false,
		},
	}
}

func (r DashboardGrafanaResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"grafana_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"outbound_ip": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r DashboardGrafanaResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DashboardGrafanaModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Dashboard.GrafanaResourceClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := grafanaresource.NewGrafanaID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.GrafanaGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identityValue := expandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))

			apiKey := grafanaresource.ApiKeyDisabled
			if model.ApiKeyEnabled {
				apiKey = grafanaresource.ApiKeyEnabled
			}

			deterministicOutboundIP := grafanaresource.DeterministicOutboundIPDisabled
			if model.DeterministicOutboundIPEnabled {
				deterministicOutboundIP = grafanaresource.DeterministicOutboundIPEnabled
			}

			publicNetworkAccess := grafanaresource.PublicNetworkAccessDisabled
			if model.PublicNetworkAccessEnabled {
				publicNetworkAccess = grafanaresource.PublicNetworkAccessEnabled
			}

			zoneRedundancy := grafanaresource.ZoneRedundancyDisabled
			if model.ZoneRedundancyEnabled {
				zoneRedundancy = grafanaresource.ZoneRedundancyEnabled
			}

			properties := &grafanaresource.ManagedGrafana{
				Identity: identityValue,
				Location: utils.String(location.Normalize(model.Location)),
				Properties: &grafanaresource.ManagedGrafanaProperties{
					ApiKey:                            &apiKey,
					AutoGeneratedDomainNameLabelScope: &model.AutoGeneratedDomainNameLabelScope,
					GrafanaConfigurations:             expandSMTPConfigurationModel(model.SMTP),
					DeterministicOutboundIP:           &deterministicOutboundIP,
					GrafanaIntegrations:               expandGrafanaIntegrationsModel(model.AzureMonitorWorkspaceIntegrations),
					GrafanaMajorVersion:               &model.GrafanaMajorVersion,
					PublicNetworkAccess:               &publicNetworkAccess,
					ZoneRedundancy:                    &zoneRedundancy,
				},
				Sku: &grafanaresource.ResourceSku{
					Name: model.Sku,
				},
				Tags: &model.Tags,
			}

			if err := client.GrafanaCreateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DashboardGrafanaResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dashboard.GrafanaResourceClient

			id, err := grafanaresource.ParseGrafanaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DashboardGrafanaModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.GrafanaGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("api_key_enabled") {
				apiKey := grafanaresource.ApiKeyDisabled
				if model.ApiKeyEnabled {
					apiKey = grafanaresource.ApiKeyEnabled
				}

				properties.Properties.ApiKey = &apiKey
			}

			if metadata.ResourceData.HasChange("auto_generated_domain_name_label_scope") {
				properties.Properties.AutoGeneratedDomainNameLabelScope = &model.AutoGeneratedDomainNameLabelScope
			}

			if metadata.ResourceData.HasChange("deterministic_outbound_ip_enabled") {
				deterministicOutboundIP := grafanaresource.DeterministicOutboundIPDisabled
				if model.DeterministicOutboundIPEnabled {
					deterministicOutboundIP = grafanaresource.DeterministicOutboundIPEnabled
				}

				properties.Properties.DeterministicOutboundIP = &deterministicOutboundIP
			}

			if metadata.ResourceData.HasChange("azure_monitor_workspace_integrations") {
				properties.Properties.GrafanaIntegrations = expandGrafanaIntegrationsModel(model.AzureMonitorWorkspaceIntegrations)
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				publicNetworkAccess := grafanaresource.PublicNetworkAccessDisabled
				if model.PublicNetworkAccessEnabled {
					publicNetworkAccess = grafanaresource.PublicNetworkAccessEnabled
				}

				properties.Properties.PublicNetworkAccess = &publicNetworkAccess
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.GrafanaCreateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DashboardGrafanaResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dashboard.GrafanaResourceClient

			id, err := grafanaresource.ParseGrafanaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GrafanaGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := DashboardGrafanaModel{
				Name:              id.GrafanaName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.NormalizeNilable(model.Location),
			}

			identityValue := flattenLegacySystemAndUserAssignedMap(model.Identity)

			if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if properties := model.Properties; properties != nil {
				if properties.ApiKey != nil {
					if *properties.ApiKey == grafanaresource.ApiKeyEnabled {
						state.ApiKeyEnabled = true
					} else {
						state.ApiKeyEnabled = false
					}
				}

				if properties.AutoGeneratedDomainNameLabelScope != nil {
					state.AutoGeneratedDomainNameLabelScope = *properties.AutoGeneratedDomainNameLabelScope
				}

				if properties.GrafanaConfigurations != nil {
					state.SMTP = flattenSMTPConfigurationModel(properties.GrafanaConfigurations.Smtp, metadata.ResourceData)
				}

				if properties.DeterministicOutboundIP != nil {
					if *properties.DeterministicOutboundIP == grafanaresource.DeterministicOutboundIPEnabled {
						state.DeterministicOutboundIPEnabled = true
					} else {
						state.DeterministicOutboundIPEnabled = false
					}
				}

				if properties.Endpoint != nil {
					state.Endpoint = *properties.Endpoint
				}

				if properties.GrafanaIntegrations != nil {
					state.AzureMonitorWorkspaceIntegrations = flattenAzureMonitorWorkspaceIntegrationModelArray(properties.GrafanaIntegrations.AzureMonitorWorkspaceIntegrations)
				}

				if properties.GrafanaVersion != nil {
					state.GrafanaVersion = *properties.GrafanaVersion
				}

				if properties.GrafanaMajorVersion != nil {
					state.GrafanaMajorVersion = *properties.GrafanaMajorVersion
				}

				if properties.OutboundIPs != nil {
					state.OutboundIPs = *properties.OutboundIPs
				}

				if properties.PublicNetworkAccess != nil {
					if *properties.PublicNetworkAccess == grafanaresource.PublicNetworkAccessEnabled {
						state.PublicNetworkAccessEnabled = true
					} else {
						state.PublicNetworkAccessEnabled = false
					}
				}

				if properties.ZoneRedundancy != nil {
					if *properties.ZoneRedundancy == grafanaresource.ZoneRedundancyEnabled {
						state.ZoneRedundancyEnabled = true
					} else {
						state.ZoneRedundancyEnabled = false
					}
				}
			}

			if model.Sku != nil {
				state.Sku = model.Sku.Name
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DashboardGrafanaResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dashboard.GrafanaResourceClient

			id, err := grafanaresource.ParseGrafanaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.GrafanaDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandSMTPConfigurationModel(input []SMTPConfigurationModel) *grafanaresource.GrafanaConfigurations {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	smtp := grafanaresource.Smtp{
		Enabled:        pointer.To(v.SMTPEnabled),
		FromAddress:    pointer.To(v.FromAddress),
		FromName:       pointer.To(v.FromName),
		Host:           pointer.To(v.Host),
		Password:       pointer.To(v.Password),
		SkipVerify:     pointer.To(v.SkipVerify),
		User:           pointer.To(v.User),
		StartTLSPolicy: (*grafanaresource.StartTLSPolicy)(pointer.To(v.StartTLSPolicy)),
	}

	return pointer.To(
		grafanaresource.GrafanaConfigurations{
			Smtp: pointer.To(smtp),
		})
}

func expandGrafanaIntegrationsModel(inputList []AzureMonitorWorkspaceIntegrationModel) *grafanaresource.GrafanaIntegrations {
	if len(inputList) == 0 {
		return nil
	}

	return &grafanaresource.GrafanaIntegrations{
		AzureMonitorWorkspaceIntegrations: expandAzureMonitorWorkspaceIntegrationModelArray(inputList),
	}
}

func expandAzureMonitorWorkspaceIntegrationModelArray(inputList []AzureMonitorWorkspaceIntegrationModel) *[]grafanaresource.AzureMonitorWorkspaceIntegration {
	var outputList []grafanaresource.AzureMonitorWorkspaceIntegration
	for _, v := range inputList {
		input := v
		output := grafanaresource.AzureMonitorWorkspaceIntegration{
			AzureMonitorWorkspaceResourceId: &input.ResourceId,
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func expandLegacySystemAndUserAssignedMap(input []interface{}) *identity.LegacySystemAndUserAssignedMap {
	identityValue, err := identity.ExpandSystemOrUserAssignedMap(input)
	if err != nil {
		return nil
	}

	return &identity.LegacySystemAndUserAssignedMap{
		Type:        identityValue.Type,
		PrincipalId: identityValue.PrincipalId,
		TenantId:    identityValue.TenantId,
		IdentityIds: identityValue.IdentityIds,
	}
}

func flattenSMTPConfigurationModel(input *grafanaresource.Smtp, data *schema.ResourceData) []SMTPConfigurationModel {
	var outputList []SMTPConfigurationModel
	if input == nil || !pointer.From(input.Enabled) {
		return outputList
	}

	var output SMTPConfigurationModel

	if input.Enabled != nil {
		output.SMTPEnabled = pointer.From(input.Enabled)
	}

	if input.FromAddress != nil {
		output.FromAddress = pointer.From(input.FromAddress)
	}

	if input.FromName != nil {
		output.FromName = pointer.From(input.FromName)
	}

	if input.Host != nil {
		output.Host = pointer.From(input.Host)
	}

	if input.SkipVerify != nil {
		output.SkipVerify = pointer.From(input.SkipVerify)
	}

	if input.User != nil {
		output.User = pointer.From(input.User)
	}

	if input.StartTLSPolicy != nil {
		output.StartTLSPolicy = string(pointer.From(input.StartTLSPolicy))
	}

	output.Password = data.Get("smtp.0.password").(string)

	outputList = append(outputList, output)

	return outputList
}

func flattenAzureMonitorWorkspaceIntegrationModelArray(inputList *[]grafanaresource.AzureMonitorWorkspaceIntegration) []AzureMonitorWorkspaceIntegrationModel {
	var outputList []AzureMonitorWorkspaceIntegrationModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := AzureMonitorWorkspaceIntegrationModel{}

		if input.AzureMonitorWorkspaceResourceId != nil {
			output.ResourceId = *input.AzureMonitorWorkspaceResourceId
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenLegacySystemAndUserAssignedMap(input *identity.LegacySystemAndUserAssignedMap) *[]interface{} {
	if input == nil {
		return &[]interface{}{}
	}

	identityValue := &identity.SystemOrUserAssignedMap{
		Type:        input.Type,
		PrincipalId: input.PrincipalId,
		TenantId:    input.TenantId,
		IdentityIds: input.IdentityIds,
	}

	output, err := identity.FlattenSystemOrUserAssignedMap(identityValue)
	if err != nil {
		return &[]interface{}{}
	}
	return output
}
