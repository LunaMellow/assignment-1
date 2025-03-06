# Country Information REST API

## Overview
This REST web service provides information about countries, including general details, historical population data, and diagnostics for service availability. It is built using Go and retrieves data from the self-hosted versions of the **CountriesNow API** and **REST Countries API**.

## Endpoints
The API exposes three main endpoints:

### 1. `/countryinfo/v1/info/{:two_letter_country_code}{?limit=10}`
Provides general information about a country based on its **ISO 3166-2** two-letter code.

#### **Request**
- **Method:** `GET`
- **Path:** `/countryinfo/v1/info/{:two_letter_country_code}{?limit=10}`
- **Parameters:**
    - `two_letter_country_code` (required): Two-letter country ISO code.
    - `limit` (optional, default=10): Limits the number of cities returned, sorted alphabetically.

#### **Example Request:**
```
GET http://localhost:8080/countryinfo/v1/info/no
```

#### **Response**
- **Content-Type:** `application/json`
- **Status Code:** `200 OK` (or appropriate error codes)

**Example Response:**
```json
{
    "name":         "Norway",
    "continents":   ["Europe"],
    "population":   4700000,
    "languages":    {"nno":"Norwegian Nynorsk","nob":"Norwegian Bokmål","smi":"Sami"},
    "borders":      ["FIN","SWE","RUS"],
    "flag":         "https://flagcdn.com/w320/no.png",
    "capital":      "Oslo",
    "cities":       ["Abelvaer","Adalsbruk","Adland"]
}
```

---

### 2. `/countryinfo/v1/population/{:two_letter_country_code}{?limit={:startYear-endYear}}`
Retrieves population data for a given country, including the mean value over a specified period.

#### **Request**
- **Method:** `GET`
- **Path:** `/countryinfo/v1/population/{:two_letter_country_code}{?limit={:startYear-endYear}}`
- **Parameters:**
    - `two_letter_country_code` (required): Two-letter country ISO code.
    - `limit` (optional): Restricts the population data to the specified start and end years.

#### **Example Requests:**
```
GET http://localhost:8080/countryinfo/v1/population/no
GET http://localhost:8080/countryinfo/v1/population/no?limit=2010-2015
```

#### **Response**
- **Content-Type:** `application/json`
- **Status Code:** `200 OK` (or appropriate error codes)

**Example Response:**
```json
{
    "mean": 5044396,
    "values": [
        {"year":2010, "value":4889252},
        {"year":2011, "value":4953088},
        {"year":2012, "value":5018573},
        {"year":2013, "value":5079623},
        {"year":2014, "value":5137232},
        {"year":2015, "value":5188607}
    ]
}
```

---

### 3. `/countryinfo/v1/status/`
Provides the current status of the APIs this service depends on and uptime information.

#### **Request**
- **Method:** `GET`
- **Path:** `/countryinfo/v1/status/`

#### **Example Request:**
```
GET http://localhost:8080/countryinfo/v1/status/
```

#### **Response**
- **Content-Type:** `application/json`
- **Status Code:** `200 OK` (or appropriate error codes)

**Example Response:**
```json
{
    "countriesnowapi": 200,
    "restcountriesapi": 200,
    "version": "v1",
    "uptime": 3600
}
```

---

## Deployment
The service is deployed on **Render** and can be accessed at:
```
https://assignment-1-huwk.onrender.com/
```

### Local Setup & Development
1. Clone the repository:
   ```sh
   git clone https://github.com/LunaMellow/assignment-1
   cd assignment-1
   ```

2. Run the application locally:
   ```sh
   go run main.go
   ```
3. The service will be available at:
   ```
   http://localhost:8080 (default port)
   ```

### Testing
Use **Postman**, **curl**, or a browser to interact with the endpoints.

Example:
```sh
curl -X GET http://localhost:8080/countryinfo/v1/info/no
```

---

## Notes
- Use **self-hosted API endpoints** as specified in the assignment instead of the public APIs.
- Handle errors gracefully and provide meaningful messages.
- Reduce API calls by efficiently retrieving and caching data where possible.
- Follow RESTful API best practices.

## License
This project is for educational purposes as part of the **IDATG2206 Computer Vision course (Spring 2025)**.

---

**Author:** Luna Sofie Bergh <br>
**Course:** IDATG2206 Computer Vision (Spring 2025)

