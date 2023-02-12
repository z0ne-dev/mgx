[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![Apache 2.0][license-shield]][license-url]

[![Go Reference][reference-shield]][reference-url]
[![Coverage][coverage-shield]][coverage-url]
[![Go Report Card][report-shield]][report-url]

[contributors-shield]: https://img.shields.io/github/contributors/z0ne-dev/mgx.svg?
[contributors-url]: https://github.com/z0ne-dev/mgx/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/z0ne-dev/mgx.svg?
[forks-url]: https://github.com/z0ne-dev/mgx/network/members
[stars-shield]: https://img.shields.io/github/stars/z0ne-dev/mgx.svg?
[stars-url]: https://github.com/z0ne-dev/mgx/stargazers
[issues-shield]: https://img.shields.io/github/issues/z0ne-dev/mgx.svg?
[issues-url]: https://github.com/z0ne-dev/mgx/issues
[coverage-shield]: https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/wiki/z0ne-dev/mgx/coverage-comment-badge.json&
[coverage-url]: https://github.com/z0ne-dev/mgx
[license-shield]: https://img.shields.io/github/license/z0ne-dev/mgx.svg?
[license-url]: https://github.com/z0ne-dev/mgx/blob/master/LICENSE.txt
[report-shield]: https://goreportcard.com/badge/github.com/z0ne-dev/mgx?
[report-url]: https://goreportcard.com/report/github.com/z0ne-dev/mgx
[reference-shield]: https://pkg.go.dev/badge/github.com/z0ne-dev/mgx.svg
[reference-url]: https://pkg.go.dev/github.com/z0ne-dev/mgx

# mgx

Simple migration system for [pgx](https://github.com/jackc/pgx).

Migrations are defined in code and are executed in order.
The migration system keeps track of which migrations have been executed and which have not.


## Getting Started

1. Install the dependency
```sh
go get -u github.com/z0ne-dev/mgx
```
2. Import the package and create a new migrator
```go
package main

import "github.com/z0ne-dev/mgx"

func main() {
    migrator, _ := mgx.New(mgx.Migrations(
		// insert migrations here 
    ))
}
```

3. Run the `Migrate(context.TODO(), pgx)` method to execute the migrations


## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request


## License

Distributed under the Apache-2.0 License. See `LICENSE` for more information.


## Acknowledgments

* [lopezator/migrator](https://github.com/lopezator/migrator) for inspiration for this package. Lots of inspiration was taken from this package, but it was not used directly. The API was designed to be similar, to reduce refactoring when switching between the two packages.
