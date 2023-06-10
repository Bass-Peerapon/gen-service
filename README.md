![Supported Go Versions](https://img.shields.io/badge/Go-1.19-lightgrey.svg)

# Gen-Service

`gen-service` is a command-line tool used for generating an innovative service directory template. It helps developers quickly set up the necessary structure for their service-based projects. 

## Installation
```bash
go install github.com/Bass-Peerapon/gen-service
```

🚀 **Features**
- Generate a service directory template
- Support for autocompletion scripts
- Easily create different types of services

## Usage

```bash
gen-service create [flags]
gen-service create [command]
```

## Available Commands

- `crudRepo` - Generate CRUD repository scripts.
- `model` - Generate a model based on a JSON file.
- `newWithParams` - Generate a function that creates a new instance with parameters.
- `service` - Generate a service.

To use a specific subcommand, use the following format:

```bash
gen-service create [command] [flags]
```

## Available Subcommands

### `crudRepo`

Generate CRUD repository scripts.

**Flags:**
- `-a, --all` - Generate SQL scripts for all CRUD operations (insert, query, update, delete).
- `-c, --create` - Generate SQL script for the insertion operation.
- `-d, --delete` - Generate SQL script for the deletion operation.
- `-f, --file string` - Specify the path to the model file.
- `-h, --help` - Display the help information for `crudRepo`.
- `-o, --output string` - Specify the path to the repository file.
- `-r, --read` - Generate SQL script for the querying operation.
- `-u, --update` - Generate SQL script for the updating operation.

**Example:**
```bash
gen-service create crudRepo --file ./path/to/model.go --output ./path/to/repository.go --all
```

This command generates SQL scripts for all CRUD operations based on the provided model file and outputs them to the specified repository file.

### `model`

Generate a model based on a JSON file.

**Flags:**
- `-f, --file_json string` - Specify the path to the JSON file.
- `-h, --help` - Display the help information for `model`.
- `-n, --model_name string` - Specify the model name.

**Example:**
```bash
gen-service create model --file_json ./path/to/model.json --model_name MyModel
```

This command generates a model based on the JSON file and assigns the specified model name to it.

### `newWithParams`

Generate a function that creates a new instance with parameters.

**Flags:**
- `-f, --file string` - Specify the path to the Go model file.
- `-h, --help` - Display the help information for `newWithParams`.

**Example:**
```bash
gen-service create newWithParams --file ./path/to/model.go
```

This command generates a function in the Go model file that creates a new instance with parameters.

### `service`

Generate a service.

**Flags:**
- `-h, --help` - Display the help information for `service`.
- `-n, --service_name string` - Specify the service name.

**Example:**
```bash
gen-service create service --service_name MyService
```

This command generates a service with the specified service name.

To learn more about each subcommand and its available flags, use the following format:

```bash
gen-service create [subcommand] --help
```

Feel free to explore and utilize the `gen-service create` command to generate the necessary components for your service. Happy coding!