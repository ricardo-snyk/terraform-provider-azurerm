package helper

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func ExtendedAuditingSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"storage_endpoint": {
					Type:     schema.TypeString,
					Required: true,
					//ValidateFunc: validation.IsURLWithHTTPS,
				},

				"storage_account_access_key_is_secondary": {
					Type:     schema.TypeBool,
					Optional: true,
				},

				"retention_in_days": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 3285),
				},
			},
		},
	}
}

func ExpandAzureRmSqlServerBlobAuditingPolicies(input []interface{}) *sql.ExtendedServerBlobAuditingPolicyProperties {
	if len(input) == 0 {
		return &sql.ExtendedServerBlobAuditingPolicyProperties{
			State: sql.BlobAuditingPolicyStateDisabled,
		}
	}
	serverBlobAuditingPolicies := input[0].(map[string]interface{})

	ExtendedServerBlobAuditingPolicyProperties := sql.ExtendedServerBlobAuditingPolicyProperties{
		State:                   sql.BlobAuditingPolicyStateEnabled,
		StorageAccountAccessKey: utils.String(serverBlobAuditingPolicies["storage_account_access_key"].(string)),
		StorageEndpoint:         utils.String(serverBlobAuditingPolicies["storage_endpoint"].(string)),
	}
	if v, ok := serverBlobAuditingPolicies["storage_account_access_key_is_secondary"]; ok {
		ExtendedServerBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse = utils.Bool(v.(bool))
	}
	if v, ok := serverBlobAuditingPolicies["retention_in_days"]; ok {
		ExtendedServerBlobAuditingPolicyProperties.RetentionDays = utils.Int32(int32(v.(int)))
	}

	return &ExtendedServerBlobAuditingPolicyProperties
}

func FlattenAzureRmSqlServerBlobAuditingPolicies(extendedServerBlobAuditingPolicy *sql.ExtendedServerBlobAuditingPolicy, d *schema.ResourceData) []interface{} {
	if extendedServerBlobAuditingPolicy == nil || extendedServerBlobAuditingPolicy.State == sql.BlobAuditingPolicyStateDisabled {
		return []interface{}{}
	}
	var storageEndpoint string
	if extendedServerBlobAuditingPolicy.StorageEndpoint != nil {
		storageEndpoint = *extendedServerBlobAuditingPolicy.StorageEndpoint
	}

	var secondKeyInUse bool
	if extendedServerBlobAuditingPolicy.IsStorageSecondaryKeyInUse != nil {
		secondKeyInUse = *extendedServerBlobAuditingPolicy.IsStorageSecondaryKeyInUse
	}
	var retentionDays int32
	if extendedServerBlobAuditingPolicy.RetentionDays != nil {
		retentionDays = *extendedServerBlobAuditingPolicy.RetentionDays
	}

	return []interface{}{
		map[string]interface{}{
			"storage_endpoint":                        storageEndpoint,
			"storage_account_access_key_is_secondary": secondKeyInUse,
			"retention_in_days":                       retentionDays,
		},
	}
}

func ExpandAzureRmSqlDBBlobAuditingPolicies(input []interface{}) *sql.ExtendedDatabaseBlobAuditingPolicyProperties {
	if len(input) == 0 {
		return &sql.ExtendedDatabaseBlobAuditingPolicyProperties{
			State: sql.BlobAuditingPolicyStateDisabled,
		}
	}
	dbBlobAuditingPolicies := input[0].(map[string]interface{})

	ExtendedDatabaseBlobAuditingPolicyProperties := sql.ExtendedDatabaseBlobAuditingPolicyProperties{
		State:                   sql.BlobAuditingPolicyStateEnabled,
		StorageAccountAccessKey: utils.String(dbBlobAuditingPolicies["storage_account_access_key"].(string)),
		StorageEndpoint:         utils.String(dbBlobAuditingPolicies["storage_endpoint"].(string)),
	}
	if v, ok := dbBlobAuditingPolicies["storage_account_access_key_is_secondary"]; ok {
		ExtendedDatabaseBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse = utils.Bool(v.(bool))
	}
	if v, ok := dbBlobAuditingPolicies["retention_in_days"]; ok {
		ExtendedDatabaseBlobAuditingPolicyProperties.RetentionDays = utils.Int32(int32(v.(int)))
	}

	return &ExtendedDatabaseBlobAuditingPolicyProperties
}

func FlattenAzureRmSqlDBBlobAuditingPolicies(extendedDatabaseBlobAuditingPolicy *sql.ExtendedDatabaseBlobAuditingPolicy, d *schema.ResourceData) []interface{} {
	if extendedDatabaseBlobAuditingPolicy == nil || extendedDatabaseBlobAuditingPolicy.State == sql.BlobAuditingPolicyStateDisabled {
		return []interface{}{}
	}
	var storageEndpoint string

	if extendedDatabaseBlobAuditingPolicy.StorageEndpoint != nil {
		storageEndpoint = *extendedDatabaseBlobAuditingPolicy.StorageEndpoint
	}
	var secondKeyInUse bool
	if extendedDatabaseBlobAuditingPolicy.IsStorageSecondaryKeyInUse != nil {
		secondKeyInUse = *extendedDatabaseBlobAuditingPolicy.IsStorageSecondaryKeyInUse
	}
	var retentionDays int32
	if extendedDatabaseBlobAuditingPolicy.RetentionDays != nil {
		retentionDays = *extendedDatabaseBlobAuditingPolicy.RetentionDays
	}

	return []interface{}{
		map[string]interface{}{
			"storage_endpoint":                        storageEndpoint,
			"storage_account_access_key_is_secondary": secondKeyInUse,
			"retention_in_days":                       retentionDays,
		},
	}
}
