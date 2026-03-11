package uncsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

type FieldDecoder interface {
	DecodeCSV(string) error
}

type FieldEncoder interface {
	EncodeCSV() (string, error)
}

type Decoder struct {
	r     io.Reader
	Comma rune
}

func (d *Decoder) SetComma(comma rune) {
	d.Comma = comma
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

type Encoder struct {
	w io.Writer
}

func decodeStructOfArrays(v any, colname_to_field map[string]int, r *csv.Reader) error {
	destT := reflect.TypeOf(v) // ignore possible nil as it was already parsed successfully once
	destV := reflect.ValueOf(v)

	if destT.Kind() == reflect.Pointer {
		destT = destT.Elem()
		destV = destV.Elem()
	}

	for rowIdx := 0; ; rowIdx++ {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("row parsing error: %w", err)
		}

		for i := range destT.NumField() {
			fieldT := destT.Field(i)

			fieldKind := fieldT.Type.Kind()
			if fieldKind != reflect.Slice && fieldKind != reflect.Array {
				return fmt.Errorf("expected fields to be arrays|slices got %s", fieldT.Type.Kind())
			}

			fieldV := destV.Field(i)
			if fieldV.Kind() == reflect.Slice && fieldV.Cap() == 0 {
				fieldV.Set(reflect.MakeSlice(fieldV.Type(), 0, 128))
			}

			tag := fieldT.Tag.Get("csv")
			if tag == "" {
				continue
			}

			colIdx, ok := colname_to_field[tag]
			if !ok {
				return fmt.Errorf("field %s: column %s not found in CSV header", fieldT.Name, tag)
			}

			elemType := fieldT.Type.Elem()
			// Optimistic that type has DecodeCSV()
			val, err := parseElement(row[colIdx], elemType)
			if err != nil {
				return fmt.Errorf("at field %s, colIdx %d: element parsing error: %w", fieldT.Name, colIdx, err)
			}

			if err := setValueAtIndex(fieldV, rowIdx, val); err != nil {
				return fmt.Errorf("failed to set field: %w", err)
			}
		}
	}

	return nil
}

func decodeArrayOfStructs(v any, colname_to_field map[string]int, r *csv.Reader) error {
	destT := reflect.TypeOf(v)
	destV := reflect.ValueOf(v)

	// Dereference pointer (v is always *[]T or *[N]T from Decode)
	destT = destT.Elem()
	destV = destV.Elem()

	elemType := destT.Elem() // the struct element type
	if elemType.Kind() != reflect.Struct {
		return fmt.Errorf("expected slice/array of structs, got slice/array of %s", elemType.Kind())
	}

	for rowIdx := 0; ; rowIdx++ {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("decoding array of structs error: %w", err)
		}

		newElem := reflect.New(elemType).Elem()

		for i := range elemType.NumField() {
			fieldT := elemType.Field(i)
			tag := fieldT.Tag.Get("csv")
			if tag == "" {
				continue
			}

			colIdx, ok := colname_to_field[tag]
			if !ok {
				return fmt.Errorf("field %s: column %s not found in CSV header", fieldT.Name, tag)
			}

			val, err := parseElement(row[colIdx], fieldT.Type)
			if err != nil {
				return fmt.Errorf("at field %s, colIdx %d: element parsing error: %w", fieldT.Name, colIdx, err)
			}
			newElem.Field(i).Set(reflect.ValueOf(val).Convert(fieldT.Type))
		}

		if err := setValueAtIndex(destV, rowIdx, newElem.Interface()); err != nil {
			return fmt.Errorf("failed to set element at row %d: %w", rowIdx, err)
		}
	}

	return nil
}

// Assumes struct of arrays
func (p *Decoder) Decode(v any) error {
	csvReader := csv.NewReader(p.r)
	csvReader.ReuseRecord = true
	if p.Comma != 0 {
		csvReader.Comma = p.Comma
	}

	header, err := csvReader.Read()
	if err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	// Strip BOM from first header column if present
	if len(header) > 0 {
		header[0] = strings.TrimPrefix(header[0], "\uFEFF")
	}

	columnNameToField := make(map[string]int)
	for i, name := range header {
		columnNameToField[name] = i
	}

	destT := reflect.TypeOf(v)
	if destT == nil {
		return fmt.Errorf("cannot decode nil")
	}

	if destT.Kind() != reflect.Pointer {
		return fmt.Errorf("expected to recieve pointer to dest recieved: %s", destT.String())
	}
	destT = destT.Elem()

	// both decide only outermost shell
	if destT.Kind() == reflect.Slice || destT.Kind() == reflect.Array {
		return decodeArrayOfStructs(v, columnNameToField, csvReader)
	} else if destT.Kind() == reflect.Struct {
		return decodeStructOfArrays(v, columnNameToField, csvReader)
	} else {
		return fmt.Errorf("well, you can put any, BUT HOW DO YOU EXPECT TO FIT CSV IN %s", destT.String())
	}
}

func (p *Encoder) Encode(v any) error {
	return errors.New("not immplemented")
}

func parseElement(strval string, elemType reflect.Type) (any, error) {
	if reflect.PointerTo(elemType).Implements(reflect.TypeFor[FieldDecoder]()) {
		newElem := reflect.New(elemType)
		if decoder, ok := newElem.Interface().(FieldDecoder); ok {
			if err := decoder.DecodeCSV(strval); err != nil {
				return nil, fmt.Errorf("custom decode failed: %w", err)
			}

			return newElem.Elem().Interface(), nil
		}
	}
	elemKind := elemType.Kind()

	elemKindSizeBits := getBitSizeFromKind(elemKind)
	switch elemKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(strval, 0, elemKindSizeBits)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse value %s as %d-bit integer: %w",
				strval, elemKindSizeBits, err,
			)
		}
		return intVal, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(strval, 0, elemKindSizeBits)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse value %s as %d-bit unsigned integer: %w",
				strval, elemKindSizeBits, err,
			)
		}
		return uintVal, nil
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(strval)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse value %s as a boolean: %w",
				strval, err,
			)
		}
		return boolVal, nil
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(strval, elemKindSizeBits)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse value %s as float%d: %w",
				strval, elemKindSizeBits, err,
			)
		}
		return floatVal, nil
	case reflect.String:
		return strval, nil
	default:
		return nil, fmt.Errorf("expected kind of element to be of simple type")
	}
}

// Rename and return bits directly
func getBitSizeFromKind(kind reflect.Kind) int {
	switch kind {
	case reflect.Int8, reflect.Uint8:
		return 8
	case reflect.Int16, reflect.Uint16:
		return 16
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return 32
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		return 64
	case reflect.Int, reflect.Uint:
		return int(unsafe.Sizeof(int(0)) * 8)
	default:
		return 0
	}
}

func setValueAtIndex(fieldValue reflect.Value, index int, value any) error {
	kind := fieldValue.Kind()

	switch kind {
	case reflect.Array:
		if index >= fieldValue.Len() {
			return fmt.Errorf("index %d out of bounds for array of length %d", index, fieldValue.Len())
		}
	case reflect.Slice:
		if index >= fieldValue.Cap() {
			newSlice := reflect.MakeSlice(fieldValue.Type(), index+1, (index+1)*2)
			reflect.Copy(newSlice, fieldValue)
			fieldValue.Set(newSlice)
		} else if index >= fieldValue.Len() {
			fieldValue.SetLen(index + 1)
		}
	default:
		return fmt.Errorf("cannot set index on kind %v", kind)
	}

	elem := fieldValue.Index(index)

	switch v := value.(type) {
	case int64:
		elem.SetInt(v)
	case uint64:
		elem.SetUint(v)
	case float64:
		elem.SetFloat(v)
	case bool:
		elem.SetBool(v)
	case string:
		elem.SetString(v)
	default:
		elem.Set(reflect.ValueOf(value))
	}

	return nil
}
