// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Package spec contains models for the elasticsearch-specification schema
// found at https://github.com/elastic/elasticsearch-specification/blob/main/output/schema/schema.json.
package spec

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Model struct {
	Types []TypeDefinition `json:"types,omitempty"`
}

type TypeDefinition struct {
	Value interface{}
}

func (t *TypeDefinition) Get() interface{} {
	return t.Value
}

func (t *TypeDefinition) UnmarshalJSON(b []byte) error {
	var objMap map[string]json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var kind string
	err = json.Unmarshal(objMap["kind"], &kind)
	if err != nil {
		return err
	}

	switch kind {
	case "interface":
		var o Interface
		err = json.Unmarshal(b, &o)
		if err != nil {
			return err
		}
		t.Value = o
	case "enum":
		var o Enum
		err = json.Unmarshal(b, &o)
		if err != nil {
			return err
		}
		t.Value = o
	case "type_alias":
		var o TypeAlias
		err = json.Unmarshal(b, &o)
		if err != nil {
			return err
		}
		t.Value = o
	case "request", "response":
		// Ignore
	default:
		return fmt.Errorf("unhandled type kind %q", kind)
	}

	return nil
}

type Interface struct {
	AttachedBehaviors []string     `json:"attachedBehaviors,omitempty"`
	Behaviors         []Inherits   `json:"behaviors,omitempty"`
	CodegenNames      []string     `json:"codegenNames,omitempty"`
	Deprecation       *Deprecation `json:"deprecation"`
	Description       string       `json:"description,omitempty"`
	DocID             string       `json:"docId,omitempty"`
	DocURL            string       `json:"docUrl,omitempty"`
	Generics          []TypeName   `json:"generics,omitempty"`
	Implements        []Inherits   `json:"implements,omitempty"`
	Inherits          Inherits     `json:"inherits"`
	Kind              string       `json:"kind,omitempty"`
	Properties        []Property   `json:"properties,omitempty"`
	ShortcutProperty  string       `json:"shortcutProperty,omitempty"`
	SpecLocation      string       `json:"specLocation,omitempty"`
	TypeName          TypeName     `json:"name"`
	VariantName       string       `json:"variantName,omitempty"`
	Variants          Variants     `json:"variants"`
}

type EnumMember struct {
	Availability *Availabilites `json:"availability,omitempty"`
	Deprecation  Deprecation    `json:"deprecation"`
	Description  string         `json:"description,omitempty"`
	Identifier   string         `json:"identifier,omitempty"`
	Name         string         `json:"name,omitempty"`
	Since        string         `json:"since,omitempty"`
}

type Enum struct {
	CodegenNames []string     `json:"codegenNames,omitempty"`
	Deprecation  *Deprecation `json:"deprecation"`
	Description  string       `json:"description,omitempty"`
	DocID        string       `json:"docId,omitempty"`
	DocURL       string       `json:"docUrl,omitempty"`
	Kind         string       `json:"kind,omitempty"`
	Members      []EnumMember `json:"members,omitempty"`
	SpecLocation string       `json:"specLocation,omitempty"`
	TypeName     TypeName     `json:"name"`
	VariantName  string       `json:"variantName,omitempty"`
}

type TypeAlias struct {
	CodegenNames []string     `json:"codegenNames,omitempty"`
	Deprecation  *Deprecation `json:"deprecation"`
	Description  string       `json:"description,omitempty"`
	DocID        string       `json:"docId,omitempty"`
	DocURL       string       `json:"docUrl,omitempty"`
	Generics     []TypeName   `json:"generics,omitempty"`
	Kind         string       `json:"kind,omitempty"`
	SpecLocation string       `json:"specLocation,omitempty"`
	Type         *ValueOf     `json:"type,omitempty"`
	TypeName     TypeName     `json:"name"`
	VariantName  string       `json:"variantName,omitempty"`
	Variants     *Variants    `json:"variants,omitempty"`
}

type TypeName struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type Inherits struct {
	TypeName `json:"type"`
	Generics []ValueOf `json:"generics,omitempty"`
}

type Deprecation struct {
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
}

type Availability struct {
	FeatureFlag *string `json:"featureFlag,omitempty"`
	Since       *string `json:"since,omitempty"`
}

type Availabilites struct {
	Serverless *Availability `json:"serverless,omitempty"`
	Stack      *Availability `json:"stack,omitempty"`
}

// ValueOf need to be referenced in other structs as a pointer
// since valueOf wraps around itself for complex nested typed
// only pointers support recursive json parsing.
type ValueOf struct {
	*InstanceOf
	*ArrayOf
	*UnionOf
	*DictionaryOf
	*UserDefinedValue
	*LiteralValue
}

func (v *ValueOf) Get() interface{} {
	value := reflect.ValueOf(v)
	ref := reflect.TypeOf(v)
	for i := 0; i < ref.NumField(); i++ {
		if !value.Field(i).IsNil() {
			return value.Field(i).Interface()
		}
	}

	return nil
}

func (v *ValueOf) Kind() string {
	if v.InstanceOf != nil {
		return v.InstanceOf.Kind
	}
	if v.ArrayOf != nil {
		return v.ArrayOf.Kind
	}
	if v.UnionOf != nil {
		return v.UnionOf.Kind
	}
	if v.DictionaryOf != nil {
		return v.DictionaryOf.Kind
	}
	if v.UserDefinedValue != nil {
		return v.UserDefinedValue.Kind
	}
	if v.LiteralValue != nil {
		return v.LiteralValue.Kind
	}

	return ""
}

func (v *ValueOf) UnmarshalJSON(b []byte) error {
	var objMap map[string]json.RawMessage
	var arrObjMap []map[string]json.RawMessage

	err := json.Unmarshal(b, &objMap)
	if err != nil {
		err = json.Unmarshal(b, &arrObjMap)
		if err != nil {
			return err
		}
		// Skip arrays. Will deserialize item per item later on.
		return nil
	}

	var kind string
	err = json.Unmarshal(objMap["kind"], &kind)
	if err != nil {
		return err
	}

	switch kind {
	case "instance_of":
		var o InstanceOf
		err = json.Unmarshal(b, &o)
		if err != nil {
			return err
		}
		*v = ValueOf{InstanceOf: &o}
	case "array_of":
		var o ArrayOf
		err = json.Unmarshal(b, &o)
		if err != nil {
			return err
		}
		*v = ValueOf{ArrayOf: &o}
	case "union_of":
		var o UnionOf
		err = json.Unmarshal(b, &o)
		if err != nil {
			return err
		}
		*v = ValueOf{UnionOf: &o}
	case "dictionary_of":
		var o DictionaryOf
		err = json.Unmarshal(b, &o)
		if err != nil {
			return err
		}
		*v = ValueOf{DictionaryOf: &o}
	case "user_defined_value":
		var o UserDefinedValue
		err = json.Unmarshal(b, &o)
		if err != nil {
			return err
		}
		*v = ValueOf{UserDefinedValue: &o}
	case "literal_value":
		var o LiteralValue
		err = json.Unmarshal(b, &o)
		if err != nil {
			return err
		}
		*v = ValueOf{LiteralValue: &o}
	default:
		return fmt.Errorf("unhandled value kind %q", kind)
	}

	return nil
}

type InstanceOf struct {
	Kind     string `json:"kind"`
	TypeName `json:"type"`
	Generics *ValueOf `json:"generics"`
}

func (i InstanceOf) String() string {
	return fmt.Sprintf("InstanceOf || Kind: %s | TypeName: %v | Generics: %v", i.Kind, i.TypeName, i.Generics)
}

type ArrayOf struct {
	Kind  string   `json:"kind"`
	Value *ValueOf `json:"Value"`
}

func (a ArrayOf) String() string {
	return fmt.Sprintf("ArrayOf || Kind: %s | Value: %v", a.Kind, a.Value)
}

type UnionOf struct {
	Kind  string     `json:"kind"`
	Items []*ValueOf `json:"items"`
}

func (u UnionOf) String() string {
	return fmt.Sprintf("UnionOf || Kind: %s | Items: %v", u.Kind, u.Items)
}

type DictionaryOf struct {
	Kind      string   `json:"kind"`
	Key       *ValueOf `json:"key"`
	Value     *ValueOf `json:"Value"`
	SingleKey bool     `json:"singleKey"`
}

func (d DictionaryOf) String() string {
	return fmt.Sprintf("DictionaryOf || Kind: %s | Key: %v | Value: %v | SingleKey: %t", d.Kind, d.Key, d.Value, d.SingleKey)
}

type UserDefinedValue struct {
	Kind string `json:"kind"`
}

func (u UserDefinedValue) String() string {
	return fmt.Sprintf("UserDefinedValue || Kind: %s", u.Kind)
}

type LiteralValue struct {
	Kind  string      `json:"kind"`
	Value interface{} `json:"Value"`
}

func (l LiteralValue) String() string {
	return fmt.Sprintf("LiteralValue || Kind: %s, Value: %v", l.Kind, l.Value)
}

type Variants struct {
	DefaultTag    string `json:"defaultTag,omitempty"`
	Kind          string `json:"kind,omitempty"`
	NonExhaustive bool   `json:"nonExhaustive,omitempty"`
	Tag           string `json:"tag,omitempty"`
}

type Property struct {
	Aliases           []string      `json:"aliases,omitempty"`
	Availability      *Availability `json:"availability,omitempty"`
	CodegenName       *string       `json:"codegenName,omitempty"`
	ContainerProperty *bool         `json:"containerProperty,omitempty"`
	Depreciation      *Deprecation  `json:"depreciation"`
	Description       *string       `json:"description,omitempty"`
	DocURL            *string       `json:"docUrl,omitempty"`
	Identifier        *string       `json:"identifier,omitempty"`
	Name              string        `json:"name,omitempty"`
	Required          *bool         `json:"required,omitempty"`
	ServerDefault     interface{}   `json:"serverDefault,omitempty"`
	Since             *string       `json:"since,omitempty"`
	Type              *ValueOf      `json:"type,omitempty"`
}
