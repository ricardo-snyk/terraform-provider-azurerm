package keyvault

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultSecretCreate,
		Read:   resourceArmKeyVaultSecretRead,
		Update: resourceArmKeyVaultSecretUpdate,
		Delete: resourceArmKeyVaultSecretDelete,
		Importer: &schema.ResourceImporter{
			State: resourceArmKeyVaultChildResourceImporter,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateKeyVaultChildName,
			},

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"value": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"not_before_date": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.RFC3339Time,
			},

			"expiration_date": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.RFC3339Time,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmKeyVaultSecretCreate(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Print("[INFO] preparing arguments for AzureRM KeyVault Secret creation.")

	name := d.Get("name").(string)
	keyVaultId := d.Get("key_vault_id").(string)

	keyVaultBaseUrl, err := azure.GetKeyVaultBaseUrlFromID(ctx, vaultClient, keyVaultId)
	if err != nil {
		return fmt.Errorf("Error looking up Secret %q vault url from id %q: %+v", name, keyVaultId, err)
	}

	if features.ShouldResourcesBeImported() {
		existing, err := client.GetSecret(ctx, keyVaultBaseUrl, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Secret %q (Key Vault %q): %s", name, keyVaultBaseUrl, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_secret", *existing.ID)
		}
	}

	value := d.Get("value").(string)
	contentType := d.Get("content_type").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := keyvault.SecretSetParameters{
		Value:            utils.String(value),
		ContentType:      utils.String(contentType),
		Tags:             tags.Expand(t),
		SecretAttributes: &keyvault.SecretAttributes{},
	}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) //validated by schema
		notBeforeUnixTime := date.UnixTime(notBeforeDate)
		parameters.SecretAttributes.NotBefore = &notBeforeUnixTime
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) //validated by schema
		expirationUnixTime := date.UnixTime(expirationDate)
		parameters.SecretAttributes.Expires = &expirationUnixTime
	}

	if _, err := client.SetSecret(ctx, keyVaultBaseUrl, name, parameters); err != nil {
		return err
	}

	// "" indicates the latest version
	read, err := client.GetSecret(ctx, keyVaultBaseUrl, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read KeyVault Secret '%s' (in key vault '%s')", name, keyVaultBaseUrl)
	}

	d.SetId(*read.ID)

	return resourceArmKeyVaultSecretRead(d, meta)
}

func resourceArmKeyVaultSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	keyVaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Print("[INFO] preparing arguments for AzureRM KeyVault Secret update.")

	id, err := azure.ParseKeyVaultChildID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := azure.GetKeyVaultIDFromBaseUrl(ctx, keyVaultClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultId == nil {
		return fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}

	ok, err := azure.KeyVaultExists(ctx, keyVaultClient, *keyVaultId)
	if err != nil {
		return fmt.Errorf("Error checking if key vault %q for Secret %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	value := d.Get("value").(string)
	contentType := d.Get("content_type").(string)
	t := d.Get("tags").(map[string]interface{})

	secretAttributes := &keyvault.SecretAttributes{}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) //validated by schema
		notBeforeUnixTime := date.UnixTime(notBeforeDate)
		secretAttributes.NotBefore = &notBeforeUnixTime
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) //validated by schema
		expirationUnixTime := date.UnixTime(expirationDate)
		secretAttributes.Expires = &expirationUnixTime
	}

	if d.HasChange("value") {
		// for changing the value of the secret we need to create a new version
		parameters := keyvault.SecretSetParameters{
			Value:            utils.String(value),
			ContentType:      utils.String(contentType),
			Tags:             tags.Expand(t),
			SecretAttributes: secretAttributes,
		}

		if _, err = client.SetSecret(ctx, id.KeyVaultBaseUrl, id.Name, parameters); err != nil {
			return err
		}

		// "" indicates the latest version
		read, err2 := client.GetSecret(ctx, id.KeyVaultBaseUrl, id.Name, "")
		if err2 != nil {
			return fmt.Errorf("Error getting Key Vault Secret %q : %+v", id.Name, err2)
		}

		if _, err = azure.ParseKeyVaultChildID(*read.ID); err != nil {
			return err
		}

		// the ID is suffixed with the secret version
		d.SetId(*read.ID)
	} else {
		parameters := keyvault.SecretUpdateParameters{
			ContentType:      utils.String(contentType),
			Tags:             tags.Expand(t),
			SecretAttributes: secretAttributes,
		}

		if _, err = client.UpdateSecret(ctx, id.KeyVaultBaseUrl, id.Name, id.Version, parameters); err != nil {
			return err
		}
	}

	return resourceArmKeyVaultSecretRead(d, meta)
}

func resourceArmKeyVaultSecretRead(d *schema.ResourceData, meta interface{}) error {
	keyVaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseKeyVaultChildID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := azure.GetKeyVaultIDFromBaseUrl(ctx, keyVaultClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultId == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	ok, err := azure.KeyVaultExists(ctx, keyVaultClient, *keyVaultId)
	if err != nil {
		return fmt.Errorf("Error checking if key vault %q for Secret %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	// Previously at this point there was a call to client.GetSecret. Avoid that
	// since we do NOT want access to the secret's value and we want to avoid
	// failing when that permission is not available to us. Instead, we can use
	// the client.GetSecrets call which requires only "list" permissions and
	// that response contains enough information for us to populate the secret
	// attributes we care about. This implementation has the downside that each
	// secret will trigger a GetSecrets request which is not efficient. The
	// response could be cached in the future, as one way to mitigate this.
	maxResults := int32(25)
	resp, err := client.GetSecrets(ctx, id.KeyVaultBaseUrl, &maxResults)
	if err != nil {
		return fmt.Errorf("Error listing Key Vault Secrets: %s", err)
	}

	// Find this secret in the list of all secrets in this Key Vault
	var secret *keyvault.SecretItem
	for _, item := range resp.Values() {
		// Parse the secret URL to extract the name. A newer version of
		// the Azure SDK Go library may expose a name attribute directly,
		// since a name attribute IS returned in the response under the hood.
		secretParts := strings.Split(strings.TrimSuffix(*item.ID, "/"), "/")
		if len(secretParts) < 4 {
			return fmt.Errorf("Unexpected Key Vault Secret ID format: %+v", *item.ID)
		}
		secretName := secretParts[len(secretParts)-1]
		if secretName == id.Name {
			secret = &item
			break
		}
	}

	// The secret was not found in the Key Vault
	if secret == nil {
		log.Printf("[DEBUG] Unable to find secret: %v", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("key_vault_id", *keyVaultId)
	d.Set("name", id.Name)
	d.Set("version", id.Version)
	d.Set("content_type", secret.ContentType)

	if attributes := secret.Attributes; attributes != nil {
		if v := attributes.NotBefore; v != nil {
			d.Set("not_before_date", time.Time(*v).Format(time.RFC3339))
		}

		if v := attributes.Expires; v != nil {
			d.Set("expiration_date", time.Time(*v).Format(time.RFC3339))
		}
	}

	return tags.FlattenAndSet(d, secret.Tags)
}

func resourceArmKeyVaultSecretDelete(d *schema.ResourceData, meta interface{}) error {
	keyVaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseKeyVaultChildID(d.Id())
	if err != nil {
		return err
	}

	keyVaultId, err := azure.GetKeyVaultIDFromBaseUrl(ctx, keyVaultClient, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("Error retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultId == nil {
		return fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", id.KeyVaultBaseUrl)
	}

	ok, err := azure.KeyVaultExists(ctx, keyVaultClient, *keyVaultId)
	if err != nil {
		return fmt.Errorf("Error checking if key vault %q for Secret %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	_, err = client.DeleteSecret(ctx, id.KeyVaultBaseUrl, id.Name)
	return err
}
