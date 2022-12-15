package s3control

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3control"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func DataSourceMultiRegionAccessPoint() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceMultiRegionAccessPointBlockRead,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: verify.ValidAccountID,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_access_block": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMultiRegionAccessPointBlockRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).S3ControlConn

	accountID := meta.(*conns.AWSClient).AccountID
	if v, ok := d.GetOk("account_id"); ok {
		accountID = v.(string)
	}

	name := d.Get("name").(string)

	input := &s3control.GetMultiRegionAccessPointInput{
		AccountId: aws.String(accountID),
		Name:      aws.String(name),
	}

	log.Printf("[DEBUG] Reading S3 Multi Region Access Point: %s", input)

	output, err := conn.GetMultiRegionAccessPoint(input)

	if err != nil {
		return diag.Errorf("error reading S3 Multi Region Access Point: %s", err)
	}

	if output == nil || output.AccessPoint == nil {
		return diag.Errorf("error reading S3 Multi Region Access Point (%s): missing access point", accountID)
	}

	d.SetId(MultiRegionAccessPointCreateResourceID(accountID, name))
	d.Set("created_at", aws.TimeValue(output.AccessPoint.CreatedAt).Format(time.RFC3339))
	d.Set("name", output.AccessPoint.Name)
	d.Set("public_access_block", output.AccessPoint.PublicAccessBlock)
	d.Set("regions", output.AccessPoint.Regions)
	d.Set("status", output.AccessPoint.Status)

	return nil
}
