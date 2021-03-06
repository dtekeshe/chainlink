package adapters_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestMultiply_Perform(t *testing.T) {
	tests := []struct {
		name      string
		params    string
		json      string
		want      string
		errored   bool
		jsonError bool
	}{
		{"string", `{"times":100}`, `{"value":"1.23"}`, "123", false, false},
		{"integer", `{"times":100}`, `{"value":123}`, "12300", false, false},
		{"float", `{"times":100}`, `{"value":1.23}`, "123", false, false},
		{"object", `{"times":100}`, `{"value":{"foo":"bar"}}`, "", true, false},
		{"zero integer string", `{"times":0}`, `{"value":"1.23"}`, "0", false, false},
		{"negative integer string", `{"times":-5}`, `{"value":"1.23"}`, "-6.15", false, false},

		{"string string", `{"times":"100"}`, `{"value":"1.23"}`, "123", false, false},
		{"string integer", `{"times":"100"}`, `{"value":123}`, "12300", false, false},
		{"string float", `{"times":"100"}`, `{"value":1.23}`, "123", false, false},
		{"string object", `{"times":"100"}`, `{"value":{"foo":"bar"}}`, "", true, false},
		{"array string", `{"times":[1, 2, 3]}`, `{"value":"1.23"}`, "", false, true},
		{"rubbish string", `{"times":"123aaa123"}`, `{"value":"1.23"}`, "", false, true},
		{"zero string string", `{"times":"0"}`, `{"value":"1.23"}`, "0", false, false},
		{"negative string string", `{"times":"-5"}`, `{"value":"1.23"}`, "-6.15", false, false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := models.RunResult{
				Data: cltest.JSONFromString(test.json),
			}
			adapter := adapters.Multiply{}
			jsonErr := json.Unmarshal([]byte(test.params), &adapter)
			result := adapter.Perform(input, nil)

			if test.jsonError {
				assert.Error(t, jsonErr)
			} else if test.errored {
				assert.Error(t, result.GetError())
				assert.NoError(t, jsonErr)
			} else {
				val, err := result.Value()
				assert.NoError(t, err)
				assert.Equal(t, test.want, val)
				assert.NoError(t, result.GetError())
				assert.NoError(t, jsonErr)
			}
		})
	}
}
