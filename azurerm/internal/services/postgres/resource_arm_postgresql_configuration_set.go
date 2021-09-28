// This is a Fugue internal resource type.
// It aggregates the values from individual Configuration resources into one map.

package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPostgreSQLConfigurationSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPostgreSQLConfigurationSetCreate,
		Read:   resourceArmPostgreSQLConfigurationSetRead,
		Delete: resourceArmPostgreSQLConfigurationSetDelete,
		Update: resourceArmPostgreSQLConfigurationSetCreate,

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
				ValidateFunc: ValidatePSQLServerName,
			},
		},
	}
}

func resourceArmPostgreSQLConfigurationSetCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmPostgreSQLConfigurationSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ConfigurationsClient
	// Manually creating the context instead of using timeouts.ForRead. The timeouts in
	// the schema above do not get used here. The way that this ResourceData gets
	// constructed by the plugin SDK, all the timeouts are set to nil, resulting in a
	// 450-second timeout (this timeout might be on Azure's side). Creating the context
	// this way reliably gets us a 30-second timeout.
	ctx, cancel := context.WithTimeout(meta.(*clients.Client).StopContext, 30*time.Second)
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
			log.Printf("[WARN] PostgreSQL Server %q was not found (resource group %q)", serverName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making List Configuration request on Azure PostgreSQL Server %q (Resource Group %q): %+v", serverName, resourceGroup, err)
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

func resourceArmPostgreSQLConfigurationSetDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
