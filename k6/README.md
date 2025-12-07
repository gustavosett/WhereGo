# WhereGo Load Testing

This directory contains k6 performance tests designed to validate the reliability and speed of WhereGo under high-concurrency scenarios, simulating multiple backend services querying the API simultaneously.

## Prerequisites

- [Docker](https://www.docker.com/)
- [k6](https://k6.io/docs/get-started/installation/)

## 1. Start the Environment

To ensure a controlled environment for benchmarks and accurate results, we strictly limit the Docker container to **8 vCPUs**.

**Run WhereGo:**

```bash
docker run --cpus="8" -p 8080:8080 --rm gustavosett/wherego:latest
```

*Note: Ensure port 8080 is free on your machine before running.*

## 2. Run the Load Test

We use a wrapper script (`run-with-report.js`) to execute the test and generate a visual HTML report.

**Run with k6:**

```bash
k6 run run-with-report.js
```

## 3. View Results

Once the test completes:

1. Look for the `index.html` file generated in this directory.
2. Open it in your web browser.
3. Use the generated charts and metrics to validate the "Ultra-fast" claims (e.g., verifying p95 latency < 50ms).

### Test Scenarios

The script (`load-test.js`) simulates the following backend traffic patterns:

1. **Ramp Up:** Scales to 50 concurrent backend consumers.
2. **Stability:** Holds 50 concurrent consumers for 1 minute.
3. **Stress:** Scales to 100 concurrent consumers.
4. **Peak:** Holds 100 concurrent consumers for 1 minute to test sustained peak load.
5. **Ramp Down:** Gracefully closes connections.
