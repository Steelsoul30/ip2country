# ip2country

`ip2country` is a Go-based service that provides country information based on IP addresses. It uses a combination of local and remote data sources to retrieve the necessary information.
This is a POC only. Do not use this code in any production environment. This is for educational purposed only.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Running the Project](#running-the-project)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [TODO](#todo)


## Prerequisites

- Go 1.23 or higher
- Docker (optional, for containerization)
- A data source for IP information. This can be a local file or an API endpoint. See zip file for an example of a local file.
 
## Configuration

The configuration for the project is managed through a `config.yaml` file together with environment variables. Capitalized variables can be overriden if set as environment variables
Below is an example configuration:

```yaml
# config.yaml

db:
  - host: "db/geolite2.zip"
    name: "local"
  - host: "https://tools.keycdn.com/geo.json?host"
    name: "api"

logger:
  level: "info"
  serviceName: "ip2country"
  serviceVersion: "0.0.1"

ACTIVE_DATA_STORE: "local"
RATE_LIMIT: 1
BURST_LIMIT: 1
port: 8080
isDebug: false
```
### Configuration Options

- \`db\`: List of data sources.
  - \`host\`: The location of the data source.
  - \`name\`: The name of the data source.
- \`logger\`: Logger configuration.
  - \`level\`: Logging level (e.g., "info", "debug").
  - \`serviceName\`: Name of the service.
  - \`serviceVersion\`: Version of the service.
- \`ACTIVE_DATA_STORE\`: The active data store to use ("local" or "api").
- \`RATE_LIMIT\`: The rate limit for requests.
- \`BURST_LIMIT\`: The burst limit for requests.
- \`port\`: The port on which the service will run.
- \`isDebug\`: Enable or disable debug mode.

## Running the Project

To run the project, follow these steps:

1. **Clone the repository:**

   ```sh
   git clone https://github.com/steelsoul30/ip2country.git
   cd ip2country
   ```

2. **Install dependencies:**

   ```sh
   go mod tidy
   ```

3. **Create and configure \`config.yaml\`:**

   Create a \`config.yaml\` file in the root directory of the project and configure it as shown in the [Configuration](#configuration) section.

4. **Run the service:**

   ```sh
   go run main.go run
   ```
  The run command is necessary.
  The service will start and listen on the port specified in the \`config.yaml\` file.
  
  You can also use the database prebuilding command to save some startup time in case you are using a local database
  ```sh
   go run main.go create-db
   ```
  For help:
  ```sh
   go run main.go
   ```

### Note:
To run with a local db, a geolite2.zip should be present with the following files:
- GeoLite2-City-Blocks-IPv4.csv
- GeoLite2-City-Locations-en.csv
You can use the file already present

## Usage
Once the service is running, you can make HTTP requests to it to get country information based on IP addresses. For example:
    
    ```sh
    curl http://localhost:8080/ip2country?ip=2.22.233.255
    ```
   
## Project Structure
- \`cmd/\`: Contained the main commands. Effectively these are the entry points
- \`config/\`: Contains configuration-related code.
-  \`db/\`: Contains the database. This is where the zip file should go.
- \`internal/\`: Contains the core logic of the application.
  - \`ip2country/\`: Contains the main functionality of the service.
    - \`handler/\`: Contains HTTP handlers.
    - \`store/\`: Contains data store implementations.
  - \`middleware/\`: Contains middleware for the service.
- \`pkg/\`: Contains shared packages.
- \`main.go\`: The entry point of the application.

## TODO

- [ ] Add more data sources for IP information, such as a relational database
- [ ] Implement caching for IP lookups. Specifically, a serializable radix trie
- [ ] Add more detailed logging.
- [ ] Improve error handling and reporting.
- [ ] Write unit tests for all components.
- [ ] Add support for IPv6 addresses.
- [ ] Create a Dockerfile for containerization.
- [ ] Set up CI/CD pipeline for automated testing and deployment.
