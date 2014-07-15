val
===
[![GoDoc](https://godoc.org/github.com/gin-gonic/gin?status.png)](http://godoc.org/github.com/michaeljs1990/val)
[![Build Status](https://travis-ci.org/michaeljs1990/val.svg?branch=master)](https://travis-ci.org/michaeljs1990/val)

Go JSON validation library.

This library was developed to meet some validation needs that I needed. However I would like to build this into a much more robust set of tools. Please feel free to ask for any feature or submit a pull request.

## Start using it
Run the following in your terminal to start using val

```
go get github.com/michaeljs1990/val
```
Then import it in your Go! code:

```
import "github.com/michaeljs1990/val"
```

## Example Usage

#### Basic example

```go
package main

import (
	"fmt"
	"github.com/michaeljs1990/val"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

	var Register struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"username" validate:"required"`
		Email    string `json:"username" validate:"required|email"`
		Notify   string `json:"username" validate:"required|in:yes,no"`
	}

	if err := val.Guaranty(&Register, r.Body); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("This validated!")
	}

}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

```

## Currently Supported Validation

#### required
This will ensure that the data is actually included in the json array.
```
Username string   `json:"username" validate:"required"`
```

#### email
This checks to see if the passed in value is a valid email, it uses the following regex "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$".
```
Username string   `json:"username" validate:"email"`
```

#### in
In support any length of arguments to validate a JSON string or int value against.
```
Username string   `json:"username" validate:"in:only,these,are,valid,strings"`
```

#### digit
Digit ensures the int passed in is exactly the right number of digits.
```
Username int   `json:"username" validate:"digit:5"`
```

#### digits_between
Digit between checks to make sure the data passed in will be between the two entered values. The following example uses the first value as the lower bound and the second as the upper bound.
```
Username string   `json:"username" validate:"digits_between:0,10"`
```

#### combinations
If you would like to ensure multiple conditions are met simply use the | character.
```
Username string   `json:"username" validate:"email|required|in:m@gmail.com,o@gmail.com"`
```
