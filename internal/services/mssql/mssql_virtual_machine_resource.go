package mssql

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	parseCompute "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlVirtualMachine() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlVirtualMachineCreateUpdate,
		Read:   resourceMsSqlVirtualMachineRead,
		Update: resourceMsSqlVirtualMachineCreateUpdate,
		Delete: resourceMsSqlVirtualMachineDelete,

		CustomizeDiff: pluginsdk.CustomizeDiffShim(resourceMsSqlVirtualMachineCustomDiff),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := sqlvirtualmachines.ParseSqlVirtualMachineID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"virtual_machine_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.VirtualMachineID,
			},

			"sql_license_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sqlvirtualmachines.SqlServerLicenseTypePAYG),
					string(sqlvirtualmachines.SqlServerLicenseTypeAHUB),
					string(sqlvirtualmachines.SqlServerLicenseTypeDR),
				}, false),
			},

			"auto_backup": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"encryption_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"encryption_password": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"manual_schedule": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"full_backup_frequency": {
										Type:             pluginsdk.TypeString,
										Required:         true,
										DiffSuppressFunc: suppress.CaseDifference,
										ValidateFunc: validation.StringInSlice([]string{
											string(sqlvirtualmachines.FullBackupFrequencyTypeDaily),
											string(sqlvirtualmachines.FullBackupFrequencyTypeWeekly),
										}, false),
									},

									"full_backup_start_hour": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 23),
									},

									"full_backup_window_in_hours": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 23),
									},

									"log_backup_frequency_in_minutes": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(5, 60),
									},
								},
							},
						},

						"retention_period_in_days": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 30),
						},

						"storage_blob_endpoint": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},

						"storage_account_access_key": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"system_databases_backup_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"auto_patching": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"day_of_week": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sqlvirtualmachines.DayOfWeekMonday),
								string(sqlvirtualmachines.DayOfWeekTuesday),
								string(sqlvirtualmachines.DayOfWeekWednesday),
								string(sqlvirtualmachines.DayOfWeekThursday),
								string(sqlvirtualmachines.DayOfWeekFriday),
								string(sqlvirtualmachines.DayOfWeekSaturday),
								string(sqlvirtualmachines.DayOfWeekSunday),
							}, false),
						},

						"maintenance_window_duration_in_minutes": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(30, 180),
						},

						"maintenance_window_starting_hour": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 23),
						},
					},
				},
			},

			"key_vault_credential": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							// api will add updated credential name, and return "sqlvmName:name1,sqlvmName:name2"
							DiffSuppressFunc: mssqlVMCredentialNameDiffSuppressFunc,
						},

						"key_vault_url": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							Sensitive:    true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},

						"service_principal_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"service_principal_secret": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"r_services_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"sql_connectivity_port": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1433,
				ValidateFunc: validation.IntBetween(1024, 65535),
			},

			"sql_connectivity_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(sqlvirtualmachines.ConnectivityTypePRIVATE),
				ValidateFunc: validation.StringInSlice([]string{
					string(sqlvirtualmachines.ConnectivityTypeLOCAL),
					string(sqlvirtualmachines.ConnectivityTypePRIVATE),
					string(sqlvirtualmachines.ConnectivityTypePUBLIC),
				}, false),
			},

			"sql_connectivity_update_password": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"sql_connectivity_update_username": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validate.SqlVirtualMachineLoginUserName,
			},

			"storage_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"disk_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sqlvirtualmachines.DiskConfigurationTypeNEW),
								string(sqlvirtualmachines.DiskConfigurationTypeEXTEND),
								string(sqlvirtualmachines.DiskConfigurationTypeADD),
							}, false),
						},
						"storage_workload_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sqlvirtualmachines.SqlWorkloadTypeGENERAL),
								string(sqlvirtualmachines.SqlWorkloadTypeOLTP),
								string(sqlvirtualmachines.SqlWorkloadTypeDW),
							}, false),
						},
						"data_settings":    helper.StorageSettingSchema(),
						"log_settings":     helper.StorageSettingSchema(),
						"temp_db_settings": helper.StorageSettingSchema(),
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceMsSqlVirtualMachineCustomDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	// ForceNew when removing the auto_backup block.
	// See https://github.com/Azure/azure-rest-api-specs/issues/12818#issuecomment-773727756
	old, new := d.GetChange("auto_backup")
	if len(old.([]interface{})) == 1 && len(new.([]interface{})) == 0 {
		return d.ForceNew("auto_backup")
	}

	encryptionEnabled := d.Get("auto_backup.0.encryption_enabled")
	v, ok := d.GetOk("auto_backup.0.encryption_password")

	if encryptionEnabled.(bool) && (!ok || v.(string) == "") {
		return fmt.Errorf("auto_backup: `encryption_password` is required when `encryption_enabled` is true")
	}

	if !encryptionEnabled.(bool) && ok && v.(string) != "" {
		return fmt.Errorf("auto_backup: `encryption_enabled` must be true when `encryption_password` is set")
	}

	return nil
}

func resourceMsSqlVirtualMachineCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualMachinesClient
	vmclient := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmId, err := parseCompute.VirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return err
	}
	id := sqlvirtualmachines.NewSqlVirtualMachineID(vmId.SubscriptionId, vmId.ResourceGroup, vmId.Name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, sqlvirtualmachines.GetOperationOptions{Expand: utils.String("*")})
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_mssql_virtual_machine", id.ID())
		}
	}

	// get location from vm
	respvm, err := vmclient.Get(ctx, id.ResourceGroupName, id.SqlVirtualMachineName, "")
	if err != nil {
		return fmt.Errorf("making Read request on Azure Virtual Machine %s: %+v", id.SqlVirtualMachineName, err)
	}

	if *respvm.Location == "" {
		return fmt.Errorf("location is empty from making Read request on Azure Virtual Machine %s: %+v", id.SqlVirtualMachineName, err)
	}

	connectivityType := sqlvirtualmachines.ConnectivityType(d.Get("sql_connectivity_type").(string))
	sqlManagement := sqlvirtualmachines.SqlManagementModeFull
	sqlServerLicenseType := sqlvirtualmachines.SqlServerLicenseType(d.Get("sql_license_type").(string))

	parameters := sqlvirtualmachines.SqlVirtualMachine{
		Location: *respvm.Location,
		Properties: &sqlvirtualmachines.SqlVirtualMachineProperties{
			AutoBackupSettings:         expandSqlVirtualMachineAutoBackupSettings(d.Get("auto_backup").([]interface{})),
			AutoPatchingSettings:       expandSqlVirtualMachineAutoPatchingSettings(d.Get("auto_patching").([]interface{})),
			KeyVaultCredentialSettings: expandSqlVirtualMachineKeyVaultCredential(d.Get("key_vault_credential").([]interface{})),
			ServerConfigurationsManagementSettings: &sqlvirtualmachines.ServerConfigurationsManagementSettings{
				AdditionalFeaturesServerConfigurations: &sqlvirtualmachines.AdditionalFeaturesServerConfigurations{
					IsRServicesEnabled: utils.Bool(d.Get("r_services_enabled").(bool)),
				},
				SqlConnectivityUpdateSettings: &sqlvirtualmachines.SqlConnectivityUpdateSettings{
					ConnectivityType:      &connectivityType,
					Port:                  utils.Int64(int64(d.Get("sql_connectivity_port").(int))),
					SqlAuthUpdatePassword: utils.String(d.Get("sql_connectivity_update_password").(string)),
					SqlAuthUpdateUserName: utils.String(d.Get("sql_connectivity_update_username").(string)),
				},
			},
			SqlManagement:                &sqlManagement,
			SqlServerLicenseType:         &sqlServerLicenseType,
			StorageConfigurationSettings: expandSqlVirtualMachineStorageConfigurationSettings(d.Get("storage_configuration").([]interface{})),
			VirtualMachineResourceId:     utils.String(d.Get("virtual_machine_id").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// Wait for the auto backup settings to take effect
	// See: https://github.com/Azure/azure-rest-api-specs/issues/12818
	if autoBackup := d.Get("auto_backup"); (d.IsNewResource() && len(autoBackup.([]interface{})) > 0) || (!d.IsNewResource() && d.HasChange("auto_backup")) {
		log.Printf("[DEBUG] Waiting for SQL Virtual Machine %q AutoBackupSettings to take effect", d.Id())
		stateConf := &pluginsdk.StateChangeConf{
			Pending:                   []string{"Retry", "Pending"},
			Target:                    []string{"Updated"},
			Refresh:                   resourceMsSqlVirtualMachineAutoBackupSettingsRefreshFunc(ctx, client, d),
			MinTimeout:                1 * time.Minute,
			ContinuousTargetOccurence: 2,
		}

		if d.IsNewResource() {
			stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
		} else {
			stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for SQL Virtual Machine %q AutoBackupSettings to take effect: %+v", d.Id(), err)
		}
	}

	// Wait for the auto patching settings to take effect
	// See: https://github.com/Azure/azure-rest-api-specs/issues/12818
	if autoPatching := d.Get("auto_patching"); (d.IsNewResource() && len(autoPatching.([]interface{})) > 0) || (!d.IsNewResource() && d.HasChange("auto_patching")) {
		log.Printf("[DEBUG] Waiting for SQL Virtual Machine %q AutoPatchingSettings to take effect", d.Id())
		stateConf := &pluginsdk.StateChangeConf{
			Pending:                   []string{"Retry", "Pending"},
			Target:                    []string{"Updated"},
			Refresh:                   resourceMsSqlVirtualMachineAutoPatchingSettingsRefreshFunc(ctx, client, d),
			MinTimeout:                1 * time.Minute,
			ContinuousTargetOccurence: 2,
		}

		if d.IsNewResource() {
			stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
		} else {
			stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for SQL Virtual Machine %q AutoPatchingSettings to take effect: %+v", d.Id(), err)
		}
	}

	return resourceMsSqlVirtualMachineRead(d, meta)
}

func resourceMsSqlVirtualMachineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, sqlvirtualmachines.GetOperationOptions{Expand: utils.String("*")})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Sql Virtual Machine %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {

			d.Set("virtual_machine_id", props.VirtualMachineResourceId)
			d.Set("sql_license_type", string(*props.SqlServerLicenseType))
			if err := d.Set("auto_backup", flattenSqlVirtualMachineAutoBackup(props.AutoBackupSettings, d)); err != nil {
				return fmt.Errorf("setting `auto_backup`: %+v", err)
			}

			if err := d.Set("auto_patching", flattenSqlVirtualMachineAutoPatching(props.AutoPatchingSettings)); err != nil {
				return fmt.Errorf("setting `auto_patching`: %+v", err)
			}

			if err := d.Set("key_vault_credential", flattenSqlVirtualMachineKeyVaultCredential(props.KeyVaultCredentialSettings, d)); err != nil {
				return fmt.Errorf("setting `key_vault_credential`: %+v", err)
			}

			if mgmtSettings := props.ServerConfigurationsManagementSettings; mgmtSettings != nil {
				if cfgs := mgmtSettings.AdditionalFeaturesServerConfigurations; cfgs != nil {
					d.Set("r_services_enabled", mgmtSettings.AdditionalFeaturesServerConfigurations.IsRServicesEnabled)
				}
				if scus := mgmtSettings.SqlConnectivityUpdateSettings; scus != nil {
					d.Set("sql_connectivity_port", mgmtSettings.SqlConnectivityUpdateSettings.Port)
					d.Set("sql_connectivity_type", mgmtSettings.SqlConnectivityUpdateSettings.ConnectivityType)
				}
			}

			// `storage_configuration.0.storage_workload_type` is in a different spot than the rest of the `storage_configuration`
			// so we'll grab that here and pass it along
			storageWorkloadType := ""
			if props.ServerConfigurationsManagementSettings != nil && props.ServerConfigurationsManagementSettings.SqlWorkloadTypeUpdateSettings != nil && props.ServerConfigurationsManagementSettings.SqlWorkloadTypeUpdateSettings.SqlWorkloadType != nil {
				storageWorkloadType = string(*props.ServerConfigurationsManagementSettings.SqlWorkloadTypeUpdateSettings.SqlWorkloadType)
			}

			if err := d.Set("storage_configuration", flattenSqlVirtualMachineStorageConfigurationSettings(props.StorageConfigurationSettings, storageWorkloadType)); err != nil {
				return fmt.Errorf("setting `storage_configuration`: %+v", err)
			}
			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceMsSqlVirtualMachineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualMachinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func resourceMsSqlVirtualMachineAutoBackupSettingsRefreshFunc(ctx context.Context, client *sqlvirtualmachines.SqlVirtualMachinesClient, d *pluginsdk.ResourceData) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(d.Id())
		if err != nil {
			return nil, "Error", err
		}

		resp, err := client.Get(ctx, *id, sqlvirtualmachines.GetOperationOptions{Expand: utils.String("*")})
		if err != nil {
			return nil, "Retry", err
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				autoBackupSettings := flattenSqlVirtualMachineAutoBackup(props.AutoBackupSettings, d)

				if len(autoBackupSettings) == 0 {
					// auto backup was nil or disabled in the response
					if v, ok := d.GetOk("auto_backup"); !ok || len(v.([]interface{})) == 0 {
						// also disabled in the config
						return resp, "Updated", nil
					}
					return resp, "Pending", nil
				}

				if v, ok := d.GetOk("auto_backup"); !ok || len(v.([]interface{})) == 0 {
					// still waiting for it to be disabled
					return resp, "Pending", nil
				}

				// check each property in the auto_backup block for drift
				for prop, val := range autoBackupSettings[0].(map[string]interface{}) {
					v := d.Get(fmt.Sprintf("auto_backup.0.%s", prop))
					switch prop {
					case "manual_schedule":
						if m := val.([]interface{}); len(m) > 0 {
							if b, ok := d.GetOk("auto_backup.0.manual_schedule"); !ok || len(b.([]interface{})) == 0 {
								// manual schedule disabled in config but still showing in response
								return resp, "Pending", nil
							}
							// check each property in the manual_schedule block for drift
							for prop2, val2 := range m[0].(map[string]interface{}) {
								v2 := d.Get(fmt.Sprintf("auto_backup.0.manual_schedule.0.%s", prop2))
								switch prop2 {
								case "full_backup_frequency":
									if !strings.EqualFold(v2.(string), val2.(string)) {
										return resp, "Pending", nil
									}
								default:
									if v2 != val2 {
										return resp, "Pending", nil
									}
								}
							}
						} else if b, ok := d.GetOk("auto_backup.0.manual_schedule"); ok || len(b.([]interface{})) > 0 {
							// manual schedule set in config but not reflecting in response
							return resp, "Pending", nil
						}
					default:
						if v != val {
							return resp, "Pending", nil
						}
					}
				}

				return resp, "Updated", nil
			}
		}

		return resp, "Retry", nil
	}
}

func expandSqlVirtualMachineAutoBackupSettings(input []interface{}) *sqlvirtualmachines.AutoBackupSettings {
	ret := sqlvirtualmachines.AutoBackupSettings{
		Enable: utils.Bool(false),
	}

	if len(input) > 0 {
		config := input[0].(map[string]interface{})
		ret.Enable = utils.Bool(true)

		if v, ok := config["retention_period_in_days"]; ok {
			ret.RetentionPeriod = utils.Int64(int64(v.(int)))
		}
		if v, ok := config["storage_blob_endpoint"]; ok {
			ret.StorageAccountUrl = utils.String(v.(string))
		}
		if v, ok := config["storage_account_access_key"]; ok {
			ret.StorageAccessKey = utils.String(v.(string))
		}

		v, ok := config["encryption_enabled"]
		enableEncryption := ok && v.(bool)
		ret.EnableEncryption = utils.Bool(enableEncryption)
		if v, ok := config["encryption_password"]; enableEncryption && ok {
			ret.Password = utils.String(v.(string))
		}

		if v, ok := config["system_databases_backup_enabled"]; ok {
			ret.BackupSystemDbs = utils.Bool(v.(bool))
		}

		backupScheduleTypeAutomated := sqlvirtualmachines.BackupScheduleTypeAutomated
		ret.BackupScheduleType = &backupScheduleTypeAutomated
		if v, ok := config["manual_schedule"]; ok && len(v.([]interface{})) > 0 {
			manualSchedule := v.([]interface{})[0].(map[string]interface{})
			backupScheduleTypeManual := sqlvirtualmachines.BackupScheduleTypeManual
			ret.BackupScheduleType = &backupScheduleTypeManual

			if v, ok := manualSchedule["full_backup_frequency"]; ok {
				fullBackupFrequencyType := sqlvirtualmachines.FullBackupFrequencyType(v.(string))
				ret.FullBackupFrequency = &fullBackupFrequencyType
			}
			if v, ok := manualSchedule["full_backup_start_hour"]; ok {
				ret.FullBackupStartTime = utils.Int64(int64(v.(int)))
			}
			if v, ok := manualSchedule["full_backup_window_in_hours"]; ok {
				ret.FullBackupWindowHours = utils.Int64(int64(v.(int)))
			}
			if v, ok := manualSchedule["log_backup_frequency_in_minutes"]; ok {
				ret.LogBackupFrequency = utils.Int64(int64(v.(int)))
			}
		}
	}

	return &ret
}

func flattenSqlVirtualMachineAutoBackup(autoBackup *sqlvirtualmachines.AutoBackupSettings, d *pluginsdk.ResourceData) []interface{} {
	if autoBackup == nil || autoBackup.Enable == nil || !*autoBackup.Enable {
		return []interface{}{}
	}

	manualSchedule := make([]interface{}, 0)
	if autoBackup.BackupScheduleType != nil && strings.EqualFold(string(*autoBackup.BackupScheduleType), string(sqlvirtualmachines.BackupScheduleTypeManual)) {
		var fullBackupStartHour int
		if autoBackup.FullBackupStartTime != nil {
			fullBackupStartHour = int(*autoBackup.FullBackupStartTime)
		}

		var fullBackupWindowHours int
		if autoBackup.FullBackupWindowHours != nil {
			fullBackupWindowHours = int(*autoBackup.FullBackupWindowHours)
		}

		var logBackupFrequency int
		if autoBackup.LogBackupFrequency != nil {
			logBackupFrequency = int(*autoBackup.LogBackupFrequency)
			// API returns 60 minutes as zero
			if logBackupFrequency == 0 {
				logBackupFrequency = 60
			}
		}

		var fullBackupFrequency string
		if autoBackup.FullBackupFrequency != nil {
			fullBackupFrequency = string(*autoBackup.FullBackupFrequency)
		}

		manualSchedule = []interface{}{
			map[string]interface{}{
				"full_backup_frequency":           fullBackupFrequency,
				"full_backup_start_hour":          fullBackupStartHour,
				"full_backup_window_in_hours":     fullBackupWindowHours,
				"log_backup_frequency_in_minutes": logBackupFrequency,
			},
		}
	}

	var retentionPeriod int
	if autoBackup.RetentionPeriod != nil {
		retentionPeriod = int(*autoBackup.RetentionPeriod)
	}

	// Password, StorageAccessKey, StorageAccountURL are not returned, so we try to copy them
	// from existing config as a best effort.
	encryptionPassword := d.Get("auto_backup.0.encryption_password").(string)
	storageKey := d.Get("auto_backup.0.storage_account_access_key").(string)
	blobEndpoint := d.Get("auto_backup.0.storage_blob_endpoint").(string)

	return []interface{}{
		map[string]interface{}{
			"encryption_enabled":              autoBackup.EnableEncryption != nil && *autoBackup.EnableEncryption,
			"encryption_password":             encryptionPassword,
			"manual_schedule":                 manualSchedule,
			"retention_period_in_days":        retentionPeriod,
			"storage_account_access_key":      storageKey,
			"storage_blob_endpoint":           blobEndpoint,
			"system_databases_backup_enabled": autoBackup.BackupSystemDbs != nil && *autoBackup.BackupSystemDbs,
		},
	}
}

func resourceMsSqlVirtualMachineAutoPatchingSettingsRefreshFunc(ctx context.Context, client *sqlvirtualmachines.SqlVirtualMachinesClient, d *pluginsdk.ResourceData) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(d.Id())
		if err != nil {
			return nil, "Error", err
		}

		resp, err := client.Get(ctx, *id, sqlvirtualmachines.GetOperationOptions{Expand: utils.String("*")})
		if err != nil {
			return nil, "Retry", err
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				autoPatchingSettings := flattenSqlVirtualMachineAutoPatching(props.AutoPatchingSettings)

				if len(autoPatchingSettings) == 0 {
					if v, ok := d.GetOk("auto_patching"); !ok || len(v.([]interface{})) == 0 {
						return resp, "Updated", nil
					}
					return resp, "Pending", nil
				}

				if v, ok := d.GetOk("auto_patching"); !ok || len(v.([]interface{})) == 0 {
					return resp, "Pending", nil
				}

				for prop, val := range autoPatchingSettings[0].(map[string]interface{}) {
					v := d.Get(fmt.Sprintf("auto_patching.0.%s", prop))
					if v != val {
						return resp, "Pending", nil
					}
				}

				return resp, "Updated", nil
			}
		}

		return resp, "Retry", nil
	}
}

func expandSqlVirtualMachineAutoPatchingSettings(input []interface{}) *sqlvirtualmachines.AutoPatchingSettings {
	if len(input) == 0 {
		return nil
	}
	autoPatchingSetting := input[0].(map[string]interface{})

	dayOfWeek := sqlvirtualmachines.DayOfWeek(autoPatchingSetting["day_of_week"].(string))

	return &sqlvirtualmachines.AutoPatchingSettings{
		Enable:                        utils.Bool(true),
		MaintenanceWindowDuration:     utils.Int64(int64(autoPatchingSetting["maintenance_window_duration_in_minutes"].(int))),
		MaintenanceWindowStartingHour: utils.Int64(int64(autoPatchingSetting["maintenance_window_starting_hour"].(int))),
		DayOfWeek:                     &dayOfWeek,
	}
}

func flattenSqlVirtualMachineAutoPatching(autoPatching *sqlvirtualmachines.AutoPatchingSettings) []interface{} {
	if autoPatching == nil || autoPatching.Enable == nil || !*autoPatching.Enable {
		return []interface{}{}
	}

	var startHour int
	if autoPatching.MaintenanceWindowStartingHour != nil {
		startHour = int(*autoPatching.MaintenanceWindowStartingHour)
	}

	var duration int
	if autoPatching.MaintenanceWindowDuration != nil {
		duration = int(*autoPatching.MaintenanceWindowDuration)
	}

	var dayOfWeek string
	if autoPatching.DayOfWeek != nil {
		dayOfWeek = string(*autoPatching.DayOfWeek)
	}

	return []interface{}{
		map[string]interface{}{
			"day_of_week":                            dayOfWeek,
			"maintenance_window_starting_hour":       startHour,
			"maintenance_window_duration_in_minutes": duration,
		},
	}
}

func expandSqlVirtualMachineKeyVaultCredential(input []interface{}) *sqlvirtualmachines.KeyVaultCredentialSettings {
	if len(input) == 0 {
		return nil
	}
	keyVaultCredentialSetting := input[0].(map[string]interface{})

	return &sqlvirtualmachines.KeyVaultCredentialSettings{
		Enable:                 utils.Bool(true),
		CredentialName:         utils.String(keyVaultCredentialSetting["name"].(string)),
		AzureKeyVaultUrl:       utils.String(keyVaultCredentialSetting["key_vault_url"].(string)),
		ServicePrincipalName:   utils.String(keyVaultCredentialSetting["service_principal_name"].(string)),
		ServicePrincipalSecret: utils.String(keyVaultCredentialSetting["service_principal_secret"].(string)),
	}
}

func flattenSqlVirtualMachineKeyVaultCredential(keyVault *sqlvirtualmachines.KeyVaultCredentialSettings, d *pluginsdk.ResourceData) []interface{} {
	if keyVault == nil || keyVault.Enable == nil || !*keyVault.Enable {
		return []interface{}{}
	}

	name := ""
	if keyVault.CredentialName != nil {
		name = *keyVault.CredentialName
	}

	keyVaultUrl := ""
	if v, ok := d.GetOk("key_vault_credential.0.key_vault_url"); ok {
		keyVaultUrl = v.(string)
	}

	servicePrincipalName := ""
	if v, ok := d.GetOk("key_vault_credential.0.service_principal_name"); ok {
		servicePrincipalName = v.(string)
	}

	servicePrincipalSecret := ""
	if v, ok := d.GetOk("key_vault_credential.0.service_principal_secret"); ok {
		servicePrincipalSecret = v.(string)
	}

	return []interface{}{
		map[string]interface{}{
			"name":                     name,
			"key_vault_url":            keyVaultUrl,
			"service_principal_name":   servicePrincipalName,
			"service_principal_secret": servicePrincipalSecret,
		},
	}
}

func mssqlVMCredentialNameDiffSuppressFunc(_, old, new string, _ *pluginsdk.ResourceData) bool {
	oldNamelist := strings.Split(old, ",")
	for _, n := range oldNamelist {
		cur := strings.Split(n, ":")
		if len(cur) > 1 && cur[1] == new {
			return true
		}
	}
	return false
}

func expandSqlVirtualMachineStorageConfigurationSettings(input []interface{}) *sqlvirtualmachines.StorageConfigurationSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	storageSettings := input[0].(map[string]interface{})

	diskConfigurationType := sqlvirtualmachines.DiskConfigurationType(storageSettings["disk_type"].(string))
	storageWorkloadType := sqlvirtualmachines.StorageWorkloadType(storageSettings["storage_workload_type"].(string))

	return &sqlvirtualmachines.StorageConfigurationSettings{
		DiskConfigurationType: &diskConfigurationType,
		StorageWorkloadType:   &storageWorkloadType,
		SqlDataSettings:       expandSqlVirtualMachineDataStorageSettings(storageSettings["data_settings"].([]interface{})),
		SqlLogSettings:        expandSqlVirtualMachineDataStorageSettings(storageSettings["log_settings"].([]interface{})),
		SqlTempDbSettings:     expandSqlVirtualMachineTempDbSettings(storageSettings["temp_db_settings"].([]interface{})),
	}
}

func flattenSqlVirtualMachineStorageConfigurationSettings(input *sqlvirtualmachines.StorageConfigurationSettings, storageWorkloadType string) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var diskType string
	if input.DiskConfigurationType != nil {
		diskType = string(*input.DiskConfigurationType)
	}

	output := map[string]interface{}{
		"storage_workload_type": storageWorkloadType,
		"disk_type":             diskType,
		"data_settings":         flattenSqlVirtualMachineStorageSettings(input.SqlDataSettings),
		"log_settings":          flattenSqlVirtualMachineStorageSettings(input.SqlLogSettings),
		"temp_db_settings":      flattenSqlVirtualMachineTempDbSettings(input.SqlTempDbSettings),
	}

	if output["storage_workload_type"].(string) == "" && output["disk_type"] == "" &&
		len(output["data_settings"].([]interface{})) == 0 &&
		len(output["log_settings"].([]interface{})) == 0 &&
		len(output["temp_db_settings"].([]interface{})) == 0 {
		return []interface{}{}
	}

	return []interface{}{output}
}

func expandSqlVirtualMachineDataStorageSettings(input []interface{}) *sqlvirtualmachines.SQLStorageSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	dataStorageSettings := input[0].(map[string]interface{})

	return &sqlvirtualmachines.SQLStorageSettings{
		Luns:            expandSqlVirtualMachineStorageSettingsLuns(dataStorageSettings["luns"].([]interface{})),
		DefaultFilePath: utils.String(dataStorageSettings["default_file_path"].(string)),
	}
}

func expandSqlVirtualMachineStorageSettingsLuns(input []interface{}) *[]int64 {
	expandedLuns := make([]int64, 0)
	for i := range input {
		if input[i] != nil {
			expandedLuns = append(expandedLuns, int64(input[i].(int)))
		}
	}

	return &expandedLuns
}

func flattenSqlVirtualMachineStorageSettings(input *sqlvirtualmachines.SQLStorageSettings) []interface{} {
	if input == nil || input.Luns == nil {
		return []interface{}{}
	}
	attrs := make(map[string]interface{})

	if input.Luns != nil {
		attrs["luns"] = *input.Luns
	}

	if input.DefaultFilePath != nil {
		attrs["default_file_path"] = *input.DefaultFilePath
	}

	return []interface{}{attrs}
}

func expandSqlVirtualMachineTempDbSettings(input []interface{}) *sqlvirtualmachines.SQLTempDbSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	dataStorageSettings := input[0].(map[string]interface{})

	return &sqlvirtualmachines.SQLTempDbSettings{
		Luns:            expandSqlVirtualMachineStorageSettingsLuns(dataStorageSettings["luns"].([]interface{})),
		DefaultFilePath: utils.String(dataStorageSettings["default_file_path"].(string)),
	}
}

func flattenSqlVirtualMachineTempDbSettings(input *sqlvirtualmachines.SQLTempDbSettings) []interface{} {
	if input == nil || input.Luns == nil {
		return []interface{}{}
	}
	attrs := make(map[string]interface{})

	if input.Luns != nil {
		attrs["luns"] = *input.Luns
	}

	if input.DefaultFilePath != nil {
		attrs["default_file_path"] = *input.DefaultFilePath
	}

	return []interface{}{attrs}
}
