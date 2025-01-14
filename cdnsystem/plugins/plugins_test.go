/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package plugins

import (
	"d7y.io/dragonfly/v2/cdnsystem/config"
	"fmt"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

func TestPluginsTestSuite(t *testing.T) {
	suite.Run(t, new(PluginsTestSuite))
}

type PluginsTestSuite struct {
	suite.Suite
	mgr Manager
}

func (s *PluginsTestSuite) SetUpSuite() {
	s.mgr = mgr
}

func (s *PluginsTestSuite) TearDownSuite() {
	mgr = s.mgr
}

func (s *PluginsTestSuite) TearDownTest() {
	mgr = s.mgr
}

func (s *PluginsTestSuite) TestSetManager() {
	tmp := &managerIml{}
	SetManager(tmp)
	s.Equal(mgr, tmp)
}

// -----------------------------------------------------------------------------

func (s *PluginsTestSuite) TestInitialize() {
	var testCase = func(cfg *config.Config, b Builder,
		pt config.PluginType, name string, hasPlugin bool, errMsg string) {
		SetManager(NewManager())
		RegisterPlugin(pt, name, b)
		err := Initialize(cfg)
		plugin := GetPlugin(pt, name)

		if errMsg != "" {
			s.NotNil(err)
			s.EqualError(err, ".*"+errMsg+".*")
			s.Nil(plugin)
		} else {
			s.Nil(err)
			if hasPlugin {
				s.Equal(plugin.Type(), pt)
				s.Equal(plugin.Name(), name)
			} else {
				s.Nil(plugin)
			}
		}
	}
	var testFunc = func(pt config.PluginType) {
		errMsg := "build error"
		name := "test"
		var createBuilder = func(err bool) Builder {
			return func(conf interface{}) (plugin Plugin, e error) {
				if err {
					return nil, fmt.Errorf(errMsg)
				}
				return &mockPlugin{pt, name}, nil
			}
		}
		var createConf = func(enabled bool) *config.Config {
			plugins := make(map[config.PluginType][]*config.PluginProperties)
			plugins[pt] = []*config.PluginProperties{{Name: name, Enable: enabled}}
			return &config.Config{Plugins: plugins}
		}
		testCase(createConf(false), createBuilder(false),
			pt, name, false, "")
		testCase(createConf(true), nil,
			pt, name, false, "cannot find builder")
		testCase(createConf(true), createBuilder(true),
			pt, name, false, errMsg)
		testCase(createConf(true), createBuilder(false),
			pt, name, true, "")
	}

	for _, pt := range config.PluginTypes {
		testFunc(pt)
	}
}

func (s *PluginsTestSuite) TestManagerIml_Builder() {
	var builder Builder = func(conf interface{}) (plugin Plugin, e error) {
		return nil, nil
	}
	manager := NewManager()

	var testFunc = func(pt config.PluginType, name string, b Builder, result bool) {
		manager.AddBuilder(pt, name, b)
		obj := manager.GetBuilder(pt, name)
		if result {
			s.NotNil(obj)
			objVal := reflect.ValueOf(obj)
			bVal := reflect.ValueOf(b)
			s.Equal(objVal.Pointer(), bVal.Pointer())
			manager.DeleteBuilder(pt, name)
		} else {
			s.Nil(obj)
		}
	}

	testFunc(config.PluginType("test"), "test", builder, false)
	for _, pt := range config.PluginTypes {
		testFunc(pt, "test", builder, true)
		testFunc(pt, "", nil, false)
		testFunc(pt, "", builder, false)
		testFunc(pt, "test", nil, false)
	}
}

func (s *PluginsTestSuite) TestManagerIml_Plugin() {
	manager := NewManager()

	var testFunc = func(p Plugin, result bool) {
		manager.AddPlugin(p)
		obj := manager.GetPlugin(p.Type(), p.Name())
		if result {
			s.NotNil(obj)
			s.Equal(obj, p)
			manager.DeletePlugin(p.Type(), p.Name())
		} else {
			s.Nil(obj)
		}
	}

	testFunc(&mockPlugin{"test", "test"}, false)
	for _, pt := range config.PluginTypes {
		testFunc(&mockPlugin{pt, "test"}, true)
		testFunc(&mockPlugin{pt, ""}, false)
	}
}

func (s *PluginsTestSuite) TestRepositoryIml() {
	type testCase struct {
		pt        config.PluginType
		name      string
		data      interface{}
		addResult bool
	}
	var createCase = func(validPlugin bool, name string, data interface{}, result bool) testCase {
		pt := config.StoragePlugin
		if !validPlugin {
			pt = config.PluginType("test-validPlugin")
		}
		return testCase{
			pt:        pt,
			name:      name,
			data:      data,
			addResult: result,
		}
	}
	var tc = func(valid bool, name string, data interface{}) testCase {
		return createCase(valid, name, data, true)
	}
	var fc = func(valid bool, name string, data interface{}) testCase {
		return createCase(valid, name, data, false)
	}
	var cases = []testCase{
		fc(true, "test", nil),
		fc(true, "", "data"),
		fc(false, "test", "data"),
		tc(true, "test", "data"),
	}

	repo := NewRepository()
	for _, v := range cases {
		repo.Add(v.pt, v.name, v.data)
		data := repo.Get(v.pt, v.name)
		if v.addResult {
			s.NotNil(data)
			s.Equal(data, v.data)
			repo.Delete(v.pt, v.name)
			data = repo.Get(v.pt, v.name)
			s.Nil(data)
		} else {
			s.Nil(data)
		}
	}
}

func (s *PluginsTestSuite) TestValidate() {
	type testCase struct {
		pt       config.PluginType
		name     string
		expected bool
	}
	var cases = []testCase{
		{config.PluginType("test"), "", false},
		{config.PluginType("test"), "test", false},
	}
	for _, pt := range config.PluginTypes {
		cases = append(cases,
			testCase{pt, "", false},
			testCase{pt, "test", true},
		)
	}
	for _, v := range cases {
		s.Equal(validate(v.pt, v.name), v.expected)
	}
}

// -----------------------------------------------------------------------------

type mockPlugin struct {
	pt   config.PluginType
	name string
}

func (m *mockPlugin) Type() config.PluginType {
	return m.pt
}

func (m *mockPlugin) Name() string {
	return m.name
}
