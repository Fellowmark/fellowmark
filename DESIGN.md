# Design Specifications

## NFRs

### Admin scenarios
1. Login from admin, student and staff
2. Signup of individual user groups
3. Create and retrieve module details
4. Enrollment (bulk inserts) and supervision registration
5. Reset password

### Assignment scenarios
1. Create assignments, questions and pairings, rubrics.
2. Students upload upload pdfs/texts.

### Review scenarios
1. Students rubric markings.
2. Students submitting appeals
3. Staff review appeals
4. Staff moderate rubrics markings.

## FRs
### API
1. CRUD opertions and email service
2. DB migrations
3. Security
4. Not for public client extension (only frontend service).

### Postgres
1. Soft delete
2. Scheduled backups

### File server
1. Store data with concurrent access
2. Scheduled backups
3. Links to files in the file system are stored in the postgres
4. Could be merged with API for a bigger monolith or replaced by Postgres 
![blobs](https://www.enterprisedb.com/postgres-tutorials/postgresql-toast-and-working-blobsclobs-explained)

### Frontend
1. Single SPA
2. Material UI
3. Admin dashboards

## System Design
For the initial production phase this architecture will deployed on a single
VM instance. Reverse proxy allows access to individual services through
a single host IP.
```
Reverse Proxy (NGINX)
______________________
| ______________      |
| | Postgres   |      |
| |_____________      |
| | API        |      |
| |_____________      |
| | File Server|      |
| |_____________      |
| | Frontend   |      |
| --------------      |
|_____________________|
```

## API code structure
.env variable to store config variables at root. API is a monolith but multithreaded.

### docker-compose
1. API-DEV profile to start api and db with dev locals
2. API profile to start api and db with production locals

### main package
Links to:
1. db package
2. route builder packages, `student`, `staff`, `admin`
3. loggers package (persistant logging in production)

### db package
1. Open connection with db
2. Migrate schemas
3. Connection pool
4. Return pointer to `*gorm.DB` for pooled interfacing

### student package
1. StudentRoute struct that stores injectable dependencies
2. RouteBuilders
3. Route specific middleware
4. Package is used to create routes students use to send and retrieve data.

### staff package
1. StaffRoute struct that stores injectable dependencies
2. RouteBuilders
3. Route specific middleware
4. Package is used to create routes staff use to send and retrieve data.

### admin package
1. AdminRoute struct that stores injectable dependencies
2. RouteBuilders
3. Route specific middleware
4. Package is used to create routes admin use to send and retrieve data.

### utils package
1. Auth utils for functions and http handlers used by auth components of route based packages
1. Middleware utils for reusable middlewares across all route packages to reduce code redundancy

## Frontend to API Interfacing
This design spec uses F2A interfacing for documentation: F2A is a custom design
spec that has three components:
1. Action: NFR scenario to be handled through series of Request, Endpoint and Responses called steps.
2. Steps: Series of RER to perform an action.
  - Request: Request from frontend with data schema (can be `NULL`)
  - Response: Response if any from API with data schema (can be `NULL`)
  - Reaction: Reaction of the frontend to the data from the API (can be `NULL`)

This doc model is a higher-level aggregation of specs like OpenAPI and gRPC.

## Assignment scenarios
### Create assignments, questions and pairings, rubric
#### Action
Staff creates all dependencies necessary for student submissions and reviews
#### Steps
<table>
  <tr>
    <th>Request</th>
    <th>Response</th>
    <th>Reaction</th>
  </tr>
  <tr>
  <td>
    POST request to create assignment 
    <pre>
    { 
      name: string, 
      module_id: int 
    }
    </pre>
  </td>
  <td>
    API returns id of new assignment 
    <pre>
    { 
      id: int 
    }
    </pre>
  </td>
  <td>
    Frontend logs assignment id to post question
  </td>
</tr>
  <tr>
  <td>
    POST request to create question
    <pre>
    { 
      assign_id: int, 
      question_number: int,
      question_text: string 
    }
    </pre>
  </td>
  <td>
    API returns id of question
    <pre>
    { 
      id: int 
    }
    </pre>
  </td>
  <td>
    Frontend logs question id to post rubric criterias
  </td>
</tr>
  <tr>
  <td>
    POST request to rubric criteria
    <pre>
    { 
      question_id: int, 
      criteria: string,
      description: string,
      min_mark: int,
      max_mark: int
    }
    </pre>
  </td>
  <td>
    API returns id of rubric criteria
    <pre>
    { 
      id: int 
    }
    </pre>
  </td>
  <td>
    Frontend keep sending rubric requests until rubric is complete
  </td>
</tr>
</table>