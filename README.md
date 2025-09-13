# WAMP

> [!WARNING]
> THIS WAMP STACK IS A WORK IN PROGRESS! ANYTHING CAN CHANGE AT ANY MOMENT WITHOUT ANY NOTICE! USE THIS WAMP STACK AT YOUR OWN RISK!

Yet another tool for managing a local WAMP (Windows, Apache, MySQL, PHP) development environment on Windows.

## Features

- **Service Management:** Start and stop Apache and MySQL services.
- **Virtual Host Management:** Easily add and remove virtual hosts (sites).
- **PHP Version Management:** Install and switch between different PHP versions.
- **Automated Configuration:** Automatically configures Apache and hosts file.

## Installation

### Prerequisites

- [Go](https://golang.org/doc/install) (for building from source)

### Building from Source

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/aziyan99/wamp.git
    cd wamp
    ```

2.  **Build the executable:**
    Use the `built.bat` script to build the application.

    - **For a development build:**
      ```sh
      built.bat dev
      ```
      This will create `build\wamp-dev.exe`.

    - **For a production build:**
      ```sh
      built.bat prod
      ```
      This will create `build\wamp.exe`.

3.  **Initialize the environment:**
    This will create the necessary directories (`bin`, `www`, `tmp`).
    ```sh
    build\wamp.exe init
    ```

The compiled `wamp.exe` is a standalone executable. You can move it to any directory on your system (e.g., `C:\wamp\`). All the files and folders that `wamp.exe` creates and manages (like `bin`, `www`, `tmp`, and configuration files) will be contained within the same directory as the executable. This makes the entire WAMP environment portable. If it is still not working, sometimes stop-start the apache again or close-open the browser fix it.


## Usage

All commands are run from the directory containing `wamp.exe`.

### General Commands

- **Initialize Application:**
  ```sh
  wamp.exe init
  ```

- **Install Services:**
  ```sh
  wamp.exe install
  ```

- **Uninstall Services:**
  ```sh
  wamp.exe uninstall
  ```

### Service Management

- **Start Apache:**
  ```sh
  wamp.exe apache start
  ```

- **Stop Apache:**
  ```sh
  wamp.exe apache stop
  ```

- **Start MySQL:**
  ```sh
  wamp.exe mysql start
  ```

- **Stop MySQL:**
  ```sh
  wamp.exe mysql stop
  ```

### Site (Virtual Host) Management

- **Add a Site:**
  ```sh
  wamp.exe site add <site-name> [--php <version>] [--ssl]
  ```
  - `<site-name>`: The desired local domain (e.g., `my-project.test`).
  - `--php` (or `-p`): Specify the PHP version to use (e.g., `php-8.3`). Defaults to `php-8.3`.
  - `--ssl` (or `-s`): Enable SSL. Defaults to `false`.

  **Example:**
  ```sh
  wamp.exe site add my-laravel-app.test --php php-8.2 --ssl true
  ```

- **Remove a Site:**
  ```sh
  wamp.exe site rm <site-name>
  ```

### PHP Management

- **Install a PHP Version:**
  Downloads and installs a specific version of PHP.
  ```sh
  wamp.exe php install <full-version-name>
  ```
  **Example:**
  ```sh
  wamp.exe php install php-8.4.9-nts-Win32-vs17-x64
  ```

## Configuration

The active versions of Apache and MySQL are configured in the `wamp.ini` file, which is created after running the `install` command.

```ini
[apache]
active = apache-2.4

[mysql]
active = mysql-8.0
```

Sometime SSL cert not working, to solve that clear the SSL cache in: Control Panels > Internet Options > Content > Clear SSL State

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.