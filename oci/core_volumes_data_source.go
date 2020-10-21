// Copyright (c) 2017, 2020, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oci_core "github.com/oracle/oci-go-sdk/v27/core"
)

func init() {
	RegisterDatasource("oci_core_volumes", CoreVolumesDataSource())
}

func CoreVolumesDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readCoreVolumes,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"availability_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"compartment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     GetDataSourceItemSchema(CoreVolumeResource()),
			},
		},
	}
}

func readCoreVolumes(d *schema.ResourceData, m interface{}) error {
	sync := &CoreVolumesDataSourceCrud{}
	sync.D = d
	sync.Client = m.(*OracleClients).blockstorageClient()

	return ReadResource(sync)
}

type CoreVolumesDataSourceCrud struct {
	D      *schema.ResourceData
	Client *oci_core.BlockstorageClient
	Res    *oci_core.ListVolumesResponse
}

func (s *CoreVolumesDataSourceCrud) VoidState() {
	s.D.SetId("")
}

func (s *CoreVolumesDataSourceCrud) Get() error {
	request := oci_core.ListVolumesRequest{}

	if availabilityDomain, ok := s.D.GetOkExists("availability_domain"); ok {
		tmp := availabilityDomain.(string)
		request.AvailabilityDomain = &tmp
	}

	if compartmentId, ok := s.D.GetOkExists("compartment_id"); ok {
		tmp := compartmentId.(string)
		request.CompartmentId = &tmp
	}

	if displayName, ok := s.D.GetOkExists("display_name"); ok {
		tmp := displayName.(string)
		request.DisplayName = &tmp
	}

	if state, ok := s.D.GetOkExists("state"); ok {
		request.LifecycleState = oci_core.VolumeLifecycleStateEnum(state.(string))
	}

	if volumeGroupId, ok := s.D.GetOkExists("volume_group_id"); ok {
		tmp := volumeGroupId.(string)
		request.VolumeGroupId = &tmp
	}

	request.RequestMetadata.RetryPolicy = getRetryPolicy(false, "core")

	response, err := s.Client.ListVolumes(context.Background(), request)
	if err != nil {
		return err
	}

	s.Res = &response
	request.Page = s.Res.OpcNextPage

	for request.Page != nil {
		listResponse, err := s.Client.ListVolumes(context.Background(), request)
		if err != nil {
			return err
		}

		s.Res.Items = append(s.Res.Items, listResponse.Items...)
		request.Page = listResponse.OpcNextPage
	}

	return nil
}

func (s *CoreVolumesDataSourceCrud) SetData() error {
	if s.Res == nil {
		return nil
	}

	s.D.SetId(GenerateDataSourceHashID("CoreVolumesDataSource-", CoreVolumesDataSource(), s.D))
	resources := []map[string]interface{}{}

	for _, r := range s.Res.Items {
		volume := map[string]interface{}{
			"compartment_id": *r.CompartmentId,
		}

		if r.AutoTunedVpusPerGB != nil {
			volume["auto_tuned_vpus_per_gb"] = strconv.FormatInt(*r.AutoTunedVpusPerGB, 10)
		}

		if r.AvailabilityDomain != nil {
			volume["availability_domain"] = *r.AvailabilityDomain
		}

		if r.DefinedTags != nil {
			volume["defined_tags"] = definedTagsToMap(r.DefinedTags)
		}

		if r.DisplayName != nil {
			volume["display_name"] = *r.DisplayName
		}

		volume["freeform_tags"] = r.FreeformTags

		if r.Id != nil {
			volume["id"] = *r.Id
		}

		if r.IsAutoTuneEnabled != nil {
			volume["is_auto_tune_enabled"] = *r.IsAutoTuneEnabled
		}

		if r.IsHydrated != nil {
			volume["is_hydrated"] = *r.IsHydrated
		}

		if r.KmsKeyId != nil {
			volume["kms_key_id"] = *r.KmsKeyId
		}

		if r.SizeInGBs != nil {
			volume["size_in_gbs"] = strconv.FormatInt(*r.SizeInGBs, 10)
		}

		if r.SizeInMBs != nil {
			volume["size_in_mbs"] = strconv.FormatInt(*r.SizeInMBs, 10)
		}

		if r.SourceDetails != nil {
			sourceDetailsArray := []interface{}{}
			if sourceDetailsMap := VolumeSourceDetailsToMap(&r.SourceDetails); sourceDetailsMap != nil {
				sourceDetailsArray = append(sourceDetailsArray, sourceDetailsMap)
			}
			volume["source_details"] = sourceDetailsArray
		} else {
			volume["source_details"] = nil
		}

		volume["state"] = r.LifecycleState

		if r.SystemTags != nil {
			volume["system_tags"] = systemTagsToMap(r.SystemTags)
		}

		if r.TimeCreated != nil {
			volume["time_created"] = r.TimeCreated.String()
		}

		if r.VolumeGroupId != nil {
			volume["volume_group_id"] = *r.VolumeGroupId
		}

		if r.VpusPerGB != nil {
			volume["vpus_per_gb"] = strconv.FormatInt(*r.VpusPerGB, 10)
		}

		resources = append(resources, volume)
	}

	if f, fOk := s.D.GetOkExists("filter"); fOk {
		resources = ApplyFilters(f.(*schema.Set), resources, CoreVolumesDataSource().Schema["volumes"].Elem.(*schema.Resource).Schema)
	}

	if err := s.D.Set("volumes", resources); err != nil {
		return err
	}

	return nil
}
