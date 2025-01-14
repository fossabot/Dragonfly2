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

package atomiccount

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(AtomicCountUtilSuite))
}

type AtomicCountUtilSuite struct{
	suite.Suite
}

func (suite *AtomicCountUtilSuite) TestAdd() {
	acount := NewAtomicInt(3)
	acount.Add(4)
	acount.Add(5)

	result := acount.Get()
	suite.Equal(result, (int32)(12))

	var nilAcount *AtomicInt
	suite.Equal(nilAcount.Add(5), (int32)(0))
}

func (suite *AtomicCountUtilSuite) TestSet() {
	acount := NewAtomicInt(3)
	acount.Add(4)
	acount.Add(5)

	_ = acount.Set(1)
	result := acount.Get()
	suite.Equal(result, (int32)(1))

	var nilAcount *AtomicInt
	suite.Equal(nilAcount.Get(), (int32)(0))
}
