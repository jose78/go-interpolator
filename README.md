# judo_interpolator

judo_interpolator helps you to interpolate your vars inside the chain. Based on the go templates and using the large list of functions provided by the [Sprin Functions Project](https://masterminds.github.io/sprig/), judo_interpolator wants to make it easier to work with strings.


## How-to install:

To install the dependency:

```bash
go get -u github.com/judoDSL/judo_interpolator
```

## How-to use:

You can see an example about how to use it. 

```go
package main

import "github.com/judoDSL/judo_interpolator"

func main () {
	values := make(map[string] interface{})
	values["name"] = "            Jose                 "
	values["main_topic"] = "restore the snyderverse"
	values["favorite_superhero"] = "batman who laughs"
	judo_interpolator.Do("I'm {{ .name | trim }} and I want to {{ .main_topic | upper  }} because I would like to see a film related with {{ .favorite_superhero | title }}", values).Println()
}
```



The result of this esecution is:

```bash
[jose78@~/ws/test_judo_interpolator] $  go run main.go 
I'm Jose and I want to RESTORE THE SNYDERVERSE because I would like to see a film related with Batman Who Laughs
```
