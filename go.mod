module github.com/terraform-providers/terraform-provider-azurerm

require (
	github.com/Azure/azure-sdk-for-go v39.3.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.3
	github.com/Azure/go-autorest/autorest/date v0.2.0
	github.com/btubbs/datetime v0.1.0
	github.com/davecgh/go-spew v1.1.1
	tfratelimiter v0.0.0
	github.com/google/uuid v1.1.2
	github.com/hashicorp/go-azure-helpers v0.10.0
	github.com/hashicorp/go-getter v1.5.3
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/go-version v1.3.0
	github.com/hashicorp/terraform-plugin-sdk v1.6.0
	github.com/satori/go.uuid v1.2.0
	github.com/satori/uuid v0.0.0-20160927100844-b061729afc07
	github.com/terraform-providers/terraform-provider-azuread v0.6.1-0.20191007035844-361c0a206ad4
	github.com/tombuildsstuff/giovanni v0.7.1
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b
	golang.org/x/net v0.0.0-20210326060303-6b1517762897
	gopkg.in/yaml.v2 v2.3.0
)

replace tfratelimiter => ../tfratelimiter

go 1.13
