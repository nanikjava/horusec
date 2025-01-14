// Copyright 2021 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
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

package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLine(t *testing.T) {
	message := &Message{
		Line: 1,
	}

	t.Run("should success get line", func(t *testing.T) {
		line := message.GetLine()

		assert.NotEmpty(t, line)
		assert.Equal(t, "1", line)
	})
}

func TestGetColumn(t *testing.T) {
	message := &Message{
		Column: 1,
	}

	t.Run("should success get column", func(t *testing.T) {
		column := message.GetColumn()

		assert.NotEmpty(t, column)
		assert.Equal(t, "1", column)
	})
}

func TestIsValidMessage(t *testing.T) {
	t.Run("should return false if invalid message", func(t *testing.T) {
		message := &Message{
			Message: "This implies that some PHP code is not scanned by PHPCS",
			Type:    "ERROR",
		}

		assert.False(t, message.IsValidMessage())
	})

	t.Run("should return true if valid message", func(t *testing.T) {
		message := &Message{
			Message: "",
			Type:    "ERROR",
		}

		assert.True(t, message.IsValidMessage())
	})
}
