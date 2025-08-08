package jsonpathmap

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type PathValue struct {
	Path  string `json:"path"`
	Value any    `json:"value"`
}

type PathValues []PathValue

// NormalizeArrayPath 将数组路径规范化，例如将 data.items[0].description 转换为 data.items[].description 用于从数据示例到文档格式
func (ps PathValues) NormalizeArrayPath() PathValues {
	var result PathValues
	re := regexp.MustCompile(`\[\d+\]`)
	for _, pv := range ps {
		path := re.ReplaceAllString(pv.Path, "[]")
		result = append(result, PathValue{Path: path, Value: pv.Value})
	}
	return result
}

// IndexArrayPath 将数组路径索引化为第一个元素，例如将 data.items[].description 转换为 data.items[0].description 用于从文档格式到数据示例格式
func (ps PathValues) IndexArrayPath() PathValues {
	var result PathValues
	for _, pv := range ps {
		path := strings.ReplaceAll(pv.Path, "[]", "[0]")
		result = append(result, PathValue{Path: path, Value: pv.Value})
	}
	return result
}

func (ps PathValues) Unqueue() PathValues {
	var result PathValues
	m := make(map[string]struct{})
	for _, pv := range ps {
		if _, ok := m[pv.Path]; !ok {
			result = append(result, pv)
			m[pv.Path] = struct{}{}
		}
	}
	return result
}

// FlattenJSON 将任意 JSON 数据拍平成 path->value 格式
func FlattenJSON(data any) (PathValues, error) {
	var result PathValues
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &data) //确保data 是 map[string]any 或者 []any 类型
	if err != nil {
		return nil, err
	}
	err = flattenValue(data, "", &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func flattenValue(data any, prefix string, result *PathValues) error {
	switch v := data.(type) {
	case map[string]any:
		for key, val := range v {
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}
			if err := flattenValue(val, path, result); err != nil {
				return err
			}
		}
	case []any:
		if len(v) == 0 {
			*result = append(*result, PathValue{Path: prefix + "[]", Value: nil})
			return nil
		}
		for i, val := range v {
			path := fmt.Sprintf("%s[%d]", prefix, i)
			if err := flattenValue(val, path, result); err != nil {
				return err
			}
		}
	default:
		*result = append(*result, PathValue{Path: prefix, Value: v})
	}
	return nil
}

// UnflattenJSON 将 path->value 恢复为 JSON
func UnflattenJSON(pvs PathValues) (map[string]any, error) {
	root := make(map[string]any)
	for _, pv := range pvs {
		if err := setValueByPath(root, pv.Path, pv.Value); err != nil {
			return nil, err
		}
	}
	return root, nil
}

func setValueByPath(root map[string]any, path string, value any) error {
	parts := parsePath(path)
	current := any(root)

	for i, part := range parts {
		isLast := i == len(parts)-1
		key, idx, isArray := parseArrayKey(part)

		switch container := current.(type) {
		case map[string]any:
			if isArray {
				arr, ok := container[key].([]any)
				if !ok {
					arr = []any{}
				}
				for len(arr) <= idx {
					arr = append(arr, nil)
				}
				if isLast {
					arr[idx] = value
				} else {
					if arr[idx] == nil {
						arr[idx] = make(map[string]any)
					}
					current = arr[idx]
				}
				container[key] = arr
			} else {
				if isLast {
					container[key] = value
				} else {
					if _, ok := container[key]; !ok {
						container[key] = make(map[string]any)
					}
					current = container[key]
				}
			}
		case []any:
			if idx >= len(container) {
				return fmt.Errorf("invalid path: index %d out of range for %v", idx, path)
			}
			if isLast {
				container[idx] = value
			} else {
				current = container[idx]
			}
		default:
			return fmt.Errorf("invalid container type for path %v", path)
		}
	}
	return nil
}

// parsePath 按 "." 分割 path，保留 [] 信息
func parsePath(path string) []string {
	parts := strings.Split(path, ".")
	return parts
}

// parseArrayKey 解析数组 key 和索引
func parseArrayKey(key string) (base string, idx int, isArray bool) {
	if strings.Contains(key, "[") && strings.HasSuffix(key, "]") {
		base = key[:strings.Index(key, "[")]
		indexStr := key[strings.Index(key, "[")+1 : len(key)-1]
		index, _ := strconv.Atoi(indexStr)
		return base, index, true
	}
	return key, 0, false
}

// UnMarshalJSON 工具函数：JSON字符串转 any
func UnMarshalJSON(jsonStr string) (any, error) {
	var v any
	err := json.Unmarshal([]byte(jsonStr), &v)
	return v, err
}

// MarshalJSON 工具函数：any 转 JSON 字符串
func MarshalJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
