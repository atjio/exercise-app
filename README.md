# exercise-app

This repository contains two applications: `SimpleApp` and `LoadBalancer`.

## How To Start

### SimpleApp
1. Navigate to the `SimpleApp` directory:
   ```bash
   cd SimpleApp
   ```
2. Run the application using the following command:
   ```bash
   go run . -port=<port_number>
   ```
   - Replace `<port_number>` with your desired port.
   - If no port is specified, it defaults to `1234`.
   - You can run multiple instances of `SimpleApp` on different ports.

### LoadBalancer
1. Navigate to the `LoadBalancer` directory:
   ```bash
   cd LoadBalancer
   ```
2. Run the application using the following command:
   ```bash
   go run . -port=<port_number>
   ```
   - Replace `<port_number>` with your desired port.
   - If no port is specified, it defaults to `8080`.

### Testing
Once both applications are running, you can use the provided `postman_collection.json` to test the functionality.

## Routes

### SimpleApp

- **POST /echo**
  - Description: Echoes back the request body.
  - Request Body: Any JSON payload.
  - Response: The same JSON payload that was sent.

- **GET /healthcheck**
  - Description: Checks the health of the service.
  - Response:
    ```json
    {
      "status": "healthy"
    }
    ```

- **GET /debug**
  - Description: Returns debugging information, including the service's port number and a list of registered services (if any).
  - Response:
    ```json
    {
      "port": "<port_number>",
      "registered_services": [] // Example, may contain registered services
    }
    ```

- **POST /simulateDelay**
  - Description: Introduces an artificial delay to simulate processing time. The delay duration is specified in the request body.
  - Request Body:
    ```json
    {
      "delay": <duration_in_milliseconds>
    }
    ```
    - Example: `{"delay": 1000}` for a 1-second delay.
  - Response:
    ```json
    {
      "message": "Delay simulated for <duration_in_milliseconds> ms"
    }
    ```

### LoadBalancer

- **GET /register**
  - Description: Manually registers a new instance of `SimpleApp` with the load balancer. (Note: `SimpleApp` instances also attempt to auto-register.)
  - Query Parameters:
    - `address`: The address of the `SimpleApp` instance (e.g., `localhost:1234`).
  - Response:
    ```json
    {
      "message": "Service registered successfully"
    }
    ```
    or an error message if registration fails.

- **POST /***
  - Description: Forwards incoming POST requests to one of the registered `SimpleApp` instances using a round-robin strategy. The path and body of the request are forwarded as is.
  - Request Path: Any path (e.g., `/echo`, `/simulateDelay`).
  - Request Body: Any JSON payload, depending on the target `SimpleApp` route.
  - Response: The response from the chosen `SimpleApp` instance.