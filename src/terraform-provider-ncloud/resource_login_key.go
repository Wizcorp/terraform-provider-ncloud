package main

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

// Virtual machines provided on
// See: https://docs.ncloud.com/en/api_new/api_new-2-1.html
func resourceLoginKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoginKeyCreate,
		Read:   resourceLoginKeyRead,
		Delete: resourceLoginKeyDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key name",
			},
			"private_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLoginKeyCreate(data *schema.ResourceData, meta interface{}) error {
	data.Partial(true)
	client := meta.(*sdk.Conn)
	name := data.Get("name").(string)

	response, err := client.CreateLoginKey(name)
	if err != nil {
		return fmt.Errorf("Failed to create login key %s", err)
	}

	data.SetId(name)
	data.Set("private_key", response.PrivateKey)
	data.SetPartial("private_key")

	return nil
}

func resourceLoginKeyRead(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*sdk.Conn)

	reqParams := new(sdk.RequestGetLoginKeyList)
	reqParams.KeyName = data.Id()

	response, err := client.GetLoginKeyList(reqParams)
	if err != nil {
		return fmt.Errorf("Failed to fetch key info %s", err)
	}

	keyInfo := response.LoginKeyList[0]
	data.Set("fingerprint", keyInfo.Fingerprint)
	data.SetPartial("fingerprint")

	return nil
}

func resourceLoginKeyDelete(data *schema.ResourceData, meta interface{}) error {
	client := meta.(*sdk.Conn)

	reqParams := new(sdk.RequestTerminateServerInstances)
	reqParams.ServerInstanceNoList = []string{}

	_, err := client.TerminateServerInstances(reqParams)
	if err != nil {
		return fmt.Errorf("Failed to delete login key %s", err)
	}

	return nil
}
