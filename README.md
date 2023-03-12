# ohlc
## Setup

1.    Install Docker on your local machine.

2.    Clone this repository:


```

git clone https://github.com/yourusername/yourproject.git

```

3.  Create a .env file in the root directory of the project and set the following environment variables:


```

DATABASE_URL=DATABASE_URL

```

4.  Build the Docker image:


```

docker build -t ohlc-image .

```

5.  Run the Docker container:

```

    docker run --env-file .env -p 8080:8080 ohlc-image

```

6.    Access the API at http://localhost:8080/<api>.

## API

The API has two endpoint:
## POST /update

Updates the OHLCData table in the database with data from a CSV file located at a specified URL or file address.

Request Body

```

POST /update HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 31

{
        "url": "data.csv"
}

```

Response Body

```json

{
  "response": "Bulk insert from CSV successful"
}
```

## GET /query

Query the data from database and Pagination and search using query string for data must be available.

This API accepts the following query string parameters:

*    page: the page number of the results to return (default is 1)
*    limit: the maximum number of results to return per page (default is 10)
*    search: a search term to filter the results by (optional)

The API queries the database for OHLC data and builds a JSON response containing the requested page of data. The search parameter is used to filter the data by symbol using a WHERE clause in the SQL query string.

Request Body

```

GET /query?limit=10&search=BTCUSDT&page=1 HTTP/1.1
Host: localhost:8080

```

Response Body

```json

[
    {
        "close": "42146.06",
        "high": "42148.32",
        "low": "42120.82",
        "open": "42123.29",
        "symbol": "BTCUSDT",
        "unix": "1.64472E+12"
    },
    {
        "close": "42123.3",
        "high": "42126.32",
        "low": "42113.07",
        "open": "42113.08",
        "symbol": "BTCUSDT",
        "unix": "1.64472E+12"
    },
    {
        "close": "42113.07",
        "high": "42130.23",
        "low": "42111.01",
        "open": "42120.8",
        "symbol": "BTCUSDT",
        "unix": "1.64472E+12"
    },
    {
        "close": "42120.8",
        "high": "42123.31",
        "low": "42102.22",
        "open": "42114.47",
        "symbol": "BTCUSDT",
        "unix": "1.64472E+12"
    },
    {
        "close": "42114.48",
        "high": "42148.24",
        "low": "42114.04",
        "open": "42148.23",
        "symbol": "BTCUSDT",
        "unix": "1.64472E+12"
    }
]

```

## License

This project is licensed under the MIT License.