package main

func ConvertNullablePointerToString(pointer *string) string {
	if pointer != nil {
		return *pointer
	} else {
		return ""
	}
}
