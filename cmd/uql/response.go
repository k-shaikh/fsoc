// Copyright 2022 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package uql

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/relvacode/iso8601"
)

// Response represents a parsed UQL response body
type Response struct {
	model    *Model
	dataSets map[string]*DataSet
	errors   []*Error
}

func (resp *Response) Model() *Model {
	return resp.model
}

func (resp *Response) Main() *DataSet {
	return resp.dataSets["d:main"]
}

func (resp *Response) DataSet(ref DataSetRef) *DataSet {
	return resp.dataSets[ref.Dataset]
}

func (resp *Response) HasErrors() bool {
	return len(resp.errors) > 0
}

func (resp *Response) Errors() []*Error {
	return resp.errors
}

// Model represents the structure of the response data
type Model struct {
	Name   string       `json:"name"`
	Fields []ModelField `json:"fields"`
}

// ModelField is a description of one column of one data set
type ModelField struct {
	Alias string `json:"alias"`
	Type  string `json:"type"`
	Form  string `json:"form"`
	Hints *Hint  `json:"hints"`
	Model *Model `json:"model"`
}

// DataSet holds the result data along with its name, structure (Model) and metadata
type DataSet struct {
	Name     string
	Model    *Model
	Metadata map[string]any
	Values   [][]any
}

// DataSetRef is a reference to another data set within the Response
type DataSetRef struct {
	JsonPath string `json:"$jsonPath"`
	Dataset  string `json:"$dataset"`
}

type Error struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (mf *ModelField) IsReference() bool {
	return mf.Form == "reference"
}

// Hint provides additional information about a ModelField such as the MELT kind, MELT field and type
type Hint struct {
	Kind  string `json:"kind"`
	Field string `json:"field"`
	Type  string `json:"type"`
}

func Errors(errors []*Error) error {
	var messages []string
	for _, e := range errors {
		messages = append(messages, fmt.Sprintf("%s: %s", e.Title, e.Detail))
	}
	return fmt.Errorf(strings.Join(messages, ", "))
}

type DataType interface {
	int | float64 | string | DataSetRef | bool | time.Time
}

type valueDeserializer[T DataType] func(json.RawMessage) (T, error)

var (
	longDeserializer valueDeserializer[int] = func(raw json.RawMessage) (int, error) {
		var value int
		if err := json.Unmarshal(raw, &value); err != nil {
			return value, err
		}
		return value, nil
	}
	doubleDeserializer valueDeserializer[float64] = func(raw json.RawMessage) (float64, error) {
		var value float64
		if err := json.Unmarshal(raw, &value); err != nil {
			return value, err
		}
		return value, nil
	}
	stringDeserializer valueDeserializer[string] = func(raw json.RawMessage) (string, error) {
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return value, err
		}
		return value, nil
	}
	booleanDeserializer valueDeserializer[bool] = func(raw json.RawMessage) (bool, error) {
		var value bool
		if err := json.Unmarshal(raw, &value); err != nil {
			return value, err
		}
		return value, nil
	}
	timestampDeserializer valueDeserializer[time.Time] = func(raw json.RawMessage) (time.Time, error) {
		var value iso8601.Time
		if err := json.Unmarshal(raw, &value); err != nil {
			return value.Time, err
		}
		return value.Time, nil
	}
	dataSetRefDeserializer valueDeserializer[DataSetRef] = func(raw json.RawMessage) (DataSetRef, error) {
		var value DataSetRef
		if err := json.Unmarshal(raw, &value); err != nil {
			return value, err
		}
		return value, nil
	}
)
