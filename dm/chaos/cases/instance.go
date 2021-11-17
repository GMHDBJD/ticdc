// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"

	config2 "github.com/pingcap/ticdc/dm/dm/config"
	"github.com/pingcap/ticdc/dm/pkg/conn"
)

// set lesser sql_mode to tolerate some SQLs generated by go-sqlsmith.
var mustExecSQLs = []string{
	`SET @@GLOBAL.SQL_MODE="NO_ENGINE_SUBSTITUTION"`,
}

// setInstancesState sets the state (like global sql_mode) for upstream and downstream DB instances.
func setInstancesState(ctx context.Context, targetCfg *config2.DBConfig, sourcesCfg ...*config2.DBConfig) error {
	targetDB, err := conn.DefaultDBProvider.Apply(targetCfg)
	if err != nil {
		return err
	}
	for _, query := range mustExecSQLs {
		_, err = targetDB.DB.ExecContext(ctx, query)
		if err != nil {
			return err
		}
	}

	for _, cfg := range sourcesCfg {
		sourceDB, err2 := conn.DefaultDBProvider.Apply(cfg)
		if err2 != nil {
			return err2
		}
		for _, query := range mustExecSQLs {
			_, err2 = sourceDB.DB.ExecContext(ctx, query)
			if err2 != nil {
				return err2
			}
		}
	}

	return nil
}
