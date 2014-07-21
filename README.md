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

Pointers are needed so go's json.Decode() call can function as one would expect. Pointers are now required in version 0.1 they were not implimented which lead to issues if you wanted json that was not required since int's would be set to 0 even if nothing was passed in causing you to potentially fail a test.

```go
package main

import (
	"fmt"
	"github.com/michaeljs1990/val"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

	var Register struct {
		Username *string `json:"username" validate:"required"`
		Password *string `json:"password" validate:"required"`
		Email    *string `json:"email" validate:"required|email"`
		Notify   *string `json:"notification" validate:"required|in:yes,no"`
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

## Performance
I have created some benchmarks to see what really needs to be improved. Currently the benchmarks run 100,000 times and the performace is as follows. Below shows that email is an expensive call due to the non-optimized regex lib in go. Hopefully this will be improved over time. Other than that I am fairly happy with the current benchmarks. They would most likely be even a bit lower since on an http server you would not have to call a function every iteration to turn a string into a io.ReadCloser.

```
val general test took: 3.9768205s to run.
val email test took: 3.3563791s to run.
val in test took: 636.4503ms to run.
val required test took: 635.4693ms to run.
val length test took: 693.4902ms to run.
val length between test took: 713.4889ms to run.

```

## Currently Supported Validation

#### required
This will ensure that the data is actually included in the json array.
```
Username *string   `json:"username" validate:"required"`
```

#### email
This checks to see if the passed in value is a valid email, it uses the following regex "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$".
```
Username *string   `json:"username" validate:"email"`
```

#### in
In support any length of arguments to validate a JSON string or int value against.
```
Username *string   `json:"username" validate:"in:only,these,are,valid,strings"`
```

#### min
Min works with ints and ensures that the number the user has entered is not under the specified min. If the number is under it will return an error.
```
Username *string   `json:"username" validate:"min:10"`
```

#### max
Max works with ints and ensures that the number the user has entered is not over the specified max. If the number is over it will return an error.
```
Username *string   `json:"username" validate:"max:243"`
```

#### regex
Regex ensures that the string the user has passed in matched the regex you have entered. Currently this is only tested with strings.
```
Username *string   `json:"username" validate:"regex:\\d+"`
```

#### length
Regex ensures that the string the user has passed in matched the regex you have entered. Currently this is only tested with strings.
```
Username *string   `json:"username" validate:"length:5"`
```

#### length_between
Regex ensures that the string the user has passed in matched the regex you have entered. Currently this is only tested with strings.
```
Username *string   `json:"username" validate:"length_between:2,5"`
```

#### combinations
If you would like to ensure multiple conditions are met simply use the | character.
```
Username *string   `json:"username" validate:"email|required|in:m@gmail.com,o@gmail.com"`
```