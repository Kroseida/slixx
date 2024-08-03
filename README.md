<img src=".github/assets/logo.png" width="200">

# Slixx


![GitHub commit activity](https://img.shields.io/github/commit-activity/w/kroseida/slixx)
| ![GitHub contributors](https://img.shields.io/github/contributors/kroseida/slixx)
![GitHub License](https://img.shields.io/github/license/kroseida/slixx)
| ![GitHub repo size](https://img.shields.io/github/repo-size/kroseida/slixx)
| ![GitHub top language](https://img.shields.io/github/languages/top/kroseida/slixx)
![GitHub commits since tagged version (branch)](https://img.shields.io/github/commits-since/kroseida/slixx/0.0.1/develop)
| ![GitHub forks](https://img.shields.io/github/forks/kroseida/slixx)



Slixx is an advanced backup application currently in the alpha stage of development. It aims to provide robust and
flexible file backup solutions using various protocols and methods. Our primary focus is to deliver a reliable tool for
data protection and recovery with a user-friendly interface and powerful features.

## Features

- **Multi-Protocol Support**: Backup files using different protocols.
    - **FTP**
    - **SFTP**
- **Backup Methods**: Choose the best backup method for your needs.
    - **Copy**: Simple file copying.
    - **Incremental**: Efficient incremental backups to save space and time.
- **Modern Technology Stack**:
    - **Frontend**: Built with Vue.js for a responsive and intuitive user interface.
    - **Backend**: Powered by GoLang for high performance and concurrency.
    - **API**: Utilizes GraphQL for flexible and efficient data querying.

## Current Status

Slixx is currently in the alpha stage and is under active development. Features and functionalities are being
continuously added and refined. We welcome any feedback and contributions from the community to help improve Slixx.

## Getting Started

To get started with Slixx, you can use Docker Compose to set up the necessary services quickly.

### Prerequisites

- Docker
- Docker Compose

## Docker Images
Slixx is available as a set of Docker images for easy deployment. You can find the images on Docker Hub.

Supervisor: https://hub.docker.com/r/kroseida/slixx.supervisor

Satellite: https://hub.docker.com/r/kroseida/slixx.satellite

### Quickstart Docker Compose

   ```yaml
version: '3'
services:
  supervisor:
    image: kroseida/slixx.supervisor:latest
    container_name: supervisor
    volumes:
      - ~/supervisor/data:/app/data
      - ~/supervisor/log:/app/log
    ports:
      - "3030:3030/tcp"
    restart: unless-stopped

  satellite1:
    image: kroseida/slixx.satellite:latest
    container_name: satellite1
    environment:
      - SATELLITE_AUTH_TOKEN={AUTH TOKEN OF SATELLITE1}
    ports:
      - "9623:9623/tcp"
    restart: unless-stopped

  satellite2:
    image: kroseida/satellite:latest
    container_name: satellite2
    environment:
      - SATELLITE_AUTH_TOKEN={AUTH TOKEN OF SATELLITE2}
    ports:
      - "9624:9623/tcp"
    restart: unless-stopped
```

   ```bash
    docker-compose up -d
```

### Usage
Once the services are up and running, you can access the Supervisor interface via http://localhost:3030. Use the provided authentication tokens to configure and manage your backups, satellites, jobs ...

### Contributing
We welcome contributions from the community! If you'd like to contribute, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them.
4. Push your changes to your fork.
5. Submit a pull request to the main repository.
6. I will review your changes and merge them if they look good.

### License
Slixx is licensed under the GPL-2.0 License. See the LICENSE file for more details.

### Acknowledgements
We'd like to thank all the contributors and users for their support and feedback. Together, we can make Slixx a powerful tool for everyone.

