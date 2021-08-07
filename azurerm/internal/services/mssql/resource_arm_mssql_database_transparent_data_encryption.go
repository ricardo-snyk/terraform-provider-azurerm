package mssql

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type TransparentDataEncryptionId struct {
	SubscriptionId string
	ResourceGroup  string
	ServerName     string
	DatabaseName   string
	Name           string
}

func resourceArmMsSqlDatabaseTransparentDataEncryption() *schema.Resource {
	return &schema.Resource{
		Create: resourceMsSqlDatabaseTransparentDataEncryptionCreateUpdate,
		Read:   resourceMsSqlDatabaseTransparentDataEncryptionRead,
		Update: resourceMsSqlDatabaseTransparentDataEncryptionCreateUpdate,
		Delete: resourceMsSqlDatabaseTransparentDataEncryptionDelete,
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
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceMsSqlDatabaseTransparentDataEncryptionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMsSqlDatabaseTransparentDataEncryptionRead(d *schema.ResourceData, meta interface{}) error {
	transparentDataEncryptionsClient := meta.(*clients.Client).MSSQL.TransparentDataEncryptionsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseTransparentDataEncryptionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := transparentDataEncryptionsClient.Get(ctx, id.ResourceGroup, id.ServerName, id.DatabaseName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request for %s: %v", id, err)
	}

	serverId := parse.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)
	d.Set("server_id", serverId.ID())

	databaseId := parse.NewDatabaseID(id.SubscriptionId, id.ResourceGroup, id.ServerName, id.DatabaseName)
	d.Set("database_id", databaseId.ID())

	state := ""
	if resp.TransparentDataEncryptionProperties != nil && resp.TransparentDataEncryptionProperties.Status != "" {
		state = string(resp.TransparentDataEncryptionProperties.Status)
	}
	if err := d.Set("state", state); err != nil {
		return fmt.Errorf("setting state`: %+v", err)
	}

	return nil
}

func resourceMsSqlDatabaseTransparentDataEncryptionDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func NewTransparentDataEncryptiontID(subscriptionId, resourceGroup, serverName, database, name string) TransparentDataEncryptionId {
	return TransparentDataEncryptionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServerName:     serverName,
		DatabaseName:   database,
		Name:           name,
	}
}

func parseTransparentDataEncryptionID(input string) (*TransparentDataEncryptionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := TransparentDataEncryptionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}
	if resourceId.DatabaseName, err = id.PopSegment("databases"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("transparentDataEncryption"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
