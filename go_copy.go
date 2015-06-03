package go_copy

import (
	"fmt"
	//"log"
	"reflect"
)

type CopyValue struct {
	CopyIfEmpty bool
}

func Copy(source interface{}, target interface{}, copyIfEmpty bool) {
	this := &CopyValue{CopyIfEmpty: copyIfEmpty}
	this.Copy(source, target)
}

func (this *CopyValue) Copy(source interface{}, target interface{}) {
	sourceTyp := reflect.TypeOf(source)
	targetTyp := reflect.TypeOf(target)

	if targetTyp.Kind() != reflect.Ptr {
		panic("Copy target must be pointer.")
	}

	sourceVal := reflect.ValueOf(source)
	targetVal := reflect.ValueOf(target)

	targetVal = this.copyValue(sourceVal, sourceTyp, targetVal, targetTyp)
}

func (this *CopyValue) copyValue(sourceVal reflect.Value, sourceTyp reflect.Type, targetVal reflect.Value, targetTyp reflect.Type) reflect.Value {
	sourceTyp = GetTypeElem(sourceTyp)
	targetTyp = GetTypeElem(targetTyp)

	if sourceTyp.Kind() != targetTyp.Kind() {
		panic("Unable to copy values from diffirent types to another.")
	}

	sourceVal = GetValueElem(sourceVal)
	targetVal = GetValueElem(targetVal)

	if !sourceVal.IsValid() {
		return targetVal
	}

	if !targetVal.IsValid() {
		targetVal = reflect.New(targetTyp).Elem()
	}

	if !targetVal.CanSet() {
		return targetVal
	}

	switch sourceTyp.Kind() {
	case reflect.Struct:
		targetVal = this.copyStruct(sourceVal, sourceTyp, targetVal, targetTyp)
	case reflect.Slice:
		targetVal = this.copySlice(sourceVal, sourceTyp, targetVal, targetTyp)
	case reflect.Map:
		targetVal = this.copyMap(sourceVal, sourceTyp, targetVal, targetTyp)
	case reflect.Chan:
		targetVal = this.copyChan(sourceVal, sourceTyp, targetVal, targetTyp)
	default:
		targetVal.Set(sourceVal)
	}

	return targetVal
}

func (this *CopyValue) copyChan(sourceVal reflect.Value, sourceTyp reflect.Type, targetVal reflect.Value, targetTyp reflect.Type) reflect.Value {
	if !sourceTyp.AssignableTo(targetTyp) {
		return targetVal
	}

	targetVal.Set(reflect.MakeChan(sourceTyp, sourceVal.Cap()))

	return targetVal
}

func (this *CopyValue) copyMap(sourceVal reflect.Value, sourceTyp reflect.Type, targetVal reflect.Value, targetTyp reflect.Type) reflect.Value {
	if !sourceTyp.AssignableTo(targetTyp) {
		return targetVal
	}

	if targetVal.Len() == 0 {
		targetVal.Set(reflect.MakeMap(sourceTyp))
	}

	keys := sourceVal.MapKeys()
	for _, key := range keys {
		val := this.duplicateValue(sourceVal.MapIndex(key))
		targetVal.SetMapIndex(key, val)
	}

	return targetVal
}

func (this *CopyValue) duplicateValue(val reflect.Value) reflect.Value {
	typ := val.Type()
	newVal := reflect.New(typ)

	result := this.copyValue(val, typ, newVal, typ)
	if typ.Kind() == reflect.Ptr {
		return result.Addr()
	}

	return result
}

func (this *CopyValue) copySlice(sourceVal reflect.Value, sourceTyp reflect.Type, targetVal reflect.Value, targetTyp reflect.Type) reflect.Value {
	if !sourceTyp.AssignableTo(targetTyp) {
		return targetVal
	}

	for i := 0; i < sourceVal.Len(); i++ {
		val := this.duplicateValue(sourceVal.Index(i))
		targetVal.Set(reflect.Append(targetVal, val))
	}

	return targetVal
}

func (this *CopyValue) copyStruct(sourceVal reflect.Value, sourceTyp reflect.Type, targetVal reflect.Value, targetTyp reflect.Type) reflect.Value {
	for i := 0; i < sourceTyp.NumField(); i++ {
		sourceField := sourceTyp.Field(i)
		sourceFieldVal := sourceVal.Field(i)

		targetField, exists := targetTyp.FieldByName(sourceField.Name)
		if !exists {
			continue
		}

		sourceFieldVal = GetValueElem(sourceFieldVal)
		//log.Println(sourceFieldVal, sourceFieldVal.IsValid(), sourceField.Name)
		if !sourceFieldVal.IsValid() {
			continue
		}

		targetFieldVal := targetVal.FieldByName(sourceField.Name)

		if !targetFieldVal.CanSet() {
			continue
		}

		sourceFieldElem := GetTypeElem(sourceField.Type)
		targetFieldElem := GetTypeElem(targetField.Type)

		targetFieldVal = this.copyValue(sourceFieldVal, sourceFieldElem, targetFieldVal, targetFieldElem)

		if targetField.Type.Kind() == reflect.Ptr {
			targetVal.Field(i).Set(targetFieldVal.Addr())
			continue
		}

		targetVal.Field(i).Set(targetFieldVal)
	}

	return targetVal
}

func GetTypeElem(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ
}

func GetValueElem(val reflect.Value) reflect.Value {
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}

type ValueException struct {
	Message string
	Kind    reflect.Kind
}

func NewValueException(kind reflect.Kind, message string) *ValueException {
	return &ValueException{
		Message: message,
		Kind:    kind,
	}
}

func (this *ValueException) String() string {
	return fmt.Sprintf("Message: %v: %v", this.Message, this.Kind.String())
}
