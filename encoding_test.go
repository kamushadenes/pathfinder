package pathfinder

import "testing"

const testString = "potato"

func TestEncodeDecodeString(t *testing.T) {
	enc, err := EncodeString(testString)

	if err != nil {
		t.Error(err.Error())
	}

	if len(enc) == 0 {
		t.Fail()
	}

	dec, err := DecodeString(enc)

	if err != nil {
		t.Error(err.Error())
	}

	if dec != testString {
		t.Fail()
	}
}
