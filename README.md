# Currency-conversion
Simple application that takes currencies rates from external third-party API and makes conversions.

# API
1. `/create` 
    - Create new project
    - returns API-Key for project

2. `/convert`
    - Convert value for symbol
    - authorized with project API-KEY in Authorization header
    - accepts json(e.g. {"symbol":"EUR/USD", value:5.65})
    - returns value converted
