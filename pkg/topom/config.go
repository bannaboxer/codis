// Copyright 2016 CodisLabs. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package topom

import (
	"bytes"

	"github.com/BurntSushi/toml"

	"github.com/CodisLabs/codis/pkg/models"
	"github.com/CodisLabs/codis/pkg/utils/bytesize"
	"github.com/CodisLabs/codis/pkg/utils/errors"
	"github.com/CodisLabs/codis/pkg/utils/log"
	"github.com/CodisLabs/codis/pkg/utils/timesize"
)

const DefaultConfig = `
##################################################
#                                                #
#                  Codis-Dashboard               #
#                                                #
##################################################

# Set Coordinator, only accept "zookeeper" & "etcd" & "filesystem".
coordinator_name = "zookeeper"
coordinator_addr = "127.0.0.1:2181"

# Set Codis Product Name/Auth.
product_name = "codis-demo"
product_auth = ""

# Set bind address for admin(rpc), tcp only.
admin_addr = "0.0.0.0:18080"

# Set arguments for data migration (only accept 'sync' & 'semi-async').
forward_method = "semi-async"
migrate_async_maxbulks = 1000
migrate_async_maxbytes = "128k"
migrate_async_pipeline = 8
migrate_async_numkeys = 32

# Set configs for redis sentinel.
sentinel_quorum = 2
sentinel_parallel_syncs = 1
sentinel_down_after = "30s"
sentinel_failover_timeout = "5m"
sentinel_notification_script = ""
sentinel_client_reconfig_script = ""
`

type Config struct {
	CoordinatorName string `toml:"coordinator_name" json:"coordinator_name"`
	CoordinatorAddr string `toml:"coordinator_addr" json:"coordinator_addr"`

	AdminAddr string `toml:"admin_addr" json:"admin_addr"`

	HostAdmin string `toml:"-" json:"-"`

	ProductName string `toml:"product_name" json:"product_name"`
	ProductAuth string `toml:"product_auth" json:"-"`

	ForwardMethod        string         `toml:"forward_method" json:"forward_method"`
	MigrateAsyncMaxBulks int            `toml:"migrate_async_maxbulks" json:"migrate_async_maxbulks"`
	MigrateAsyncMaxBytes bytesize.Int64 `toml:"migrate_async_maxbytes" json:"migrate_async_maxbytes"`
	MigrateAsyncPipeline int            `toml:"migrate_async_pipeline" json:"migrate_async_pipeline"`
	MigrateAsyncNumKeys  int            `toml:"migrate_async_numkeys" json:"migrate_async_numkeys"`

	SentinelQuorum               int               `toml:"sentinel_quorum" json:"sentinel_quorum"`
	SentinelParallelSyncs        int               `toml:"sentinel_parallel_syncs" json:"sentinel_parallel_syncs"`
	SentinelDownAfter            timesize.Duration `toml:"sentinel_down_after" json:"sentinel_down_after"`
	SentinelFailoverTimeout      timesize.Duration `toml:"sentinel_failover_timeout" json:"sentinel_failover_timeout"`
	SentinelNotificationScript   string            `toml:"sentinel_notification_script" json:"sentinel_notification_script"`
	SentinelClientReconfigScript string            `toml:"sentinel_client_reconfig_script" json:"sentinel_client_reconfig_script"`
}

func NewDefaultConfig() *Config {
	c := &Config{}
	if _, err := toml.Decode(DefaultConfig, c); err != nil {
		log.PanicErrorf(err, "decode toml failed")
	}
	if err := c.Validate(); err != nil {
		log.PanicErrorf(err, "validate config failed")
	}
	return c
}

func (c *Config) LoadFromFile(path string) error {
	_, err := toml.DecodeFile(path, c)
	if err != nil {
		return errors.Trace(err)
	}
	return c.Validate()
}

func (c *Config) String() string {
	var b bytes.Buffer
	e := toml.NewEncoder(&b)
	e.Indent = "    "
	e.Encode(c)
	return b.String()
}

func (c *Config) Validate() error {
	if c.CoordinatorName == "" {
		return errors.New("invalid coordinator_name")
	}
	if c.CoordinatorAddr == "" {
		return errors.New("invalid coordinator_addr")
	}
	if c.AdminAddr == "" {
		return errors.New("invalid admin_addr")
	}
	if c.ProductName == "" {
		return errors.New("invalid product_name")
	}
	if _, ok := models.ParseForwardMethod(c.ForwardMethod); !ok {
		return errors.New("invalid forward_method")
	}
	if c.MigrateAsyncMaxBulks <= 0 {
		return errors.New("invalid migrate_async_maxbulks")
	}
	if c.MigrateAsyncMaxBytes <= 0 {
		return errors.New("invalid migrate_async_maxbytes")
	}
	if c.MigrateAsyncPipeline <= 0 {
		return errors.New("invalid migrate_async_pipeline")
	}
	if c.MigrateAsyncNumKeys <= 0 {
		return errors.New("invalid migrate_async_numkeys")
	}
	if c.SentinelQuorum <= 0 {
		return errors.New("invalid sentinel_quorum")
	}
	if c.SentinelParallelSyncs <= 0 {
		return errors.New("invalid sentinel_parallel_syncs")
	}
	if c.SentinelDownAfter <= 0 {
		return errors.New("invalid sentinel_down_after")
	}
	if c.SentinelFailoverTimeout <= 0 {
		return errors.New("invalid sentinel_failover_timeout")
	}
	return nil
}
