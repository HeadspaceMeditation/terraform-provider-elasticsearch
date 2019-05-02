package main

import (
	"net/url"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/olivere/elastic"
)

var awsUrlRegexp = regexp.MustCompile(`([a-z0-9-]+).es.amazonaws.com$`)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ELASTICSEARCH_URL", nil),
				Description: "Elasticsearch URL",
			},

			"aws_access_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The access key for use with AWS Elasticsearch Service domains",
			},

			"aws_secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The secret key for use with AWS Elasticsearch Service domains",
			},

			"aws_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The session token for use with AWS Elasticsearch Service domains",
			},

			"cacert_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A Custom CA certificate",
			},

			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable SSL verification of API calls",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"elasticsearch_index_template": resourceElasticsearchIndexTemplate(),
			// "elasticsearch_snapshot_repository": resourceElasticsearchSnapshotRepository(),
			// "elasticsearch_kibana_object":       resourceElasticsearchKibanaObject(),
		},

		ConfigureFunc: providerConfigure,
	}

}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	rawUrl := d.Get("url").(string)
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(rawUrl),
		elastic.SetScheme(parsedUrl.Scheme),
	}

	return elastic.NewSimpleClient(opts...)
}
