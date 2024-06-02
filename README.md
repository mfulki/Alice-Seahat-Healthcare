# Alice-Seahat-Healthcare
This Code has immigrated from GitLab.
The code was made by : 
1. Muhammad Luthfi Taufiqurrahman
2. Muhammad Fulki Aslam Adriswan 
3. Obie Krisnanto
4. Evan Susanto 

## Initiate DB

1. go to the base path
2. sudo -u postgres psql
3. \c postgres
4. \i database/sql/ddl.sql
5. \i database/sql/dml.sql

## How to Usage

1. Install go dependencies with:

```bash
go mod tidy
```

2. Then, copy `.env.example` to `.env`

```bash
cp .env.example .env
```

3. Then, fill your local environment to .env
4. Last, run the server.

```bash
go run .
```

Optional, you can start server with Makefile configuration. Ensure, you have been installed `nodemon`.

```bash
make dev
```
