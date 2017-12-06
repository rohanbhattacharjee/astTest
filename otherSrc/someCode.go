// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package otherSrc

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

var (
	FiveMinutes = 5 * time.Minute

	defaultTimeout = &schema.ResourceTimeout{
		Create: &FiveMinutes,
		Update: &FiveMinutes,
		Delete: &FiveMinutes,
	}
)

type ConsoleHistoryResourceCrud struct {
	var1 int
	var2 string
}

func ConsoleHistoryResource() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: defaultTimeout,
		Create:   createConsoleHistory,
		Read:     readConsoleHistory,
		Delete:   deleteConsoleHistory,
		Schema: map[string]*schema.Schema{
			"availability_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compartment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createConsoleHistory(d *schema.ResourceData, m interface{}) (e error) {
	e = nil
	return
}

func readConsoleHistory(d *schema.ResourceData, m interface{}) (e error) {
	e = nil
	return
}

func deleteConsoleHistory(d *schema.ResourceData, m interface{}) (e error) {
	e = nil
	return
}

func (s *ConsoleHistoryResourceCrud) ID() string {
	return "foobar5"
}

func (s *ConsoleHistoryResourceCrud) CreatedPending() []string {
	return []string{"foobar3", "foobar4"}
}

func (s *ConsoleHistoryResourceCrud) CreatedTarget() []string {
	return []string{"foobar2"}
}

func (s *ConsoleHistoryResourceCrud) State() string {
	return "foobar1"
}

func (s *ConsoleHistoryResourceCrud) Create() (e error) {
	e = nil
	return
}

func (s *ConsoleHistoryResourceCrud) Get() (e error) {
	e = nil
	return
}

func (s *ConsoleHistoryResourceCrud) Delete() (e error) {
	e = nil
	return
}

func (s *ConsoleHistoryResourceCrud) SetData() {
}
