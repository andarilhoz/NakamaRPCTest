package payload

import "encoding/json"

type PayloadRequest struct {
	RequestType    string  `json:"type"`
	RequestVersion string  `json:"version"`
	RequestHash    *string `json:"hash"`
}

type Response struct {
	DataType    string  `json:"type"`
	DataVersion string  `json:"version"`
	DataHash    string  `json:"hash"`
	DataContent *string `json:"content"`
}

func DeserializePayload(payload string) (PayloadRequest, error) {
	var request PayloadRequest

	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		return PayloadRequest{}, err
	}

	request.PopulateDefaultValues()
	return request, nil
}

func GenerateResponse(request PayloadRequest, request_hash string, content string, equalHashes bool) (string, error) {
	var dataContent *string
	if equalHashes {
		dataContent = &content
	} else {
		dataContent = nil
	}

	responseObject := Response{
		DataType:    request.RequestType,
		DataVersion: request.RequestVersion,
		DataHash:    request_hash,
		DataContent: dataContent,
	}

	response, err := json.Marshal(responseObject)
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func (request *PayloadRequest) PopulateDefaultValues() {
	if request.RequestType == "" {
		request.RequestType = "core"
	}

	if request.RequestVersion == "" {
		request.RequestVersion = "1.0.0"
	}
}
