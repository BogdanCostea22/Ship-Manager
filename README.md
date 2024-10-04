# Ship Manager

Ship Manager is a Go-based web application that helps calculate the optimal number of packs needed to fulfill an order, based on available pack sizes.

## Features

- Add and manage pack sizes
- Calculate the optimal pack combination for a given order size
- Clear all pack sizes
- Simple and intuitive web interface

## Live Demo

The application is deployed and accessible at:

https://ship-manager.fly.dev/calculator

**Note:** Due to the use of shared CPU resources, there might be a slight delay when accessing the application for the first time as the instance starts up.

## Local Development

### Prerequisites

- Go 1.23 or later
- templ
- Tailwind CSS
- Make

### Setting Up and Running the Application

1. Clone the repository
```
git clone https://github.com/BogdanCostea22/Ship-Manager.git
```
2. Install Go dependencies:
   ```
   go mod tidy
   ```
3. Build the application:
   ```
   make all build
   ```
4. Run the application:
   ```
   make run
   ```
5. Open `http://localhost:8080/calculator` in your browser

For development with live reload:
```
make watch
```

## Makefile Commands

- `make all build`: Run all make commands with clean tests and build the application
- `make build`: Build the application
- `make run`: Run the application
- `make watch`: Run the application with live reload
- `make test`: Run the test suite
- `make clean`: Clean up binary from the last build

For a full list of available commands, refer to the Makefile in the project root.

## Deployment

This project is deployed on [Fly.io](https://fly.io/). For deployment instructions, refer to the Fly.io documentation.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).