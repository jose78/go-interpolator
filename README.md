# interpolator

The _interpolator_ project is part as _judoctl.sh_, the main goal of _interpolator_ is to help you to interpolate your _**nested vars**_ inside the string and evaluate functions related with this vars. It has been implemented to use 'go templates'  and to be highly configurable, enabling you to incorporate your own functions or 3rd parties, like [Sprig Functions Project](https://masterminds.github.io/sprig/).


## How-to install:

To install the dependency:

```bash
go get -u github.com/judoctl/interpolator
```

## How-to use:

In this example we are configuring the interpolator to use the [Sprig Functions Project](https://masterminds.github.io/sprig/) and templating the the input string with the **nested** variables.

```go
package main

import (
	"github.com/judoctl/interpolator"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
	)

func main () {
	
    customFuncMap := func() template.FuncMap  {
		return sprig.FuncMap()
	}

	runner := Configure(Configuration{FnProviderFunction: customFuncMap })

  values := make(map[string] interface{})
	values["name"] = "            Jose                 "
	values["main_topic"] = "restore the snyderverse"
	values["hero"] = "{{ .hero_redirect }}"
	values["favorite_superhero"] = "{{ .hero | upper }} who laughs"
	values["hero_redirect"] = "batman"
	
	str, _ := runner("I'm {{ .name | trim }} and I want to {{ .main_topic | upper  }} because I would like to see a film related with {{ .favorite_superhero | title }}", values)

	fmt.Print(str)
}
```


The result of this execution is:

```bash
[jose78@~/ws/interpolator] $  go run main.go 
I'm Jose and I want to RESTORE THE SNYDERVERSE because I would like to see a film related with BATMAN Who Laughs
```

> [!IMPORTANT]
>
> Please, note the value of **favorite_superhero**, it has been converted in **_BATMAN Who Laughs_**, interpolating the nested vars: *hero* and *hero_redirect*.

