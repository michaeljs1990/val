package main

import (
	"fmt"
	"github.com/michaeljs1990/val"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

// Make string into io.ReadCloser
// this is just for convinience
func jsonFactory(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}

func TestSpeedAverage(t *testing.T) {

	type Register struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required|email"`
		Type     string `json:"type" validate:"required|in:admin,user,guest"`
	}

	startTime := time.Now()

	// Benchmark an average JSON request
	// Currently ~3.7s for 100k requests
	for i := 0; i < 100000; i++ {

		var mockRegister Register

		testJSON := jsonFactory(`{"username": "michaeljs1990", "password": "secret", "email": "michaeljs1990@gmail.com", "type": "admin"}`)

		if err := val.Guaranty(&mockRegister, testJSON); err != nil {
			t.Error(err)
		}

	}

	endTime := time.Now()

	fmt.Printf("val general test took: %v to run.\n", endTime.Sub(startTime))
}

func TestSpeedEmail(t *testing.T) {

	type Email struct {
		Email string `json:"email" validate:"required|email"`
	}

	startTime := time.Now()

	// Benchmark an average JSON request
	// Currently ~3.7s for 100k requests
	for i := 0; i < 100000; i++ {

		var mockEmail Email

		testJSON := jsonFactory(`{"email": "michaeljs1990@gmail.com"}`)

		if err := val.Guaranty(&mockEmail, testJSON); err != nil {
			t.Error(err)
		}

	}

	endTime := time.Now()

	fmt.Printf("val email test took: %v to run.\n", endTime.Sub(startTime))
}

func TestSpeedIn(t *testing.T) {

	type In struct {
		Type string `json:"type" validate:"in:admin,user,guest"`
	}

	startTime := time.Now()

	// Benchmark an average JSON request
	// Currently ~3.7s for 100k requests
	for i := 0; i < 100000; i++ {

		var mockIn In

		testJSON := jsonFactory(`{"type": "admin"}`)

		if err := val.Guaranty(&mockIn, testJSON); err != nil {
			t.Error(err)
		}

	}

	endTime := time.Now()

	fmt.Printf("val in test took: %v to run.\n", endTime.Sub(startTime))
}

func TestSpeedRequired(t *testing.T) {

	type Required struct {
		Type string `json:"type" validate:"required"`
	}

	startTime := time.Now()

	// Benchmark an average JSON request
	// Currently ~3.7s for 100k requests
	for i := 0; i < 100000; i++ {

		var mockRequired Required

		testJSON := jsonFactory(`{"type": "admin"}`)

		if err := val.Guaranty(&mockRequired, testJSON); err != nil {
			t.Error(err)
		}

	}

	endTime := time.Now()

	fmt.Printf("val required test took: %v to run.\n", endTime.Sub(startTime))
}

func TestSpeedDigitInt(t *testing.T) {

	type DigitInt struct {
		Number int `json:"number" validate:"digit:4"`
	}

	startTime := time.Now()

	// Benchmark an average JSON request
	// Currently ~3.7s for 100k requests
	for i := 0; i < 100000; i++ {

		var mockDigitInt DigitInt

		testJSON := jsonFactory(`{"number": 1000}`)

		if err := val.Guaranty(&mockDigitInt, testJSON); err != nil {
			t.Error(err)
		}

	}

	endTime := time.Now()

	fmt.Printf("val digit int test took: %v to run.\n", endTime.Sub(startTime))
}

func TestSpeedDigitBetweenInt(t *testing.T) {

	type DigitBetweenInt struct {
		Number int `json:"number" validate:"digit_between:4,6"`
	}

	startTime := time.Now()

	// Benchmark an average JSON request
	// Currently ~3.7s for 100k requests
	for i := 0; i < 100000; i++ {

		var mockDigitBetweenInt DigitBetweenInt

		testJSON := jsonFactory(`{"number": 10000}`)

		if err := val.Guaranty(&mockDigitBetweenInt, testJSON); err != nil {
			t.Error(err)
		}

	}

	endTime := time.Now()

	fmt.Printf("val digit string test took: %v to run.\n", endTime.Sub(startTime))
}
