Gecho
==============================

[![License](https://img.shields.io/badge/license-MIT-green.svg)](https://git.thebarrens.nu/wolvie/gecho/blob/master/LICENSE)

A Simple echo server written in Go

Usage
-----

Just compile (or build the container image)  and run the application, the only supported parameters are

* `MAX_REQUESTS` environment variable which defaults to 500 if not set, used to limit the number of requests the server will accept
* `TOKEN` environment variable which defaults to `token` if not set, used as a authentication token for the /reset endpoint

