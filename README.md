# jsonpathma

一个 Go 库，用于在 JSON 与 `path->value` 格式之间进行相互转换，支持数组、嵌套对象，并可生成无数组索引的路径模板。

## 安装

```bash
go get github.com/suifengpiao14/jsonpathma
```
功能
1. FlattenJSON：JSON → []PathValue
2. UnflattenJSON：[]PathValue → JSON
3. NormalizeArrayPath：去除数组索引，替换为 []

## 示例
```go
package main

import (
	"fmt"
	"github.com/suifengpiao14/jsonpathmap"
)

func main() {
	jsonStr := `{"data":{"items":[{"name":"a"},{"name":"b"}]}}`
	data, _ := jsonpathmap.MarshalJSONStr(jsonStr)

	pvs, _ := jsonpathmap.FlattenJSON(data)
	fmt.Println(pvs)

	recovered, _ := jsonpathmap.UnflattenJSON(pvs)
	fmt.Println(jsonpathmap.ToJSONStr(recovered))

	fmt.Println(jsonpathmap.NormalizeArrayPath("data.items[0].name"))
}

```