package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestAccumulatorAdd(t *testing.T) {
	a := NewAccumulator(3)

	a.Add("I")
	a.Add("love")
	a.Add("tacos")
	a.Add("they")
	a.Add("are")
	a.Add("so")
	a.Add("yum")
	a.Add("they")
	a.Add("are")
	a.Add("so")
	a.Add("delicious")

	if _, ok := a.groupings["I love tacos"]; !ok {
		t.Errorf("collection is missing key: %s -- %v", "I love tacos", a.groupings)
	}

	if _, ok := a.groupings["love tacos they"]; !ok {
		t.Errorf("collection is missing key: %s -- %v", "love tacos they", a.groupings)
	}

	if _, ok := a.groupings["tacos they are"]; !ok {
		t.Errorf("collection is missing key: %s -- %v", "tacos they are", a.groupings)
	}
}

func TestScanWords(t *testing.T) {
	tests := []struct {
		data  []byte
		eof   bool
		adv   int
		token []byte
		err   error
	}{
		{[]byte(`tacos `), true, 6, []byte(`tacos`), nil},
		{[]byte(`1231 `), true, 5, []byte(``), nil},
		{[]byte(`     `), true, 5, []byte(``), nil},
		{[]byte(`blar's`), true, 6, []byte(`blars`), nil},
		{[]byte(`well, `), true, 6, []byte(`well`), nil},
		{[]byte(`\n`), true, 2, []byte(`n`), nil},
	}

	for i, tc := range tests {
		adv, token, err := ScanWords(tc.data, tc.eof)
		if adv != tc.adv || !bytes.Equal(token, tc.token) || !errors.Is(err, tc.err) {
			t.Errorf("Test at index %d failed. Expected %d, '%s', %s but received %d, '%s', %s", i, tc.adv, tc.token, tc.err, adv, token, err)
		}
	}
}
