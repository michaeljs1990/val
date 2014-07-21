package val

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

// Make string into io.ReadCloser
// this is just for convinience
func jsonSpeedFactory(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}

func TestSpeedAverage(t *testing.T) {

	type Register struct {
		Username *string `json:"username" validate:"required"`
		Password *string `json:"password" validate:"required"`
		Email    *string `json:"email" validate:"required|email"`
		Type     *string `json:"type" validate:"required|in:admin,user,guest"`
	}

	startTime := time.Now()

	// Benchmark an average JSON request
	// Currently ~3.7s for 100k requests
	for i := 0; i < 100000; i++ {

		var mockRegister Register

		testJSON := jsonSpeedFactory(`{"username": "michaeljs1990", "password": "secret", "email": "michaeljs1990@gmail.com", "type": "admin"}`)

		if err := Bind(testJSON, &mockRegister); err != nil {
			t.Error(err)
		}

	}

	endTime := time.Now()

	fmt.Printf("val general test took: %v to run.\n", endTime.Sub(startTime))
}

func TestSpeedEmail(t *testing.T) {

	type Email struct {
		Email *string `json:"email" validate:"required|email"`
	}

	startTime := time.Now()

	// Benchmark an average JSON request
	// Currently ~3.7s for 100k requests
	for i := 0; i < 100000; i++ {

		var mockEmail Email

		testJSON := jsonSpeedFactory(`{"email": "michaeljs1990@gmail.com"}`)

		if err := Bind(testJSON, &mockEmail); err != nil {
			t.Error(err)
		}

	}

	endTime := time.Now()

	fmt.Printf("val email test took: %v to run.\n", endTime.Sub(startTime))
}

func TestSpeedIn(t *testing.T) {

	type In struct {
		Type *string `json:"type" validate:"in:admin,user,guest"`
	}

	startTime := time.Now()

	// Benchmark an average JSON request
	// Currently ~3.7s for 100k requests
	for i := 0; i < 100000; i++ {

		var mockIn In

		testJSON := jsonSpeedFactory(`{"type": "admin"}`)

		if err := Bind(testJSON, &mockIn); err != nil {
			t.Error(err)
		}

	}

	endTime := time.Now()

	fmt.Printf("val in test took: %v to run.\n", endTime.Sub(startTime))
}

func TestSpeedRequired(t *testing.T) {

	type Required struct {
		Type *string `json:"type" validate:"required"`
	}

	startTime := time.Now()

	// Benchmark an average JSON request
	// Currently ~3.7s for 100k requests
	for i := 0; i < 100000; i++ {

		var mockRequired Required

		testJSON := jsonSpeedFactory(`{"type": "admin"}`)

		if err := Bind(testJSON, &mockRequired); err != nil {
			t.Error(err)
		}

	}

	endTime := time.Now()

	fmt.Printf("val required test took: %v to run.\n", endTime.Sub(startTime))
}

func TestSpeedLength(t *testing.T) {

	type DigitInt struct {
		Number *string `json:"number" validate:"length:4"`
	}

	startTime := time.Now()

	// Benchmark an average JSON request
	// Currently ~3.7s for 100k requests
	for i := 0; i < 100000; i++ {

		var mockDigitInt DigitInt

		testJSON := jsonSpeedFactory(`{"number": "1000"}`)

		if err := Bind(testJSON, &mockDigitInt); err != nil {
			t.Error(err)
		}

	}

	endTime := time.Now()

	fmt.Printf("val length test took: %v to run.\n", endTime.Sub(startTime))
}

func TestSpeedLengthBetween(t *testing.T) {

	type DigitBetweenInt struct {
		Number *string `json:"number" validate:"length_between:4,6"`
	}

	startTime := time.Now()

	// Benchmark an average JSON request
	// Currently ~3.7s for 100k requests
	for i := 0; i < 100000; i++ {

		var mockDigitBetweenInt DigitBetweenInt

		testJSON := jsonSpeedFactory(`{"number": "aaaa"}`)

		if err := Bind(testJSON, &mockDigitBetweenInt); err != nil {
			t.Error(err)
		}

	}

	endTime := time.Now()

	fmt.Printf("val length between test took: %v to run.\n", endTime.Sub(startTime))
}
