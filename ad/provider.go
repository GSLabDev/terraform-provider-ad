package ad

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{

			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Domain in which AD Server resides",
				DefaultFunc: schema.EnvDefaultFunc("AD_DOMAIN", nil),
			},

			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP of the AD Server",
				DefaultFunc: schema.EnvDefaultFunc("AD_IP", nil),
			},

			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user name of the AD Server",
				DefaultFunc: schema.EnvDefaultFunc("AD_USER", nil),
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The user password of the AD Server",
				DefaultFunc: schema.EnvDefaultFunc("AD_PASSWORD", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"ad_computer":       resourceComputer(),
			"ad_computer_to_ou": resourceComputerToOU(),
			"ad_group_to_ou":    resourceGroupToOU(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	config := Config{
		Domain:   d.Get("domain").(string),
		IP:       d.Get("ip").(string),
		Username: d.Get("user").(string),
		Password: d.Get("password").(string),
	}
	log.Printf("[DEBUG] Connecting to AD")
	return config.Client()
}
