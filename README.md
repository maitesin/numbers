# Numbers

Application that receives multiple connections and writes down to file deduplicated numbers from the ones received through the connections.

## Project structure

- **cmd/generator**: contains an executable that generates random numbers of 9 digits and prints them to stdout.
- **cmd/numbers**: contains the main executable of the project.
- **config**: contains the general configuration of the project. It follows the [12 factor](https://12factor.net/config) configuration design.
- **internal**: contains the Domain Driven Design approach to implement the requirements of the task.
- **internal/app**: contains the application layer, where the use cases for the project have been implemented.
- **internal/domain**: contains the domain layer.
- **internal/infra**: contains the infrastructure layer, where the network client manager is implemented.

## How to run the project

The following dependencies are required for the project to be executed. `make`, and `go`.

### Run the binary

Now that the DB is up and running you can start the execution of the project binary:

```bash
make run
```

## How to test the project

Besides, the mandatory dependencies to run the project we will need the following extra dependencies to run test and run linting checks in the project. `golangci-lint`, and `moq`. They can be installed with the following:

```bash
make tools
```

To run the linting checks it is similar to the previous command:

```bash
make lint
```

### Unit test

In order to run the unit test you just need the following command:

```bash
make test
```

### Performance test

You can use the generate main contained in the repository to test the throughput of the solution. If you have `netcat` installed in your system, while the numbers binary is running, explained in the section above, you can execute the following command to send request for 21 seconds.

```bash
make send-numbers-for-21-seconds
```

In my home desktop (Intel i5-4670K and 16GB RAM) I obtained a throughput of 4M.

```bash
$ make run
cd cmd/numbers && go run main.go
Received 3281026 unique numbers, 5463 duplicates. Unique total: 3286489
Received 3958772 unique numbers, 21250 duplicates. Unique total: 7266511
```

## Configuration

As mentioned above the project follows the [12 factor](https://12factor.net/config) configuration design. Therefore, it uses environment variables in order to configure several aspects of the project:

* `NUMBERS_CONCURRENT_CLIENTS`: sets the number of concurrent clients allowed. It's default value is `5`.
* `NUMBERS_TIME_BETWEEN_REPORTS`: sets the time between reports (in seconds). It's default value is `10`.
* `NUMBERS_HOST`: sets the host name to listen for TCP requests. It's default value is `127.0.0.1`.
* `NUMBERS_PORT`: sets the port to listen for TCP requests. It's default value is `4000`.

## Exit codes

As it is common in the Unix-like system applications, if the `numbers` binary stops the execution by entering the `terminate` command. It will return 0 as its exit code. Please, see the list below with the different exit codes and what do they mean:

- 0: gracefully shutdown performed.
- 1: failure to open `numbers.log` file.
- 2: configuration failure.
- 3: failure to bind to the port in the interface specified.
