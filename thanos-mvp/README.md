# ROS OCP Thanos PoC - Manual TSDB Block Uploads to MINIO

## Overview
This PoC demonstrates setting up **Thanos** with **MinIO** for long-term storage of Prometheus TSDB blocks. The setup includes:
- **MINIO** as the object storage backend.
- **Thanos Store Gateway** to access historical data from MinIO.
- **Thanos Querier** to aggregate and query data.

### 1. Start the containers
Run:
```sh
export MINIO_SECRET_KEY="minioadmin"
export MINIO_ACCESS_KEY="minioadmin"
docker-compose up
```
This will start MINIO, Thanos Store Gateway, Thanos Querier, and Kruize.

### 2. Verify MINIO and Buckets
MINIO runs at **http://localhost:9001** (Console UI).
- Access MINIO using:
  - **Username:** `minioadmin`
  - **Password:** `minioadmin`
- The bucket `rosocp-tsdb` should be created.

### 3. Query Data Using Thanos
Once blocks are uploaded:
- Thanos Querier runs at **http://localhost:19192**
- View on UI:
  ```sh
  http://localhost:19192/stores
  ```
thanos-store-gateway:19090 should show up here

### 4. Run the script
Manually upload TSDB blocks using the function:
```go
go run main.go
```
The script will upload the TSDB blocks, create a metric profile for Kruize and send a Bulk API request.

This assumes the TSDB blocks are located in `data/` and will be pushed to MinIO.
Get in touch with the contributor for TSDB blocks which would generate a kruize recommendation.


---

## Stopping & Cleaning Up
To stop and remove containers:
```sh
docker-compose down --remove-orphans
docker volume ls -qf dangling=true | xargs docker volume rm
```

---

## Response

```sh
Block <ULID> uploaded successfully!
Metric Profile Creation Status: 201 Created
Bulk API Request: {"filter":{"include":{"namespace":[],"workload":[],"containers":[],"labels":{}},"exclude":{"namespace":[],"workload":[],"containers":[],"labels":{}}},"time_range":{},"datasource":"thanos"}
Received ID: <JOB-ID>
```

## Notes
- Check Kruize Job details
```sh
curl localhost:8080/bulk\?job_id=<JOB-ID>&include=summary,experiments | jq
```
- In case there are some permission issues with executables, ideally there should be a command available in the stdout
```sh
docker-compose logs <service-name>
```
