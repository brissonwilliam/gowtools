# gowtools

gowtools (go-willow-tools) is a collection of tools designed to simplify the usage of Go programming language. This project aims to provide developers with useful utilities and functionalities that enhance their Go development workflow. Please note that gowtools is currently in an incomplete state, and all contributions are more than welcome.


## Features

### Algo

- LSH (Locality Sensitive Hashing) index
- Jacard similarity

### Async

- WorkGroup: process a slice of data into a set number of workers

### Slice operations

- break into chunks
- paginate

### Sql

- build named values
- build on duplicate clauses from columns
- build on duplicate key increment values query
- prefix columns by table name

### Math

- min
- max
- modulus (handle negative values like other languages)
- round floats to decimal

### IP

- ip strings to uint
- uint to string ip
- ip block (masking)

### Countries

- Validate an iso alpha 3 country code

## Installation

To use gowtools, you need to have Go installed on your system. You can install gowtools by following these steps:

1. Open a terminal or command prompt.

2. Run the following command to install gowtools:

   ```shell
   go get github.com/brisonwilliam/gowtools

## Contributing

Contributions to gowtools are highly encouraged

## Licence

This project is licensed under the MIT Licence
