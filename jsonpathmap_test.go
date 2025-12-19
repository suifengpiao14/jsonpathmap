package jsonpathmap_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/jsonpathmap"
)

func TestFlattenAndUnflatten(t *testing.T) {
	jsonStr := `{
		"code": "",
		"message": "",
		"data": {
			"items": [
				{"description": "a", "dbComment": "c1", "regexDescription": "r1"},
				{"description": "b", "dbComment": "c2", "regexDescription": "r2"}
			],
			"pageIndex": "1",
			"pageSize": "10",
			"total": "2"
		}
	}`
	pvs, _ := jsonpathmap.FlattenJSON([]byte(jsonStr))

	for _, pv := range pvs {
		t.Logf("%s = %v", pv.Path, pv.Value)
	}

	recovered, _ := jsonpathmap.UnflattenJSON(pvs)
	b, err := json.MarshalIndent(recovered, "", "  ")
	require.NoError(t, err)
	require.JSONEq(t, jsonStr, string(b))

}
