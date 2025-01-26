# BugBounty-Automation-Tool

A TUI-based tool for automating bug bounty and penetration testing tasks, built in Golang.

![GitHub Stars](https://img.shields.io/github/stars/yourusername/bugbounty-automation-tool)
![GitHub Forks](https://img.shields.io/github/forks/yourusername/bugbounty-automation-tool)
![GitHub Issues](https://img.shields.io/github/issues/yourusername/bugbounty-automation-tool)
![GitHub License](https://img.shields.io/github/license/yourusername/bugbounty-automation-tool)

## Overview

This project is designed to streamline the workflow of bug bounty hunters and penetration testers by automating the execution of various security testing tools. It provides a user-friendly terminal interface to select and run tools, view real-time output, and manage results efficiently.

## Features

- **Multi-Tool Support**: Integrates with popular tools across different categories:
  - **Reconnaissance**: assetfinder, subfinder, amass, dnsx, httpx, hakrawler, gospider, sniper
  - **Scanning**: nmap, rustscan, naabu, nuclei, sqlmap, arjun, xssstrike, ffuf, paramspider, cariddi
  - **Crawling**: hakrawler, gospider, sniper
  - **Exploitation**: metasploit

- **Real-Time Output**: View the output of commands as they execute.

- **Input Validation**: Ensures that required inputs like target domains and output directories are valid before executing tasks.

- **Piping Functionality**: Allows users to pipe the output of one tool to another, enabling complex workflows.

- **Customizable Arguments**: Users can specify custom arguments for each tool.

- **Output Management**: Results can be saved to a specified directory.

## Installation

### Prerequisites

- **Golang**: Ensure you have Go installed on your system.
- **Tools**: Install the required tools (e.g., nmap, rustscan, naabu, etc.) and ensure they are accessible from your PATH.

### Setup

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/yourusername/bugbounty-automation-tool.git
   cd bugbounty-automation-tool
   ```

2. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

3. **Configure Tools**:
   - Update `config.json` with the correct paths to your installed tools.
   - Modify `target_domains.json` to include your target domains.

### Running the Application

1. **Execute the Application**:
   ```bash
   go run main.go
   ```

2. **Usage**:
   - **Navigation**: Use arrow keys to navigate the task list.
   - **Selection**: Press Enter to select a task.
   - **Arguments**: Enter custom arguments in the provided text box.
   - **Target Domains**: Input target domain(s) in the domain input field.
   - **Output Directory**: Specify where results should be saved.
   - **Piping**: Toggle piping mode with the Pipe button to chain tool executions.

## Configuration

### config.json

This file defines the tools and their paths. Example:

```json
{
    "tools": {
        "recon": {
            "assetfinder": "/usr/local/bin/assetfinder",
            "subfinder": "/usr/local/bin/subfinder",
            "amass": "/usr/local/bin/amass",
            "dnsx": "/usr/local/bin/dnsx",
            "httpx": "/usr/local/bin/httpx",
            "hakrawler": "/usr/local/bin/hakrawler",
            "gospider": "/usr/local/bin/gospider",
            "sniper": "/usr/local/bin/sniper"
        },
        "scanning": {
            "nmap": "/usr/bin/nmap",
            "rustscan": "/usr/local/bin/rustscan",
            "naabu": "/usr/local/bin/naabu",
            "nuclei": "/usr/local/bin/nuclei",
            "sqlmap": "/usr/local/bin/sqlmap",
            "arjun": "/usr/local/bin/arjun",
            "xssstrike": "/usr/local/bin/xssstrike",
            "ffuf": "/usr/local/bin/ffuf",
            "paramspider": "/usr/local/bin/paramspider",
            "cariddi": "/usr/local/bin/cariddi"
        },
        "crawling": {
            "hakrawler": "/usr/local/bin/hakrawler",
            "gospider": "/usr/local/bin/gospider",
            "sniper": "/usr/local/bin/sniper"
        },
        "exploit": {
            "metasploit": "/usr/local/bin/msfconsole"
        }
    },
    "bbot": {
        "target": "example.com"
    }
}
```

### target_domains.json

This file lists the target domains. Example:

```json
{
    "domains": [
        "example.com",
        "sub.example.com",
        "api.example.com",
        "admin.example.com"
    ]
}
```

## Usage

1. **Select a Task**:
   - Navigate through the list of available tools using arrow keys.
   - Press Enter to select a tool.

2. **Configure Task**:
   - **Arguments**: Enter any additional arguments for the selected tool.
   - **Target Domains**: Input the target domain(s) separated by commas.
   - **Output Directory**: Specify the directory to save results.

3. **Execute Task**:
   - Click the "Start" button to execute the selected tool with the configured arguments.

4. **Piping Functionality**:
   - Enable piping mode by clicking the "Pipe" button.
   - Select the first tool and configure it.
   - The output of the first tool will be piped to the next selected tool.

## Known Issues and Limitations

- **Basic Error Handling**: The current implementation includes fundamental error handling. More robust error handling and recovery mechanisms are planned for future updates.
- **Limited Piping Complexity**: The piping functionality is in its early stages and supports simple tool chaining. More complex workflows will be added in future releases.

## Contributing

Contributions are welcome! If you have any feature requests, bug reports, or improvements, please:

1. **Submit an Issue**: Describe your request or issue on the [GitHub Issues](https://github.com/yourusername/bugbounty-automation-tool/issues) page.
2. **Fork the Repository**: Create your feature branch (`git checkout -b feature/your-feature-name`).
3. **Commit Changes**: Implement your changes and commit them.
4. **Push to the Branch**: Push your changes to your forked repository.
5. **Open a Pull Request**: Submit a Pull Request against the `main` branch of this repository.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the need for efficient automation in bug bounty and penetration testing workflows.
- Built using [Golang](https://golang.org/) and [Charmbraceful](https://github.com/charmbraceful/charmbraceful) for the TUI.

---

**BugBounty-Automation-Tool** is a powerful utility designed to simplify and accelerate your security testing processes. With its extensible architecture and user-friendly interface, it's an essential tool for any security professional.# Goforgold
