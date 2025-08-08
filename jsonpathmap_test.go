package jsonpathmap_test

import (
	"testing"

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

	data, _ := jsonpathmap.UnMarshalJSON(jsonStr)
	pvs, _ := jsonpathmap.FlattenJSON(data)

	for _, pv := range pvs {
		t.Logf("%s = %v", pv.Path, pv.Value)
	}

	recovered, _ := jsonpathmap.UnflattenJSON(pvs)
	t.Log(jsonpathmap.MarshalJSON(recovered))

	if jsonpathmap.MarshalJSON(data) != jsonpathmap.MarshalJSON(recovered) {
		t.Error("flatten and unflatten mismatch")
	}
}
