// Copyright 2022 EMQ Technologies Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlgen

import (
	"fmt"
	"github.com/lf-edge/ekuiper/pkg/cast"
	"github.com/xo/dburl"
)

type SqlQueryGenerator interface {
	IndexValuer
	SqlQueryStatement() (string, error)
	UpdateMaxIndexValue(rows map[string]interface{})
}

type IndexValuer interface {
	SetIndexValue(interface{})
	GetIndexValue() interface{}
}

const DATETIME_TYPE = "DATETIME"

type InternalSqlQueryCfg struct {
	Table          string      `json:"table"`
	Limit          int         `json:"limit"`
	IndexField     string      `json:"indexField"`
	IndexValue     interface{} `json:"indexValue"`
	IndexFieldType string      `json:"indexFieldType"`
	DateTimeFormat string      `json:"dateTimeFormat"`
}

func (i *InternalSqlQueryCfg) SetIndexValue(val interface{}) {
	i.IndexValue = val
}

func (i *InternalSqlQueryCfg) GetIndexValue() interface{} {
	return i.IndexValue
}

type TemplateSqlQueryCfg struct {
	TemplateSql    string      `json:"templateSql"`
	IndexField     string      `json:"indexField"`
	IndexValue     interface{} `json:"indexValue"`
	IndexFieldType string      `json:"indexFieldType"`
	DateTimeFormat string      `json:"dateTimeFormat"`
}

func (i *TemplateSqlQueryCfg) SetIndexValue(val interface{}) {
	i.IndexValue = val
}

func (i *TemplateSqlQueryCfg) GetIndexValue() interface{} {
	return i.IndexValue
}

type sqlConfig struct {
	TemplateSqlQueryCfg *TemplateSqlQueryCfg `json:"templateSqlQueryCfg"`
	InternalSqlQueryCfg *InternalSqlQueryCfg `json:"internalSqlQueryCfg"`
}

func (cfg *sqlConfig) Init(props map[string]interface{}) error {
	err := cast.MapToStruct(props, &cfg)
	if err != nil {
		return fmt.Errorf("read properties %v fail with error: %v", props, err)
	}

	if cfg.TemplateSqlQueryCfg == nil && cfg.InternalSqlQueryCfg == nil {
		return fmt.Errorf("read properties %v fail with error: %v", props, err)
	}
	return nil
}

func GetQueryGenerator(u *dburl.URL, props map[string]interface{}) (SqlQueryGenerator, error) {
	cfg := &sqlConfig{}
	err := cfg.Init(props)
	if err != nil {
		return nil, err
	}

	if cfg.TemplateSqlQueryCfg != nil {
		ge, err := NewTemplateSqlQuery(cfg.TemplateSqlQueryCfg)
		if err != nil {
			return nil, err
		} else {
			return ge, nil
		}
	}

	switch u.Driver {
	case "sqlserver":
		return NewSqlServerQuery(cfg.InternalSqlQueryCfg), nil
	default:
		return NewCommonSqlQuery(cfg.InternalSqlQueryCfg), nil
	}
}
