# go-hue

A typesafe, fully tested, go package for the Phillip's Hue client API. 

Note: This is a work in project and *NOT* ready to be used in production yet. I'm hoping to have it completed by the end of September 2013 and am actively working on the project.

Please see [godoc.org](http://godoc.org/github.com/bcurren/go-hue) for a detailed API description.

For more information about the Hue API see the [Hue Developer API Specification](http://developers.meethue.com/).

# Packages

This contains multiple packages. See the README in each package for more details. Here is a list for you reference.

* [huetest](huetest/) - Stub hue.API implementation to simplify testing.
* [strand](strand/) - Provides a LightStrand that maps location on a strand to a light id in hue. 

## Todo

- [X] Lights API
- [~] Configuration API
- [X] Discovery API
- [~] Documentation
- [ ] Groups API
- [ ] Schedules API

## Usage

* [Setup your go environment](http://golang.org/doc/code.html)
* ```go get github.com/bcurren/go-hue```
* Write code using the library.

## How to contribute
* Fork
* Write tests and code
* Run go fmt
* Submit a pull request
* Drink a beer. Good job!
