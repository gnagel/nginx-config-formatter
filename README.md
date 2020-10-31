# Nginx Config Formatter
nginx config file formatter/beautifier written in Go. 
This is a port of the original formatter written in python here: [1connect/nginx-config-formatter](https://github.com/1connect/nginx-config-formatter)

This Go app script formats *nginx* configuration files in consistent way, described below:

* all lines are indented in uniform manner, with 4 spaces per level
* neighbouring empty lines are collapsed to at most two empty lines
* curly braces placement follows Java convention
* whitespaces are collapsed, except in comments an quotation marks
* whitespaces in variable designators are removed: `${  my_variable }` is collapsed to `${my_variable}`

## Installation

Go get to run this app:

    go get github.com/gnagel/nginx-config-formatter
    go tool github.com/gnagel/nginx-config-formatter ...  


## Usage

```
usage: go tool github.com/gnagel/nginx-config-formatter [-h] [-v] [-b] config_files [config_files ...]

Script formats nginx configuration file.

positional arguments:
  config_files          configuration files to format

optional arguments:
  -h, --help            show this help message and exit
  -v, --verbose         show formatted file names
  -b, --backup-original
                        backup original config file
```

## Credits

Copyright 2020 G. Nagel, credit for original inspiration goes to Michał Słomkowski.