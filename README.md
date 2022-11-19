# Facial Recognition Photo Sorter

[![Go](https://github.com/jtowe1/photo-sort/actions/workflows/go.yml/badge.svg)](https://github.com/jtowe1/photo-sort/actions/workflows/go.yml)


This is a command line application that takes as input the path to a folder full of photos </br>
and sorts them into albums using facial recognition.

## Setup
This application uses AWS Rekognition, so you need credentials setup for that.

Create the following files:

`~/.aws/credentials`
```
[photo-sort]
aws_access_key_id = <your access key id here>
aws_secret_access_key = <your access key here>
```

`~/.aws/config`
```
[profile photo-sort]
region = <your region here>
```

## Usage
```bash
$ go run ./cmd/cli --pathToPhotos ~/path/to/photos

"duration:  20.196546208s"
"photos sorted and placed in /path/to/photos/sorted"

```

You can also enable debug mode for more output
```bash
$ go run ./cmd/cli --pathToPhotos ~/path/to/photos --debug
```