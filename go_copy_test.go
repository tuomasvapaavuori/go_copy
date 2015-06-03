package go_copy

import (
	"testing"
)

func TestCopyStruct(t *testing.T) {
	type Test struct {
		Field1 string
		Field2 int
		Test   *Test
	}

	source := Test{}
	target := &Test{Field1: "BUHAA"}

	//source.Field1 = "HOGUS"
	source.Field2 = 12345
	source.Test = &Test{}
	source.Test.Field1 = "BOGUS"
	source.Test.Field2 = 98765

	Copy(source, &target, true)

	if target.Field1 != source.Field1 {
		t.Errorf("Expected %v, got %v", source.Field1, target.Field1)
	}

	if target.Field2 != source.Field2 {
		t.Errorf("Expected %v, got %v", source.Field2, target.Field2)
	}

	if target.Test.Field1 != source.Test.Field1 {
		t.Errorf("Expected %v, got %v", source.Test.Field1, target.Test.Field1)
	}

	if target.Test.Field2 != source.Test.Field2 {
		t.Errorf("Expected %v, got %v", source.Test.Field2, target.Test.Field2)
	}

	if target.Test.Test != source.Test.Test {
		t.Errorf("Expected %v, got %v", source.Test.Test, target.Test.Test)
	}
}

func TestCopySlice(t *testing.T) {
	var sourceSlice []int64
	var targetSlice []int64

	sourceSlice = append(sourceSlice, 124885)
	sourceSlice = append(sourceSlice, 5959595)

	Copy(sourceSlice, &targetSlice, true)

	if len(targetSlice) != 2 {
		t.Errorf("Expected %v, got %v", 2, len(targetSlice))
	}

	if targetSlice[0] != 124885 {
		t.Errorf("Expected %v, got %v", 124885, targetSlice[0])
	}
}

func TestCopyMap(t *testing.T) {
	sourceMap := make(map[string]string)
	var targetMap map[string]string

	sourceMap["This"] = "That"
	sourceMap["Hap"] = "Dhap"
	sourceMap["HapHap"] = "Dhap"

	Copy(sourceMap, &targetMap, true)

	if len(targetMap) != 3 {
		t.Errorf("Expected %v, got %v", 3, len(targetMap))
	}
}

func TestCopyChanMap(t *testing.T) {
	sourceMap := make(map[int]chan string)
	var targetMap map[int]chan string

	sourceMap[0] = make(chan string, 10)
	sourceMap[1] = make(chan string)

	Copy(sourceMap, &targetMap, true)

	if len(targetMap) != 2 {
		t.Errorf("Expected %v, got %v", 2, len(targetMap))
	}

	const (
		expectedStr1 = "HOGUS"
		expectedStr2 = "BOGUS"
	)

	var (
		targetValue1 string
		targetValue2 string
		chan1Closed  bool
		chan2Closed  bool
	)
	go func() {
		for {
			if chan1Closed && chan2Closed {
				return
			}

			select {
			case val, ok := <-targetMap[0]:
				if !ok {
					chan1Closed = true
					continue
				}

				targetValue1 = val
			case val, ok := <-targetMap[1]:
				if !ok {
					chan2Closed = true
					continue
				}

				targetValue2 = val
			}
		}
	}()

	targetMap[0] <- expectedStr1
	close(targetMap[0])
	targetMap[1] <- expectedStr2
	close(targetMap[1])

	if targetValue1 != expectedStr1 {
		t.Errorf("Expected %v, got %v", expectedStr1, targetValue1)
	}

	if targetValue2 != expectedStr2 {
		t.Errorf("Expected %v, got %v", expectedStr2, targetValue2)
	}
}
