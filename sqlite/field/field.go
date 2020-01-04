/*
 * BSD 2-Clause License
 *
 *	Copyright (c) 2019, Piotr Pszczółkowski
 *	All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 * list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 * this list of conditions and the following disclaimer in the documentation
 * and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package field

import (
	"errors"
	"strings"

	"Timelancer/sqlite/vtc"
)

type Field struct {
	Name      string        `json:"name"`
	Value     interface{}   `json:"value"`
	ValueType vtc.ValueType `json:"type"`
}

func New(name string) *Field {
	if name := strings.TrimSpace(name); name != "" {
		return &Field{Name: name}
	}
	return nil
}

func NewWithValue(name string, v interface{}) *Field {
	if f := New(name); f != nil {
		return f.SetValue(v)
	}
	return nil
}

func (f *Field) Valid() bool {
	name := strings.TrimSpace(f.Name)
	if name == "" || name != f.Name {
		return false
	}
	switch f.ValueType {
	case vtc.Text, vtc.Int, vtc.Float, vtc.Blob, vtc.Null:
		return true
	default:
		return false
	}
}

func (f *Field) SetValue(v interface{}) *Field {
	f.Value = v

	switch v.(type) {
	case string:
		f.ValueType = vtc.Text
	case int64:
		f.ValueType = vtc.Int
	case float32:
		f.Value = float64(v.(float32))
		f.ValueType = vtc.Float
	case float64:
		f.ValueType = vtc.Float
	case []byte:
		f.ValueType = vtc.Blob
	case bool:
		v := v.(bool)
		if v == true {
			f.Value = int64(1)
		} else {
			f.Value = int64(0)
		}
		f.ValueType = vtc.Int
	default:
		f.ValueType = vtc.Null
	}
	return f
}

func (f *Field) Int64() (int, error) {
	if v, ok := f.Value.(int64); ok {
		return int(v), nil
	}
	return 0, errors.New("can't convert field value to int64")
}

func (f *Field) UInt64() (uint64, error) {
	if v, ok := f.Value.(int64); ok {
		return uint64(v), nil
	}
	return 0, errors.New("can't convert field value to uint64")
}

func (f *Field) Int32() (int32, error) {
	if v, ok := f.Value.(int64); ok {
		return int32(v), nil
	}
	return 0, errors.New("can't convert field value to int32")
}

func (f *Field) UInt32() (uint32, error) {
	if v, ok := f.Value.(int64); ok {
		return uint32(v), nil
	}
	return 0, errors.New("can't convert field value to uint32")
}

func (f *Field) Float32() (float32, error) {
	if v, ok := f.Value.(float64); ok {
		return float32(v), nil
	}
	return 0.0, errors.New("can't convert field value to float32")
}

func (f *Field) Float64() (float64, error) {
	if v, ok := f.Value.(float64); ok {
		return v, nil
	}
	return 0.0, errors.New("can't convert field value to float64")
}

func (f *Field) Text() (string, error) {
	if v, ok := f.Value.(string); ok {
		return v, nil
	}
	return "", errors.New("can't convert field value to string")
}

func (f *Field) Blob() ([]byte, error) {
	if v, ok := f.Value.([]byte); ok {
		return v, nil
	}
	return nil, errors.New("can't convert field value to []byte")
}

func (f *Field) Bool() (bool, error) {
	v, err := f.Int64()
	if err != nil {
		return false, err
	}

	if v == 1 {
		return true, nil
	}
	return false, nil
}
