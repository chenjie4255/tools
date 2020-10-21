package structmap

import (
	"encoding/json"
)

func Convert2JSONMap(o interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	ret := map[string]interface{}{}
	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}
