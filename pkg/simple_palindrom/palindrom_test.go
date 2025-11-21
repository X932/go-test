package simple_palindrom

import (
	"fmt"
	"os"
	"testing"
)

// for test IsSimplePalindromList
// stringsForPalindrom := []string{"g", "goog", "erww"}

func TestPalindrom(t *testing.T) {
	inputForPalindrom := "gog"
	want := true
	isPalindrom, err := IsSimplePalindrom(inputForPalindrom)

	if want != isPalindrom || err != nil {
		t.Errorf(`IsSimplePalindrom(%q) = %v, %v, want match for %#v`, inputForPalindrom, isPalindrom, err, want)
	}

	fmt.Println("====== DATABASE_URL ======", os.Getenv("DATABASE_URL"))
}

func TestPalindromEmpty(t *testing.T) {
	inputForPalindrom := ""
	want := false
	isPalindrom, err := IsSimplePalindrom(inputForPalindrom)

	if want != isPalindrom || err == nil {
		t.Errorf(`IsSimplePalindrom(%q) = %v, %v, want match for %#v`, inputForPalindrom, isPalindrom, err, want)
	}
}
