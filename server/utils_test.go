package main

import (
	"reflect"
	"testing"
)

func TestString2TagSlice(t *testing.T) {
	want := make([]string, 3)
	want[0] = "tag1"
	want[1] = "tag2"
	want[2] = "tag3"
	tagString := "#tag1 #tag2 #tag3"
	if got := String2TagSlice(tagString); !reflect.DeepEqual(want, got) {
		t.Fatalf("Expected: %q, got: %q", want, got)
	}

	tagString = " #tag1        #tag2#tag3"
	if got := String2TagSlice(tagString); !reflect.DeepEqual(want, got) {
		t.Fatalf("Expected: %q, got: %q", want, got)
	}

	tagString = " #tag1#tag2 #tag3 "
	if got := String2TagSlice(tagString); !reflect.DeepEqual(want, got) {
		t.Fatalf("Expected: %q, got: %q", want, got)
	}

	wantEmptyStrSlice := make([]string, 1)
	wantEmptyStrSlice[0] = ""
	tagString = ""
	if got := String2TagSlice(tagString); !reflect.DeepEqual(wantEmptyStrSlice, got) {
		t.Fatalf("Expected: %q, got: %q", want, got)
	}

	tagString = "#"
	if got := String2TagSlice(tagString); !reflect.DeepEqual(wantEmptyStrSlice, got) {
		t.Fatalf("Expected: %q, got: %q", want, got)
	}

	tagString = "tag1tag2tag3"
	if got := String2TagSlice(tagString); !reflect.DeepEqual(wantEmptyStrSlice, got) {
		t.Fatalf("Expected: %q, got: %q", want, got)
	}

	tagString = "notexist#tag1#tag2"
	wantNoBeforNumChar := make([]string, 2)
	wantNoBeforNumChar[0] = "tag1"
	wantNoBeforNumChar[1] = "tag2"
	if got := String2TagSlice(tagString); !reflect.DeepEqual(wantNoBeforNumChar, got) {
		t.Fatalf("Expected: %q, got: %q", want, got)
	}

	tagString = "   notexist    #tag1 #tag2 "
	if got := String2TagSlice(tagString); !reflect.DeepEqual(wantNoBeforNumChar, got) {
		t.Fatalf("Expected: %q, got: %q", want, got)
	}
}
