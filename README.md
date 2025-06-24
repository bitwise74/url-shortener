# URL Shortener

A simple URL shortener service written in Go.

## Features

- Shorten long URLs to compact links
- Redirect short links to original URLs
- RESTful API endpoints
- In-memory storage (can be extended)

## Getting Started

### Prerequisites

- Go 1.18 or newer

### Installation

```bash
git clone https://github.com/bash06/url-shortener.git
cd url-shortener
go build -o url-shortener
```

### Usage

Start the server:

```bash
./url-shortener
```

By default, the server runs on `localhost:8080`.

### API Endpoints

- **POST /shorten**  
        Shorten a URL.  
        **Request:**  
        ```json
        { "url": "https://example.com" }
        ```
        **Response:**  
        ```json
        { "short_url": "http://localhost:8080/abc123" }
        ```

- **GET /{short_code}**  
        Redirects to the original URL.

## Configuration

You can configure the server port and other settings via environment variables or the `config.toml` file. Every available config option is explained in the `config.example.toml` file

## License

MIT License

---