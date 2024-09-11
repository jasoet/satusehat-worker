package hash

import "testing"

func TestConcatHashWithMD5(t *testing.T) {
	// Test case 1: Concatenate empty slice of strings
	var input1 []string
	expectedOutput1 := "d41d8cd98f00b204e9800998ecf8427e"
	actualOutput1, err := WithMD5(input1)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if actualOutput1 != expectedOutput1 {
		t.Errorf("Incorrect output: expected %s but got %s", expectedOutput1, actualOutput1)
	}

	// Test case 2: Concatenate single string
	input2 := []string{"foo"}
	expectedOutput2 := "acbd18db4cc2f85cedef654fccc4a4d8"
	actualOutput2, err := WithMD5(input2)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if actualOutput2 != expectedOutput2 {
		t.Errorf("Incorrect output: expected %s but got %s", expectedOutput2, actualOutput2)
	}

	// Test case 3: Concatenate multiple strings
	input3 := []string{"foo", "bar", "baz"}
	expectedOutput3 := "6df23dc03f9b54cc38a0fc1483df6e21"
	actualOutput3, err := WithMD5(input3)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if actualOutput3 != expectedOutput3 {
		t.Errorf("Incorrect output: expected %s but got %s", expectedOutput3, actualOutput3)
	}
}
