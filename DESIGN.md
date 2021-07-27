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

