//
// Copyright 2017 SilencerCo, LLC
// Copyright 2015-2017 Pedro Salgado
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cfg

import (
	"fmt"
	"reflect"
)

// Config configuration interface.
type Config interface {
	// Value returns the value for the key associated with this configuration, otherwise nil.
	Value(key interface{}) interface{}
}

// emptyConfig empty implementation of configuration.
type emptyConfig uint8

// String returns the string representation of an emptyConfig context.
func (e *emptyConfig) String() string {
	switch e {
	case root:
		return "config.Root"
	}
	return "unknown config"
}

// Value returns the value for the key associated with this configuration, or nil
// if no value is associated with key.
func (*emptyConfig) Value(key interface{}) interface{} {
	return nil
}

// root
var root = new(emptyConfig)

// Root returns the configuration root.
func Root() Config {
	return root
}

// WithValue returns a new configuration with
// a new key/value pair added.
func WithValue(parent Config, k, v interface{}) Config {
	if k == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(k).Comparable() {
		panic("key is not comparable")
	}
	return &keyValuePair{parent, k, v}
}

// A keyValuePair carries a key-value pair.
// It implements Value for that key and delegates all other calls to the parent configuration.
type keyValuePair struct {
	Config
	key, val interface{}
}

// String returns a string representation of this struct.
func (c *keyValuePair) String() string {
	return fmt.Sprintf("%v.WithValue(%#v, %#v)", c.Config, c.key, c.val)
}

// Value returns the value for the key associated with this configuration, or nil
// if no value is associated with key (in the pair or in its parent configuration).
func (c *keyValuePair) Value(key interface{}) interface{} {
	if c.key == key {
		return c.val
	}

	// find key in parent
	return c.Config.Value(key)
}
