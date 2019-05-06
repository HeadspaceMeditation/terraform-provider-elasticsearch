package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/olivere/elastic"
)

func resourceElasticsearchIndexTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceElasticsearchIndexTemplateCreate,
		Read:   resourceElasticsearchIndexTemplateRead,
		Update: resourceElasticsearchIndexTemplateUpdate,
		Delete: resourceElasticsearchIndexTemplateDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"body": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: diffSuppressIndexTemplate,
			},
		},
	}
}

func resourceElasticsearchIndexTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	err := resourceElasticsearchPutIndexTemplate(d, meta, true)
	if err != nil {
		return err
	}
	d.SetId(d.Get("name").(string))
	return nil
}

func resourceElasticsearchIndexTemplateRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	var result string
	var err error
	switch meta.(type) {
	case *elastic.Client:
		client := meta.(*elastic.Client)
		result, err = elastic6IndexGetTemplate(client, id)
	default:
		return errors.New("No elastic client made")
	}
	if err != nil {
		return err
	}

	d.Set("name", d.Id())
	d.Set("body", result)
	return nil
}

func elastic6IndexGetTemplate(client *elastic.Client, id string) (string, error) {
	res, err := client.IndexGetTemplate(id).Do(context.Background())
	if err != nil {
		return "", err
	}

	t := res[id]
	tj, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(tj), nil
}

func resourceElasticsearchIndexTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceElasticsearchPutIndexTemplate(d, meta, false)
}

func resourceElasticsearchIndexTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	var err error
	switch meta.(type) {
	case *elastic.Client:
		client := meta.(*elastic.Client)
		err = elastic6IndexDeleteTemplate(client, id)
	default:
		return errors.New("No elastic client made")
	}

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func elastic6IndexDeleteTemplate(client *elastic.Client, id string) error {
	_, err := client.IndexDeleteTemplate(id).Do(context.Background())
	return err
}

func resourceElasticsearchPutIndexTemplate(d *schema.ResourceData, meta interface{}, create bool) error {
	name := d.Get("name").(string)
	body := d.Get("body").(string)

	var err error
	switch meta.(type) {
	case *elastic.Client:
		client := meta.(*elastic.Client)
		err = elastic6IndexPutTemplate(client, name, body, create)
	default:
		return errors.New("No elastic client made")
	}

	return err
}

func elastic6IndexPutTemplate(client *elastic.Client, name string, body string, create bool) error {
	_, err := client.IndexPutTemplate(name).BodyString(body).Create(create).Do(context.Background())
	return err
}
