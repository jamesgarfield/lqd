package lqdsql

import "testing"
import "reflect"

func Test_buildUpdate(t *testing.T) {
	tests := []struct {
		table  string
		values map[string]interface{}
		expect SQL
	}{
		{
			"users",
			map[string]interface{}{"name": "bob"},
			SQL{"UPDATE \"users\"\nSET\n\"name\" = $1", []interface{}{"bob"}},
		},
		{
			"items",
			map[string]interface{}{"type": "to do", "status": 5},
			SQL{"UPDATE \"items\"\nSET\n\"type\" = $1,\n\"status\" = $2", []interface{}{"to do", 5}},
		},
	}

	for i, x := range tests {
		result, err := buildUpdate(x.table, x.values)
		if err != nil {
			t.Errorf("unexpected error at index %d, %+v", i, err)
			continue
		}
		if !reflect.DeepEqual(result, x.expect) {
			t.Errorf("expected result at index %d to be %+v, got %+v", i, x.expect, result)
		}
	}
}
