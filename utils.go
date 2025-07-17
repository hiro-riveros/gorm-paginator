package paginator

import "reflect"

// makeSlice generates a new slice of the model type: []*Model
func makeSlice(model interface{}) interface{} {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	sliceType := reflect.SliceOf(reflect.PtrTo(modelType))
	return reflect.New(sliceType).Interface()
}

// dereferenceSlice converts *slice to slice (for example: *([]*User) → []*User)
func dereferenceSlice(ptr interface{}) interface{} {
	return reflect.ValueOf(ptr).Elem().Interface()
}
