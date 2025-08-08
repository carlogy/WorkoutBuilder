package services

import (
	"encoding/json"
	"fmt"

	"github.com/sqlc-dev/pqtype"
)

func ConvertMapToRawJSON[T string | int](mapToConvert map[T]T) pqtype.NullRawMessage {
	if mapToConvert == nil {
		nullval := map[string]string{"null_val": "null"}
		jsonBytes, err := json.Marshal(nullval)
		if err != nil {
			fmt.Println(err)
		}
		return pqtype.NullRawMessage{RawMessage: jsonBytes, Valid: false}
	}

	jsonBytes, err := json.Marshal(mapToConvert)
	if err != nil {
		fmt.Printf("Error converting map to RawJSON: %v\n", err)
	}

	var NullRawMessage pqtype.NullRawMessage

	NullRawMessage.RawMessage = jsonBytes
	NullRawMessage.Valid = true

	return NullRawMessage
}

func ConvertMapsToRawJSON(data any) (pqtype.NullRawMessage, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return pqtype.NullRawMessage{}, err
	}

	return pqtype.NullRawMessage{
		RawMessage: jsonBytes,
		Valid:      true,
	}, nil
}

func ConvertRawJSONTOMap[T string | int](data pqtype.NullRawMessage) map[T]T {
	if !data.Valid {
		return nil
	}

	var mapData map[T]T
	err := json.Unmarshal(data.RawMessage, &mapData)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return mapData
}
