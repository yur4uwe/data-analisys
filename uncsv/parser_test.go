package uncsv

import (
	"strconv"
	"strings"
	"testing"
)

// assertSlicesEqual checks that two slices have identical length and elements.
func assertSlicesEqual[T comparable](t *testing.T, name string, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s: length mismatch: got %d, want %d", name, len(got), len(want))
		return
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("%s[%d]: got %v, want %v", name, i, got[i], want[i])
		}
	}
}

// TestNewDecoder tests the NewDecoder constructor
func TestNewDecoder(t *testing.T) {
	csvData := "header1,header2\nvalue1,value2"
	reader := strings.NewReader(csvData)

	decoder := NewDecoder(reader)

	if decoder == nil {
		t.Fatal("NewDecoder returned nil")
	}
}

// TestDecodeEmptyData tests decoding with empty CSV data
func TestDecodeEmptyData(t *testing.T) {
	csvData := ""
	reader := strings.NewReader(csvData)
	decoder := NewDecoder(reader)

	var result struct {
		Field1 []string
		Field2 []int
	}

	err := decoder.Decode(&result)
	// Should handle empty data gracefully
	if err != nil {
		t.Logf("Decode with empty data returned error: %v", err)
	}

	assertSlicesEqual(t, "Field1", result.Field1, []string{})
	assertSlicesEqual(t, "Field2", result.Field2, []int{})
}

// TestDecodeSimpleStruct tests decoding into a simple struct of arrays
func TestDecodeSimpleStruct(t *testing.T) {
	csvData := "name,age\nAlice,30\nBob,25\nCharlie,35"
	reader := strings.NewReader(csvData)
	decoder := NewDecoder(reader)

	var result struct {
		Name []string `csv:"name"`
		Age  []int    `csv:"age"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Errorf("Decode failed: %v", err)
	}

	assertSlicesEqual(t, "Name", result.Name, []string{"Alice", "Bob", "Charlie"})
	assertSlicesEqual(t, "Age", result.Age, []int{30, 25, 35})
}

// TestDecodeWithFloats tests decoding struct with float arrays
func TestDecodeWithFloats(t *testing.T) {
	csvData := "id,value,price\n1,10.5,20.99\n2,15.3,30.50\n3,12.7,25.75"
	reader := strings.NewReader(csvData)
	decoder := NewDecoder(reader)

	var result struct {
		ID    []int     `csv:"id"`
		Value []float64 `csv:"value"`
		Price []float64 `csv:"price"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Errorf("Decode with floats failed: %v", err)
	}

	assertSlicesEqual(t, "ID", result.ID, []int{1, 2, 3})
	assertSlicesEqual(t, "Value", result.Value, []float64{10.5, 15.3, 12.7})
	assertSlicesEqual(t, "Price", result.Price, []float64{20.99, 30.50, 25.75})
}

// TestDecodeSingleRow tests decoding a single row
func TestDecodeSingleRow(t *testing.T) {
	csvData := "name,status\nJohn,active"
	reader := strings.NewReader(csvData)
	decoder := NewDecoder(reader)

	var result struct {
		Name   []string `csv:"name"`
		Status []string `csv:"status"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Errorf("Decode single row failed: %v", err)
	}

	assertSlicesEqual(t, "Name", result.Name, []string{"John"})
	assertSlicesEqual(t, "Status", result.Status, []string{"active"})
}

// TestDecodeMultipleFields tests decoding with many fields
func TestDecodeMultipleFields(t *testing.T) {
	csvData := "a,b,c,d,e\n1,2,3,4,5\n6,7,8,9,10"
	reader := strings.NewReader(csvData)
	decoder := NewDecoder(reader)

	var result struct {
		A []int `csv:"a"`
		B []int `csv:"b"`
		C []int `csv:"c"`
		D []int `csv:"d"`
		E []int `csv:"e"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Errorf("Decode multiple fields failed: %v", err)
	}

	assertSlicesEqual(t, "A", result.A, []int{1, 6})
	assertSlicesEqual(t, "B", result.B, []int{2, 7})
	assertSlicesEqual(t, "C", result.C, []int{3, 8})
	assertSlicesEqual(t, "D", result.D, []int{4, 9})
	assertSlicesEqual(t, "E", result.E, []int{5, 10})
}

// TestDecodeWithMissingValues tests decoding with missing/empty values
func TestDecodeWithMissingValues(t *testing.T) {
	csvData := "name,value\nAlice,100\nBob,\nCharlie,200"
	reader := strings.NewReader(csvData)
	decoder := NewDecoder(reader)

	var result struct {
		Name  []string `csv:"name"`
		Value []string `csv:"value"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Errorf("Decode with missing values failed: %v", err)
	}

	assertSlicesEqual(t, "Name", result.Name, []string{"Alice", "Bob", "Charlie"})
	assertSlicesEqual(t, "Value", result.Value, []string{"100", "", "200"})
}

// TestDecodeIntoNil tests Decode with nil pointer
func TestDecodeIntoNil(t *testing.T) {
	csvData := "a,b\n1,2"
	reader := strings.NewReader(csvData)
	decoder := NewDecoder(reader)

	err := decoder.Decode(nil)
	if err == nil {
		t.Error("Decode(nil) should return an error")
	}
}

// TestDecodeIntoNonStruct tests Decode with non-struct type
func TestDecodeIntoNonStruct(t *testing.T) {
	csvData := "a,b\n1,2"
	reader := strings.NewReader(csvData)
	decoder := NewDecoder(reader)

	var result int
	err := decoder.Decode(&result)
	if err == nil {
		t.Error("Decode into non-struct should return an error")
	}
}

// TestDecodeIntoNonArrayFields tests Decode where struct fields are not arrays
func TestDecodeIntoNonArrayFields(t *testing.T) {
	csvData := "name,age\nAlice,30"
	reader := strings.NewReader(csvData)
	decoder := NewDecoder(reader)

	var result struct {
		Name string `csv:"name"` // Not an array
		Age  int    `csv:"age"`  // Not an array
	}

	err := decoder.Decode(&result)
	if err == nil {
		t.Error("Decode into non-array fields should return an error")
	}
}

// TestDecodeFieldDecoderInterface tests FieldDecoder interface implementation
func TestDecodeFieldDecoderInterface(t *testing.T) {
	csvData := "value\n100\n200\n300"
	reader := strings.NewReader(csvData)
	decoder := NewDecoder(reader)

	var result struct {
		Value []CustomType `csv:"value"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Errorf("Decode with FieldDecoder implementation failed: %v", err)
	}

	if len(result.Value) != 3 {
		t.Fatalf("Value: length mismatch: got %d, want 3", len(result.Value))
	}
	want := []int{100, 200, 300}
	for i, v := range result.Value {
		if v.Value != want[i] {
			t.Errorf("Value[%d]: got %d, want %d", i, v.Value, want[i])
		}
	}
}

// CustomType implements FieldDecoder for testing
type CustomType struct {
	Value int
}

func (c *CustomType) DecodeCSV(s string) error {
	// Simple implementation for testing
	val, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	c.Value = val
	return err
}

// TestDecodeWithSpecialCharacters tests decoding CSV with special characters
func TestDecodeWithSpecialCharacters(t *testing.T) {
	csvData := "description\n\"Hello, World\"\n\"Line1\nLine2\"\n\"Quote\"\"Test\""
	reader := strings.NewReader(csvData)
	decoder := NewDecoder(reader)

	var result struct {
		Description []string `csv:"description"`
	}

	err := decoder.Decode(&result)
	if err != nil {
		t.Errorf("Decode with special characters failed: %v", err)
	}

	assertSlicesEqual(t, "Description", result.Description, []string{
		"Hello, World",
		"Line1\nLine2",
		"Quote\"Test",
	})
}

// TestDecodeMultipleTimes tests calling Decode multiple times
func TestDecodeMultipleTimes(t *testing.T) {
	csvData := "name,value\nAlice,100\nBob,200"
	reader := strings.NewReader(csvData)
	decoder := NewDecoder(reader)

	var result1 struct {
		Name  []string `csv:"name"`
		Value []int    `csv:"value"`
	}

	err := decoder.Decode(&result1)
	if err != nil {
		t.Errorf("First Decode failed: %v", err)
	}
	assertSlicesEqual(t, "result1.Name", result1.Name, []string{"Alice", "Bob"})
	assertSlicesEqual(t, "result1.Value", result1.Value, []int{100, 200})

	// Second decode on exhausted reader should return an error (EOF on header read)
	var result2 struct {
		Name  []string `csv:"name"`
		Value []int    `csv:"value"`
	}
	err = decoder.Decode(&result2)
	if err != nil {
		t.Error("Second Decode on exhausted reader should not return an error")
	}
}

type PersonRow struct {
	Name  string  `csv:"name"`
	Age   int     `csv:"age"`
	Score float64 `csv:"score"`
}

// TestDecodeArrayOfStructs tests decoding into a slice of structs
func TestDecodeArrayOfStructs(t *testing.T) {
	csvData := "name,age,score\nAlice,30,9.5\nBob,25,8.0\nCharlie,35,7.3"
	reader := strings.NewReader(csvData)
	decoder := NewDecoder(reader)

	var result []PersonRow
	err := decoder.Decode(&result)
	if err != nil {
		t.Fatalf("Decode array of structs failed: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("length mismatch: got %d, want 3", len(result))
	}

	want := []PersonRow{
		{"Alice", 30, 9.5},
		{"Bob", 25, 8.0},
		{"Charlie", 35, 7.3},
	}
	for i, got := range result {
		if got != want[i] {
			t.Errorf("result[%d]: got %+v, want %+v", i, got, want[i])
		}
	}
}

// TestDecodeArrayOfStructsSingleRow tests decoding a single row into a slice of structs
func TestDecodeArrayOfStructsSingleRow(t *testing.T) {
	csvData := "name,age,score\nZara,28,10.0"
	reader := strings.NewReader(csvData)

	var result []PersonRow
	err := NewDecoder(reader).Decode(&result)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("length mismatch: got %d, want 1", len(result))
	}
	if result[0] != (PersonRow{"Zara", 28, 10.0}) {
		t.Errorf("result[0]: got %+v, want {Zara 28 10.0}", result[0])
	}
}

// TestDecodeArrayOfStructsHeaderOnly tests decoding a CSV with only a header row
func TestDecodeArrayOfStructsHeaderOnly(t *testing.T) {
	csvData := "name,age,score"
	reader := strings.NewReader(csvData)

	var result []PersonRow
	err := NewDecoder(reader).Decode(&result)
	if err != nil {
		t.Fatalf("Decode header-only CSV failed: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(result))
	}
}

// TestDecodeArrayOfStructsMissingColumn tests error when a tagged column is absent
func TestDecodeArrayOfStructsMissingColumn(t *testing.T) {
	csvData := "name,age\nAlice,30" // "score" column missing
	reader := strings.NewReader(csvData)

	var result []PersonRow
	err := NewDecoder(reader).Decode(&result)
	if err == nil {
		t.Error("expected error for missing column, got nil")
	}
}

// TestDecodeArrayOfStructsIgnoresUntaggedFields tests that fields without csv tags are skipped
func TestDecodeArrayOfStructsIgnoresUntaggedFields(t *testing.T) {
	type PartialRow struct {
		Name    string `csv:"name"`
		Ignored int    // no tag — should be left as zero value
	}

	csvData := "name,age\nAlice,30\nBob,25"
	reader := strings.NewReader(csvData)

	var result []PartialRow
	err := NewDecoder(reader).Decode(&result)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("length mismatch: got %d, want 2", len(result))
	}
	assertSlicesEqual(t, "Name", []string{result[0].Name, result[1].Name}, []string{"Alice", "Bob"})
	if result[0].Ignored != 0 || result[1].Ignored != 0 {
		t.Errorf("untagged field should be zero, got %d and %d", result[0].Ignored, result[1].Ignored)
	}
}

// TestDecodeArrayOfStructsNonStructElem tests error when slice element is not a struct
func TestDecodeArrayOfStructsNonStructElem(t *testing.T) {
	csvData := "value\n1\n2"
	reader := strings.NewReader(csvData)

	var result []int
	err := NewDecoder(reader).Decode(&result)
	if err == nil {
		t.Error("expected error for non-struct slice element, got nil")
	}
}
