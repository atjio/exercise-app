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

## Running with Docker

This project can also be run using Docker and Docker Compose, which simplifies deployment and scaling.

### Prerequisites
- Docker: [Install Docker](https://docs.docker.com/get-docker/)
- Docker Compose: [Install Docker Compose](https://docs.docker.com/compose/install/) (usually included with Docker Desktop)

### Building and Running
1.  **Navigate to the root of the project directory.**
    This is the directory containing the `docker-compose.yml` file.

2.  **Build and start the services:**
    ```bash
    docker-compose up
    ```
    This command will build the images for `SimpleApp` and `LoadBalancer` (if not already built) and then start the containers. By default, it starts one `SimpleApp` instance and the `LoadBalancer`.

3.  **Running multiple instances of SimpleApp:**
    To run a specific number of `SimpleApp` instances, use the `--scale` option:
    ```bash
    docker-compose up --scale simpleapp=<count>
    ```
    Replace `<count>` with the desired number of `SimpleApp` instances (e.g., `docker-compose up --scale simpleapp=3` to run 3 instances). The `LoadBalancer` will distribute requests among these instances.

    The `LoadBalancer` will be accessible at `http://localhost:8080`. `SimpleApp` instances will register themselves with the `LoadBalancer`.

### Stopping the Services
To stop and remove the containers, networks, and volumes created by `docker-compose up`, run:
```bash
docker-compose down
```
If you want to remove the images as well, you can do so manually using `docker rmi <image_id_or_name>`.

### Automated Image Publishing
This repository is configured with a GitHub Actions workflow to automatically build and publish Docker images to Docker Hub when a new version tag (e.g., `v1.0.0`, `v0.1.0`) is pushed.

The images are published to the following locations:
-   **SimpleApp:** `atjio/atjio-exercise-app:<tagname>`
-   **LoadBalancer:** `atjio/atjio-exercise-loadbalancer:<tagname>`

You can pull these pre-built images instead of building them locally if you prefer, for example:
```bash
docker pull atjio/atjio-exercise-app:v1.0.0
docker pull atjio/atjio-exercise-loadbalancer:v1.0.0
```
(Replace `v1.0.0` with the desired tag.)
