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

// Unpack JSON and call the validate function if no
// errors are found when unpacking it.
func Guaranty(obj interface{}, input io.ReadCloser) error {
	if b, err := ioutil.ReadAll(input); err == nil && string(b) != "{}" && string(b) != "" {
		// Turn our string back into a io.Reader if it's valid
		decoder := json.NewDecoder(bytes.NewReader(b))

		if err := decoder.Decode(obj); err == nil {
			return validate(obj)
		} else {
			return err
		}
	} else if err == nil {
		return errors.New("Nothing was passed in or JSON featured an empty object.")
	} else {
		return err
	}
}

// validate handles ensuting that all passed in values
// are safe to use before calling out to other functions
// to ensure the proper validation and type
func validate(obj interface{}) error {

	typ := reflect.TypeOf(obj)
	value := reflect.ValueOf(obj)

	// Check to ensure we are getting a valid
	// pointer for manipulation.
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = value.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {

		field := typ.Field(i)
		fieldValue := value.Field(i).Interface()
		zero := reflect.Zero(field.Type).Interface()

		// Validate nested and embedded structs (if pointer, only do so if not nil)
		if field.Type.Kind() == reflect.Struct ||
			(field.Type.Kind() == reflect.Ptr && !reflect.DeepEqual(zero, fieldValue)) {
			if err := validate(fieldValue); err != nil {
				return err
			}
		}

		if field.Tag.Get("validate") != "" {
			// Break validate field into array
			array := strings.Split(field.Tag.Get("validate"), "|")

			// Do the hard work of checking all assertions
			for setting := range array {

				match := array[setting]

				switch {
				case "required" == match:
					if err := required(field, fieldValue, zero); err != nil {
						return err
					}
				case "email" == match:
					if err := email(fieldValue); err != nil {
						return err
					}
				case "in:" == match[0:3]:
					if err := in(match, fieldValue); err != nil {
						return err
					}
				default:
					// Temp logging to check for errors
					errors.New(array[setting] + " is not a valid validation type.")
				}
			}
		}
	}

	return nil
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
func email(value interface{}) error {
	// Email Regex Checker
	var emailRegex string = `^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`

	if data, ok := value.(string); ok {
		if match, _ := regexp.Match(emailRegex, []byte(data)); match {
			return nil
		} else {
			return errors.New("A valid email address was not entered.")
		}
	} else {
		return errors.New("Email was not able to convert the passed in data to a []byte.")
	}
}

// Check that the passed in field is a valid email
// Need to improve error logging for this method
// Currently only supports strings, ints
func in(field string, value interface{}) error {

	if data, ok := value.(string); ok {

		valid := strings.Split(field[3:], ",")

		for option := range valid {
			if valid[option] == data {
				return nil
			}
		}

		return errors.New("In did not match any of the expected values.")

	} else if data, ok := value.(int); ok {
		// This will run with passed in data is an int
		valid := strings.Split(field[3:], ",")

		for option := range valid {
			// Check for convertion to valid int
			if valint, err := strconv.ParseInt(valid[option], 0, 64); err == nil {
				if valint == int64(data) {
					return nil
				}
			}
		}

		return errors.New("In did not match any of the expected values.")

	} else {
		return errors.New("in, was not able to convert the data passed in to a string.")
	}

}
