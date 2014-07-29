/*
This package allows for easy validation of passed in json.
Val does not intend to be a robust solution but does seek to cover 95% of use cases.
Val requires a structure to use pointers for validation. This may seem odd but if a pointer is
not used you will run into some strange issues since json.Decode() will pass an int type back 
set as 0 giving no way to tell if a 0 was actually passed in or not. Using a pointer allows to
check for a nil value before doing the validation and lets you have optional json parameters.

Basic Struct Example.

    var Register struct {
        Username *string `json:"username" validate:"required"`
        Password *string `json:"password" validate:"required"`
        Email    *string `json:"email" validate:"required|email"`
        Notify   *string `json:"notify" validate:"required|in:yes,no"`
    }

Normal Use Case.

    if err := val.Bind(r.Body, &Register); err != nil {
        fmt.Println(err)
    }

*/

package val

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Unpack JSON and call the validate function if no errors are found when unpacking it.
// Bind kicks of the validation process. Note that Request.Body impliments an io.ReadCloser.
// Look into ReadAll http://jmoiron.net/blog/crossing-streams-a-love-letter-to-ioreader/
func Bind(input io.ReadCloser, obj interface{}) error {
	// Don't go through any logic if nothing was passed in.
	if b, err := ioutil.ReadAll(input); err == nil && string(b) != "{}" && string(b) != "" {
		// Turn our string back into a io.Reader if it's valid
		decoder := json.NewDecoder(bytes.NewReader(b))

		if err := decoder.Decode(obj); err == nil {
			return Validate(obj)
		} else {
			return err
		}
	} else if err == nil {
		return errors.New("Nothing was passed in or JSON featured an empty object.")
	} else {
		return err
	}
}

// In version 1.0 I exported the Validation function. This can be used when you may
// not need to or want to have JSON first converted into a struct.
func Validate(obj interface{}) error {

	typ := reflect.TypeOf(obj)
	value := reflect.ValueOf(obj)

	// Check to ensure we are getting a valid
	// pointer for manipulation.
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = value.Elem()
	}

	// Kill process if obj did not pass in a scruct.
	// This happens when a pointer passed in.
	if value.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < typ.NumField(); i++ {

		field := typ.Field(i)
		fieldValue := value.Field(i).Interface()
		zero := reflect.Zero(field.Type).Interface()

		// Validate nested and embedded structs (if pointer, only do so if not nil)
		if field.Type.Kind() == reflect.Struct ||
			(field.Type.Kind() == reflect.Ptr && !reflect.DeepEqual(zero, fieldValue)) {
			if err := Validate(fieldValue); err != nil {
				return err
			}
		}

		if field.Tag.Get("validate") != "" || field.Tag.Get("binding") != "" {
			// Break validate field into array
			array := strings.Split(field.Tag.Get("validate"), "|")

			// Legacy Support for binding.
			if array[0] == "" {
				array = strings.Split(field.Tag.Get("binding"), "|")
			}

			// Do the hard work of checking all assertions
			for setting := range array {

				match := array[setting]

				//Check that value was passed in and is not required.
				if match != "required" && null(fieldValue) == true {
					return nil
				}

				switch {
				case "required" == match:
					if err := required(field, fieldValue, zero); err != nil {
						return err
					}
				case "email" == match:
					if err := regex(`regex:^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`, fieldValue); err != nil {
						return err
					}
				case "url" == match:
					if err := regex(`regex:/^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/`, fieldValue); err != nil {
						return err
					}
				case "alpha" == match:
					if err := regex(`regex:\p{L}`, fieldValue); err != nil {
						return err
					}
				case "alphadash" == match:
					if err := regex(`regex:^[a-zA-Z0-9_]*$`, fieldValue); err != nil {
						return err
					}
				case "alphanumeric" == match:
					if err := regex(`regex:/[0-9a-zA-Z]/`, fieldValue); err != nil {
						return err
					}
				case strings.HasPrefix(match, "min:"):
					if err := min(match, fieldValue); err != nil {
						return err
					}
				case strings.HasPrefix(match, "max:"):
					if err := max(match, fieldValue); err != nil {
						return err
					}
				case strings.HasPrefix(match, "in:"):
					if err := in(match, fieldValue); err != nil {
						return err
					}
				case strings.HasPrefix(match, "regex:"):
					if err := regex(match, fieldValue); err != nil {
						return err
					}
				case strings.HasPrefix(match, "length:"):
					if err := length(match, fieldValue); err != nil {
						return err
					}
				case strings.HasPrefix(match, "length_between:"):
					if err := length_between(match, fieldValue); err != nil {
						return err
					}
				default:
					panic("The field " + match + " is not a valid validation check.")
				}
			}
		}
	}

	return nil
}

// Ensure that the value being passed in is not of type nil.
func null(value interface{}) bool {
	if reflect.ValueOf(value).IsNil() {
		return true
	}

	return false
}

// Check that the following function features
// the required field. May need to check for
// more special cases like since passing in null
// is the same as 0 for int type checking.
func required(field reflect.StructField, value, zero interface{}) error {

	if reflect.DeepEqual(zero, value) {
		if _, ok := value.(int); !ok {
			return errors.New("The required field " + field.Name + " was not submitted.")
		}
	}

	return nil
}

// Check that the passed in field is a valid email
// Need to improve error logging for this method
// Currently only supports strings, ints
func in(field string, value interface{}) error {

	if data, ok := value.(*string); ok {
		if len(*data) == 0 {
			return nil
		}

		valid := strings.Split(field[3:], ",")

		for option := range valid {
			if valid[option] == *data {
				return nil
			}
		}

	} else {
		return errors.New("The value passed in for IN could not be converted to a string.")
	}

	return errors.New("In did not match any of the expected values.")
}

func min(field string, value interface{}) error {

	if data, ok := value.(*int); ok {

		min := field[strings.Index(field, ":")+1:]

		if minNum, ok := strconv.ParseInt(min, 0, 64); ok == nil {

			if int64(*data) >= minNum {
				return nil
			} else {
				return errors.New("The data you passed in was smaller then the allowed minimum.")
			}

		}
	}

	return errors.New("The value passed in for MIN could not be converted to an int.")
}

func max(field string, value interface{}) error {

	if data, ok := value.(*int); ok {

		max := field[strings.Index(field, ":")+1:]

		if maxNum, ok := strconv.ParseInt(max, 0, 64); ok == nil {
			if int64(*data) <= maxNum {
				return nil
			} else {
				return errors.New("The data you passed in was larger than the maximum.")
			}

		}
	}

	return errors.New("The value passed in for MAX could not be converted to an int.")
}

// Regex handles the general regex call and also handles
// the regex email.
func regex(field string, value interface{}) error {

	reg := field[strings.Index(field, ":")+1:]

	if data, ok := value.(*string); ok {
		if len(*data) == 0 {
			return nil
		} else if err := match_regex(reg, []byte(*data)); err != nil {
			return err
		}
	} else if data, ok := value.(*int); ok {
		if err := match_regex(reg, []byte(strconv.Itoa(*data))); err != nil {
			return err
		}
	} else {
		return errors.New("The value passed in for REGEX could not be converted to a string or int.")
	}

	return nil
}

// Helper function for regex.
func match_regex(reg string, data []byte) error {

	if match, err := regexp.Match(reg, []byte(data)); err == nil && match {
		return nil
	} else {
		return errors.New("Your regex did not match or was not valid.")
	}
}

// Check passed in json length string is exact value passed in.
// Also checks if passed in values is between two different ones.
func length(field string, value interface{}) error {

	length := field[strings.Index(field, ":")+1:]

	if data, ok := value.(*string); ok {
		if intdata, intok := strconv.Atoi(length); intok == nil {
			if len(*data) == intdata {
				return nil
			} else {
				return errors.New("The data passed in was not equal to the expected length.")
			}
		} else {
			return errors.New("The value passed in for LENGTH could not be converted to an int.")
		}
	} else {
		return errors.New("The value passed in for LENGTH could not be converted to a string.")
	}
}

// Check if the strings length is between high,low.
func length_between(field string, value interface{}) error {

	length := field[strings.Index(field, ":")+1:]
	vals := strings.Split(length, ",")

	if len(vals) == 2 {

		if data, ok := value.(*string); ok {

			if lowerbound, lowok := strconv.Atoi(vals[0]); lowok == nil {

				if upperbound, upok := strconv.Atoi(vals[1]); upok == nil {

					if lowerbound <= len(*data) && upperbound >= len(*data) {
						return nil
					} else {
						return errors.New("The value passed in for LENGTH BETWEEN was not in bounds.")
					}

				} else {
					return errors.New("The value passed in for LENGTH BETWEEN could not be converted to an int.")
				}

			} else {
				return errors.New("The value passed in for LENGTH BETWEEN could not be converted to an int.")
			}

		} else {
			return errors.New("The value passed in for LENGTH BETWEEN could not be converted to a string.")
		}
	} else {
		return errors.New("LENGTH BETWEEN requires exactly two paramaters.")
	}
}
