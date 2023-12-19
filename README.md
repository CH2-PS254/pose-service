# Pose Service

Golang service to manage pose data

## Setup

1. Copy `.env.example` to `.env`

   ```sh
   cp .env.example .env
   ```

2. Fill `.env` based on your credentials

3. Install package

   ```sh
   go get .
   ```

4. Run service

   ```sh
   go run .
   ```

## Deployments

### App Engine

1. Create `app.yaml` file in project root

   ```yaml
   runtime: go121

   env_variables:
     GIN_MODE: release
     CLOUDSQL_CONNECTION_NAME: PROJECT_ID:REGION_ID:INSTANCE_ID
     CLOUDSQL_USER: postgres
     CLOUDSQL_PASSWORD: ""
     CLOUDSQL_DATABASE_NAME: postgres
     JWT_SECRET: ""
   ```

2. Deploy to App Engine using gcloud CLI

   ```sh
   gcloud app deploy
   ```
