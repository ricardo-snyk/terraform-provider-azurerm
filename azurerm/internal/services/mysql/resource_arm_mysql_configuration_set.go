// This is a Fugue internal resource type.  
// It aggregates the values from individual Configuration resources into one map.

package mysql

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMySQLConfigurationSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMySQLConfigurationSetCreate,
		Read:   resourceArmMySQLConfigurationSetRead,
		Delete: resourceArmMySQLConfigurationSetDelete,
		Update: resourceArmMySQLConfigurationSetCreate,

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
			"config_map": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				ForceNew: false,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMySqlServerName,
			},
		},
	}
}

func resourceArmMySQLConfigurationSetCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmMySQLConfigurationSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	resp, err := client.ListByServer(ctx, resourceGroup, serverName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] MySQL Server %q was not found (resource group %q)", serverName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making List Configuration request on Azure MySQL Server %q (Resource Group %q): %+v", serverName, resourceGroup, err)
	}
	configMap := make(map[string]interface{})
	configs := resp.Value
	for _, conf := range *configs {
		key := conf.Name
		value := conf.ConfigurationProperties.Value
		configMap[*key] = *value
	}
	d.Set("server_name", serverName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("config_map", configMap)

	return nil
}

func resourceArmMySQLConfigurationSetDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
