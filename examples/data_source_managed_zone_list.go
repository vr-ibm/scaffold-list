
package dns

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func DataSourceManagedZoneList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceManagedZoneListRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"managed_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						
"description": {
	Type:     schema.TypeString,
	Computed: true,
},
						
"dns_name": {
	Type:     schema.TypeString,
	Computed: true,
},
						
"dnssec_config": {
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			
"kind": {
	Type:     schema.TypeString,
	Computed: true,
},
			
"non_existence": {
	Type:     schema.TypeString,
	Computed: true,
},
			
"state": {
	Type:     schema.TypeString,
	Computed: true,
},
			
"default_key_specs": {
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Schema{Type: schema.TypeString},
},
		},
	},
},
						
"managed_zone_id": {
	Type:     schema.TypeInt,
	Computed: true,
},
						
"name": {
	Type:     schema.TypeString,
	Computed: true,
},
						
"name_servers": {
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Schema{Type: schema.TypeString},
},
						
"creation_time": {
	Type:     schema.TypeString,
	Computed: true,
},
						
"labels": {
	Type:     schema.TypeString,
	Computed: true,
},
						
"visibility": {
	Type:     schema.TypeString,
	Computed: true,
},
						
"private_visibility_config": {
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			
"gke_clusters": {
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Schema{Type: schema.TypeString},
},
			
"networks": {
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Schema{Type: schema.TypeString},
},
		},
	},
},
						
"forwarding_config": {
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			
"target_name_servers": {
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Schema{Type: schema.TypeString},
},
		},
	},
},
						
"peering_config": {
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			
"target_network": {
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			
"network_url": {
	Type:     schema.TypeString,
	Computed: true,
},
		},
	},
},
		},
	},
},
						
"reverse_lookup": {
	Type:     schema.TypeBool,
	Computed: true,
},
						
"service_directory_config": {
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			
"namespace": {
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			
"namespace_url": {
	Type:     schema.TypeString,
	Computed: true,
},
		},
	},
},
		},
	},
},
						
"cloud_logging_config": {
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			
"enable_logging": {
	Type:     schema.TypeBool,
	Computed: true,
},
		},
	},
},
					},
				},
			},
		},
	}
}

func init() {
	registry.Schema{
		Name:        "google_dns_managed_zone_list",
		ProductName: "dns",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceManagedZoneList(),
	}.Register()
}

func dataSourceManagedZoneListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return diag.FromErr(err)
	}
	baseURL := fmt.Sprintf("https://dns.googleapis.com/dns/v1/projects/%s/managedZones", project)
	var items []interface{}
	token := ""
	for {
		url := baseURL
		if token != "" {
			url = fmt.Sprintf("%s?pageToken=%s", baseURL, token)
		}
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    url,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("error listing ManagedZones: %w", err))
		}
		pageItems, ok := res["managedZones"].([]interface{})
		if !ok {
			break
		}
		for _, raw := range pageItems {
			item, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			items = append(items, flattenManagedZoneListItem(item, project))
		}
		token, _ = res["nextPageToken"].(string)
		if token == "" {
			break
		}
	}
	if err := d.Set("managed_zones", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("projects/%s/managedZones", project))
	return nil
}

func flattenManagedZoneListItem(item map[string]interface{}, project string) map[string]interface{} {
	name, _ := item["name"].(string)
	idStr, _ := item["id"].(string)
	idInt, _ := strconv.ParseInt(idStr, 10, 64)
	return map[string]interface{}{
		"id":              fmt.Sprintf("projects/%s/managedZones/%s", project, name),
		"name":            name,
		"dns_name":        item["dnsName"],
		"description":     item["description"],
		"visibility":      item["visibility"],
		"managed_zone_id": idInt,
		"project":         project,
	}
}
