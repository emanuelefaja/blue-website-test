---
title: Location Custom Field
description: Create location fields to store geographic coordinates for records
category: Custom Fields
---

Location custom fields store geographic coordinates (latitude and longitude) for records. They support precise coordinate storage, geospatial queries, and efficient location-based filtering.

## Basic Example

Create a simple location field:

```graphql
mutation CreateLocationField {
  createCustomField(input: {
    name: "Meeting Location"
    type: LOCATION
    projectId: "proj_123"
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create a location field with description:

```graphql
mutation CreateDetailedLocationField {
  createCustomField(input: {
    name: "Office Location"
    type: LOCATION
    projectId: "proj_123"
    description: "Primary office location coordinates"
  }) {
    id
    name
    type
    description
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the location field |
| `type` | CustomFieldType! | ✅ Yes | Must be `LOCATION` |
| `description` | String | No | Help text shown to users |

## Setting Location Values

Location fields store latitude and longitude coordinates:

```graphql
mutation SetLocationValue {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    latitude: 40.7128
    longitude: -74.0060
  })
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record to update |
| `customFieldId` | String! | ✅ Yes | ID of the location custom field |
| `latitude` | Float! | ✅ Yes | Latitude coordinate (-90 to 90) |
| `longitude` | Float! | ✅ Yes | Longitude coordinate (-180 to 180) |

## Coordinate Validation

### Valid Ranges

| Coordinate | Range | Description |
|------------|-------|-------------|
| Latitude | -90 to 90 | North/South position |
| Longitude | -180 to 180 | East/West position |

### Example Coordinates

| Location | Latitude | Longitude |
|----------|----------|-----------|
| New York City | 40.7128 | -74.0060 |
| London | 51.5074 | -0.1278 |
| Sydney | -33.8688 | 151.2093 |
| Tokyo | 35.6762 | 139.6503 |
| São Paulo | -23.5505 | -46.6333 |

## Creating Records with Location Values

When creating a new record with location data:

```graphql
mutation CreateRecordWithLocation {
  createTodo(input: {
    title: "Site Visit"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "location_field_id"
      value: "40.7128,-74.0060"  # Format: "latitude,longitude"
    }]
  }) {
    id
    title
    customFields {
      id
      customField {
        name
        type
      }
      latitude
      longitude
    }
  }
}
```

### Input Format for Creation

When creating records, location values use comma-separated format:

| Format | Example | Description |
|--------|---------|-------------|
| `"latitude,longitude"` | `"40.7128,-74.0060"` | Standard coordinate format |
| `"51.5074,-0.1278"` | London coordinates | No spaces around comma |
| `"-33.8688,151.2093"` | Sydney coordinates | Negative values allowed |

## Response Fields

### TodoCustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field value |
| `customField` | CustomField! | The custom field definition |
| `latitude` | Float | Latitude coordinate |
| `longitude` | Float | Longitude coordinate |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last modified |

## Important Limitations

### No Built-in Geocoding

Location fields store only coordinates - they do **not** include:
- Address-to-coordinates conversion
- Reverse geocoding (coordinates-to-address)
- Address validation or search
- Integration with mapping services
- Place name lookup

### External Services Required

For address functionality, you'll need to integrate external services:
- **Google Maps API** for geocoding
- **OpenStreetMap Nominatim** for free geocoding
- **MapBox** for mapping and geocoding
- **Here API** for location services

### Example Integration

```javascript
// Client-side geocoding example (not part of Blue API)
async function geocodeAddress(address) {
  const response = await fetch(
    `https://maps.googleapis.com/maps/api/geocode/json?address=${encodeURIComponent(address)}&key=${API_KEY}`
  );
  const data = await response.json();
  
  if (data.results.length > 0) {
    const { lat, lng } = data.results[0].geometry.location;
    
    // Now set the location field in Blue
    await setTodoCustomField({
      todoId: "todo_123",
      customFieldId: "location_field_456",
      latitude: lat,
      longitude: lng
    });
  }
}
```

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create location field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update location field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Set location value | Standard record edit permissions |
| View location value | Standard record view permissions |

## Error Responses

### Invalid Coordinates
```json
{
  "errors": [{
    "message": "Invalid coordinates: latitude must be between -90 and 90",
    "extensions": {
      "code": "CUSTOM_FIELD_VALUE_PARSE_ERROR"
    }
  }]
}
```

### Invalid Longitude
```json
{
  "errors": [{
    "message": "Invalid coordinates: longitude must be between -180 and 180",
    "extensions": {
      "code": "CUSTOM_FIELD_VALUE_PARSE_ERROR"
    }
  }]
}
```

## Best Practices

### Data Collection
- Use GPS coordinates for precise locations
- Validate coordinates before storing
- Consider coordinate precision needs (6 decimal places ≈ 10cm accuracy)
- Store coordinates in decimal degrees (not degrees/minutes/seconds)

### User Experience
- Provide map interfaces for coordinate selection
- Show location previews when displaying coordinates
- Validate coordinates client-side before API calls
- Consider timezone implications for location data

### Performance
- Use spatial indexes for efficient queries
- Limit coordinate precision to needed accuracy
- Consider caching for frequently accessed locations
- Batch location updates when possible

## Common Use Cases

1. **Field Operations**
   - Equipment locations
   - Service call addresses
   - Inspection sites
   - Delivery locations

2. **Event Management**
   - Event venues
   - Meeting locations
   - Conference sites
   - Workshop locations

3. **Asset Tracking**
   - Equipment positions
   - Facility locations
   - Vehicle tracking
   - Inventory locations

4. **Geographic Analysis**
   - Service coverage areas
   - Customer distribution
   - Market analysis
   - Territory management

## Integration Features

### With Lookups
- Reference location data from other records
- Find records by geographic proximity
- Aggregate location-based data
- Cross-reference coordinates

### With Automations
- Trigger actions based on location changes
- Create geofenced notifications
- Update related records when locations change
- Generate location-based reports

### With Formulas
- Calculate distances between locations
- Determine geographic centers
- Analyze location patterns
- Create location-based metrics

## Limitations

- No built-in geocoding or address conversion
- No mapping interface provided
- Requires external services for address functionality
- Limited to coordinate storage only
- No automatic location validation beyond range checking

## Related Resources

- [Custom Fields Overview](/custom-fields/list-custom-fields) - General concepts
- [Google Maps API](https://developers.google.com/maps) - Geocoding services
- [OpenStreetMap Nominatim](https://nominatim.org/) - Free geocoding
- [MapBox API](https://docs.mapbox.com/) - Mapping and geocoding services

Location custom fields allow you to store geographic coordinates (latitude and longitude) for records. They support coordinate validation, geographic filtering, and can be used to track addresses, delivery locations, or any location-based information.

## Basic Example

Create a simple location field:

```graphql
mutation CreateLocationField {
  createCustomField(input: {
    name: "Delivery Address"
    type: LOCATION
    projectId: "proj_123"
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create a location field with description:

```graphql
mutation CreateLocationField {
  createCustomField(input: {
    name: "Meeting Location"
    type: LOCATION
    projectId: "proj_123"
    description: "Physical location where the meeting will take place"
  }) {
    id
    name
    type
    description
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the location field |
| `type` | CustomFieldType! | ✅ Yes | Must be `LOCATION` |
| `description` | String | No | Help text shown to users |

## Setting Location Values

Location fields store latitude and longitude coordinates:

### Direct Coordinates

```graphql
mutation SetLocationCoordinates {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    latitude: 40.7128
    longitude: -74.0060
  })
}
```

### Using String Format

```graphql
mutation SetLocationFromString {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    locationString: "40.7128,-74.0060"
  })
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record to update |
| `customFieldId` | String! | ✅ Yes | ID of the location custom field |
| `latitude` | Float | No | Latitude coordinate (-90 to 90) |
| `longitude` | Float | No | Longitude coordinate (-180 to 180) |
| `locationString` | String | No | Comma-separated "latitude,longitude" string |

**Note**: Either provide both `latitude` and `longitude`, or use `locationString`. Both coordinates are required for a valid location.

## Coordinate Validation

### Latitude Range
- **Valid range**: -90 to 90 degrees
- **Examples**: 
  - `40.7128` (New York City)
  - `-33.8688` (Sydney)
  - `0.0` (Equator)

### Longitude Range
- **Valid range**: -180 to 180 degrees
- **Examples**:
  - `-74.0060` (New York City)
  - `151.2093` (Sydney)
  - `0.0` (Prime Meridian)

### Coordinate Formats

| Format | Example | Description |
|--------|---------|-------------|
| Decimal Degrees | `40.7128, -74.0060` | Most common format |
| String Format | `"40.7128,-74.0060"` | Comma-separated string |
| High Precision | `40.712776, -74.005974` | Up to 6 decimal places |

## Creating Records with Location Values

When creating a new record with location values:

```graphql
mutation CreateRecordWithLocation {
  createTodo(input: {
    title: "Site Visit"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "location_field_id"
      value: "40.7128,-74.0060"
    }]
  }) {
    id
    title
    customFields {
      id
      customField {
        name
        type
      }
      latitude
      longitude
      value
    }
  }
}
```

### Supported Input Formats

| Format | Example | Result |
|--------|---------|---------|
| Comma-separated | `"40.7128,-74.0060"` | Latitude: 40.7128, Longitude: -74.0060 |
| Space-separated | `"40.7128 -74.0060"` | Latitude: 40.7128, Longitude: -74.0060 |
| With labels | `"lat:40.7128,lng:-74.0060"` | Parsed to decimal coordinates |

## Response Fields

### TodoCustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field value |
| `customField` | CustomField! | The custom field definition |
| `latitude` | Float | Latitude coordinate |
| `longitude` | Float | Longitude coordinate |
| `value` | Object | Combined location object (see below) |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last modified |

### Value Object Structure

```json
{
  "value": {
    "latitude": 40.7128,
    "longitude": -74.0060
  }
}
```

**Note**: If either coordinate is missing, the value will be `null`.

## Geographic Filtering

Location fields support advanced geographic filtering:

```graphql
query FilterByLocation {
  todos(filter: {
    customFields: [{
      customFieldId: "location_field_id"
      polygon: [
        { latitude: 40.7000, longitude: -74.0200 },
        { latitude: 40.7200, longitude: -74.0200 },
        { latitude: 40.7200, longitude: -73.9800 },
        { latitude: 40.7000, longitude: -73.9800 }
      ]
      operator: WITHIN
    }]
  }) {
    id
    title
    customFields {
      latitude
      longitude
    }
  }
}
```

### Supported Operators

| Operator | Description |
|----------|-------------|
| `WITHIN` | Location is within the specified polygon |
| `OUTSIDE` | Location is outside the specified polygon |
| `NULL` | Location field is empty |
| `NOT_NULL` | Location field has coordinates |

### Polygon Structure

| Field | Type | Description |
|-------|------|-------------|
| `latitude` | Float! | Latitude of polygon point |
| `longitude` | Float! | Longitude of polygon point |

**Note**: Polygon must have at least 3 points to define a valid area.

## Location Display

The system provides intelligent location display:

| Scenario | Display Format |
|----------|----------------|
| Full coordinates | `40.7128, -74.0060` |
| High precision | `40.712776, -74.005974` |
| Missing coordinates | `Not set` |
| Invalid coordinates | `Invalid location` |

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create location field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update location field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Set location value | Standard record edit permissions |
| View location value | Standard record view permissions |
| Use geographic filtering | Standard record query permissions |

## Error Responses

### Invalid Latitude
```json
{
  "errors": [{
    "message": "Latitude must be between -90 and 90 degrees",
    "extensions": {
      "code": "INVALID_LATITUDE"
    }
  }]
}
```

### Invalid Longitude
```json
{
  "errors": [{
    "message": "Longitude must be between -180 and 180 degrees",
    "extensions": {
      "code": "INVALID_LONGITUDE"
    }
  }]
}
```

### Invalid Coordinate Format
```json
{
  "errors": [{
    "message": "Invalid coordinate format. Use 'latitude,longitude' or separate latitude/longitude fields",
    "extensions": {
      "code": "CUSTOM_FIELD_VALUE_PARSE_ERROR"
    }
  }]
}
```

### Missing Coordinates
```json
{
  "errors": [{
    "message": "Both latitude and longitude are required",
    "extensions": {
      "code": "INCOMPLETE_COORDINATES"
    }
  }]
}
```

## Best Practices

### Coordinate Precision
- Use 4-6 decimal places for most applications
- Higher precision (6+ places) for precise positioning
- Consider your accuracy requirements

### Data Entry
- Provide map picker in your UI
- Show coordinate validation in real-time
- Allow multiple input formats for flexibility

### Geographic Features
- Use polygon filtering for area-based queries
- Combine with other filters for complex searches
- Consider performance for large datasets

## Common Use Cases

1. **Field Operations**
   - Service locations
   - Delivery addresses
   - Site visits
   - Equipment locations

2. **Event Management**
   - Venue locations
   - Meeting points
   - Event addresses
   - Parking locations

3. **Asset Tracking**
   - Equipment locations
   - Facility management
   - Inventory tracking
   - Mobile asset monitoring

4. **Sales & Marketing**
   - Lead locations
   - Territory mapping
   - Store locations
   - Customer addresses

## Integration Features

### With Maps
- Display locations on interactive maps
- Calculate distances between points
- Route optimization
- Geofencing capabilities

### With Automations
- Trigger actions based on location changes
- Send notifications when entering/leaving areas
- Update status based on geographic rules

### With Lookups
- Find nearest locations
- Calculate geographic averages
- Group by geographic regions

### With Forms
- Map picker integration
- Address lookup and validation
- Coordinate conversion

## Limitations

- No built-in address-to-coordinate conversion
- No distance calculation functions
- Polygon filtering requires manual coordinate definition
- No support for 3D coordinates (elevation)
- Performance considerations for large polygon queries

## Related Resources

- [Custom Fields Overview](/custom-fields/list-custom-fields) - General custom field concepts
- [Automations API](/api/automations/index) - Create location-based automations
- [Forms API](/api/forms) - Include location fields in forms
- [Filtering API](/api/filtering) - Advanced location filtering techniques