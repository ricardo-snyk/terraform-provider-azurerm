package network

import (
	"fmt"
	"time"

	networkNew "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-11-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type NetworkWatcherFlowLogAccountID struct {
	azure.ResourceID
	NetworkWatcherName string
	FlowLogName        string
}

func ParseNetworkWatcherFlowLogID(id string) (*NetworkWatcherFlowLogAccountID, error) {

	flowLogID, err := azure.ParseAzureResourceID(id)
	if err != nil {
		return nil, err
	}
	watcherName, ok := flowLogID.Path["networkWatchers"]
	if !ok {
		return nil, fmt.Errorf("Error: Watcher name not found in Flow Log ID: %s", id)
	}
	flowLogName, ok := flowLogID.Path["flowLogs"]
	if !ok {
		return nil, fmt.Errorf("Error: Flow Log name not found in Flow Log ID: %s", id)
	}

	return &NetworkWatcherFlowLogAccountID{
		ResourceID:         *flowLogID,
		NetworkWatcherName: watcherName,
		FlowLogName:        flowLogName,
	}, nil
}

func resourceArmNetworkWatcherFlowLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkWatcherFlowLogCreateUpdate,
		Read:   resourceArmNetworkWatcherFlowLogRead,
		Update: resourceArmNetworkWatcherFlowLogCreateUpdate,
		Delete: resourceArmNetworkWatcherFlowLogDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"network_watcher_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"network_security_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"retention_policy": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:             schema.TypeBool,
							Required:         true,
							DiffSuppressFunc: azureRMSuppressFlowLogRetentionPolicyEnabledDiff,
						},

						"days": {
							Type:             schema.TypeInt,
							Required:         true,
							DiffSuppressFunc: azureRMSuppressFlowLogRetentionPolicyDaysDiff,
						},
					},
				},
			},

			"traffic_analytics": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"workspace_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},

						"workspace_region": {
							Type:             schema.TypeString,
							Required:         true,
							StateFunc:        azure.NormalizeLocation,
							DiffSuppressFunc: azure.SuppressLocationDiff,
						},

						"workspace_resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
						},
					},
				},
			},

			"version": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 2),
			},
		},
	}
}

func azureRMSuppressFlowLogRetentionPolicyEnabledDiff(k, old, new string, d *schema.ResourceData) bool {
	// Ignore if flow log is disabled as the returned flow log configuration
	// returns default value `false` which may differ from config
	return old != "" && !d.Get("enabled").(bool)
}

func azureRMSuppressFlowLogRetentionPolicyDaysDiff(k, old, new string, d *schema.ResourceData) bool {
	// Ignore if flow log is disabled as the returned flow log configuration
	// returns default value `0` which may differ from config
	return old != "" && !d.Get("enabled").(bool)
}

func resourceArmNetworkWatcherFlowLogCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	flowLogsClient := meta.(*clients.Client).Network.FlowLogsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	networkWatcherName := d.Get("network_watcher_name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	networkSecurityGroupID := d.Get("network_security_group_id").(string)
	storageAccountID := d.Get("storage_account_id").(string)
	enabled := d.Get("enabled").(bool)

	parameters := networkNew.FlowLog{
		FlowLogPropertiesFormat: &networkNew.FlowLogPropertiesFormat{
			StorageID:        &storageAccountID,
			Enabled:          &enabled,
			TargetResourceID: &networkSecurityGroupID,
			RetentionPolicy:  expandAzureRmNetworkWatcherFlowLogRetentionPolicy(d),
		},
	}

	if _, ok := d.GetOk("traffic_analytics"); ok {
		parameters.FlowAnalyticsConfiguration = expandAzureRmNetworkWatcherFlowLogTrafficAnalytics(d)
	}

	if version, ok := d.GetOk("version"); ok {
		format := &networkNew.FlowLogFormatParameters{
			Version: utils.Int32(int32(version.(int))),
		}

		parameters.FlowLogPropertiesFormat.Format = format
	}

	future, err := flowLogsClient.CreateOrUpdate(ctx, resourceGroupName, networkWatcherName, networkSecurityGroupID, parameters)
	if err != nil {
		return fmt.Errorf("Error creating or updating Flow Log on Watcher %q: %+v", networkWatcherName, err)
	}

	if err = future.WaitForCompletionRef(ctx, flowLogsClient.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of setting Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
	}

	resp, err := flowLogsClient.Get(ctx, resourceGroupName, networkWatcherName, networkSecurityGroupID)
	if err != nil {
		return fmt.Errorf("Cannot read Network Watcher %q (Resource Group %q) err: %+v", networkWatcherName, resourceGroupName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Flow Log %q is nil (Resource Group %q)", networkWatcherName, resourceGroupName)
	}

	d.SetId(*resp.ID)

	return resourceArmNetworkWatcherFlowLogRead(d, meta)
}

func resourceArmNetworkWatcherFlowLogRead(d *schema.ResourceData, meta interface{}) error {
	flowLogsClient := meta.(*clients.Client).Network.FlowLogsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseNetworkWatcherFlowLogID(d.Id())
	if err != nil {
		return err
	}

	resp, err := flowLogsClient.Get(ctx, id.ResourceGroup, id.NetworkWatcherName, id.FlowLogName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Flow Log %q (Resource Group %q): %+v",
			d.Id(), id.ResourceGroup, err)
	}

	d.Set("network_watcher_name", id.NetworkWatcherName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("network_security_group_id", resp.TargetResourceID)

	flowAnalytics := flattenAzureRmNetworkWatcherFlowLogTrafficAnalytics(resp.FlowAnalyticsConfiguration)
	if err := d.Set("traffic_analytics", flowAnalytics); err != nil {
		return fmt.Errorf("Error setting `traffic_analytics`: %+v", err)
	}

	if props := resp.FlowLogPropertiesFormat; props != nil {
		d.Set("enabled", props.Enabled)

		if format := props.Format; format != nil {
			d.Set("version", format.Version)
		}

		// Azure API returns "" when flow log is disabled
		// Don't overwrite to prevent storage account ID diff when that is the case
		if props.StorageID != nil && *props.StorageID != "" {
			d.Set("storage_account_id", props.StorageID)
		}

		if err := d.Set("retention_policy", flattenAzureRmNetworkWatcherFlowLogRetentionPolicy(props.RetentionPolicy)); err != nil {
			return fmt.Errorf("Error setting `retention_policy`: %+v", err)
		}
	}

	return nil
}

func resourceArmNetworkWatcherFlowLogDelete(d *schema.ResourceData, meta interface{}) error {
	flowLogsClient := meta.(*clients.Client).Network.FlowLogsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseNetworkWatcherFlowLogID(d.Id())
	if err != nil {
		return err
	}

	resp, err := flowLogsClient.Get(ctx, id.ResourceGroup, id.NetworkWatcherName, id.FlowLogName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Flow Log %q (Resource Group %q): %+v",
			d.Id(), id.ResourceGroup, err)
	}

	if props := resp.FlowLogPropertiesFormat; props != nil {
		if props.Enabled != nil && *props.Enabled {
			delFuture, err := flowLogsClient.Delete(ctx, id.ResourceGroup, id.NetworkWatcherName, id.FlowLogName)
			if err != nil {
				return fmt.Errorf("Error deleting Flow Log %q: %+v", d.Id(), err)
			}
			if err = delFuture.WaitForCompletionRef(ctx, flowLogsClient.Client); err != nil {
				return fmt.Errorf("Error waiting for completion of Flow Log deletion %q: %+v", d.Id(), err)
			}
		}
	}

	return nil
}

func expandAzureRmNetworkWatcherFlowLogRetentionPolicy(d *schema.ResourceData) *networkNew.RetentionPolicyParameters {
	vs := d.Get("retention_policy").([]interface{})
	if len(vs) < 1 || vs[0] == nil {
		return nil
	}

	v := vs[0].(map[string]interface{})
	enabled := v["enabled"].(bool)
	days := v["days"].(int)

	return &networkNew.RetentionPolicyParameters{
		Enabled: utils.Bool(enabled),
		Days:    utils.Int32(int32(days)),
	}
}

func flattenAzureRmNetworkWatcherFlowLogRetentionPolicy(input *networkNew.RetentionPolicyParameters) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if input.Enabled != nil {
		result["enabled"] = *input.Enabled
	}
	if input.Days != nil {
		result["days"] = *input.Days
	}

	return []interface{}{result}
}

func flattenAzureRmNetworkWatcherFlowLogTrafficAnalytics(input *networkNew.TrafficAnalyticsProperties) []interface{} {
	if input == nil || input.NetworkWatcherFlowAnalyticsConfiguration == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if cfg := input.NetworkWatcherFlowAnalyticsConfiguration; cfg != nil {
		if cfg.Enabled != nil {
			result["enabled"] = *cfg.Enabled
		}
		if cfg.WorkspaceID != nil {
			result["workspace_id"] = *cfg.WorkspaceID
		}
		if cfg.WorkspaceRegion != nil {
			result["workspace_region"] = *cfg.WorkspaceRegion
		}
		if cfg.WorkspaceResourceID != nil {
			result["workspace_resource_id"] = *cfg.WorkspaceResourceID
		}
	}

	return []interface{}{result}
}

func expandAzureRmNetworkWatcherFlowLogTrafficAnalytics(d *schema.ResourceData) *networkNew.TrafficAnalyticsProperties {
	vs := d.Get("traffic_analytics").([]interface{})

	v := vs[0].(map[string]interface{})
	enabled := v["enabled"].(bool)
	workspaceID := v["workspace_id"].(string)
	workspaceRegion := v["workspace_region"].(string)
	workspaceResourceID := v["workspace_resource_id"].(string)

	return &networkNew.TrafficAnalyticsProperties{
		NetworkWatcherFlowAnalyticsConfiguration: &networkNew.TrafficAnalyticsConfigurationProperties{
			Enabled:             utils.Bool(enabled),
			WorkspaceID:         utils.String(workspaceID),
			WorkspaceRegion:     utils.String(workspaceRegion),
			WorkspaceResourceID: utils.String(workspaceResourceID),
		},
	}
}
