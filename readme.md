# Virtual Ethiopia Digital Nation Prototype

This prototype demonstrates a blockchain-based digital nation with citizen registration and voting capabilities.

## Prerequisites

### Install Docker

1. For Windows/Mac:
   - Download Docker Desktop from [https://www.docker.com/products/docker-desktop](https://www.docker.com/products/docker-desktop)
   - Follow the installation wizard
   - Start Docker Desktop

2. For Ubuntu Linux:
```bash
# Update package index
sudo apt-get update

# Install prerequisites
sudo apt-get install ca-certificates curl gnupg lsb-release

# Add Docker's official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Set up the stable repository
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Add your user to the docker group (optional, to run docker without sudo)
sudo usermod -aG docker $USER
```

Verify Docker installation:
```bash
docker --version
docker compose version
```

## Running the Prototype

1. Extract the project:
```bash
clone https://github.com/ten2the6/virtual_ethiopia_dap.git
cd virtual_ethiopia_dap
```

2. Start the nodes:
```bash
docker-compose -f docker/docker-compose.yml up --build
```

## Testing the Digital Nation Features

### 1. Register Three Citizens

```bash
# Register Citizen 1
curl -X POST http://localhost:3001/citizens/register \
-H "Content-Type: application/json" \
-d '{
    "name": "John Doe",
    "dateOfBirth": "1990-01-01",
    "publicKey": "citizen1_key"
}'

# Register Citizen 2
curl -X POST http://localhost:3001/citizens/register \
-H "Content-Type: application/json" \
-d '{
    "name": "Jane Smith",
    "dateOfBirth": "1992-05-15",
    "publicKey": "citizen2_key"
}'

# Register Citizen 3
curl -X POST http://localhost:3001/citizens/register \
-H "Content-Type: application/json" \
-d '{
    "name": "Bob Johnson",
    "dateOfBirth": "1985-11-30",
    "publicKey": "citizen3_key"
}'

# Check registered citizens
curl http://localhost:3001/citizens
```

### 2. Approve Citizens
Note: Use the actual citizen IDs returned from the registration responses

```bash
# Approve Citizen 1
curl -X POST http://localhost:3001/citizens/approve \
-H "Content-Type: application/json" \
-d '{
    "citizenId": "CITIZEN1_ID_FROM_REGISTRATION",
    "approverKey": "GENESIS_ADMIN"
}'

# Approve Citizen 2
curl -X POST http://localhost:3001/citizens/approve \
-H "Content-Type: application/json" \
-d '{
    "citizenId": "CITIZEN2_ID_FROM_REGISTRATION",
    "approverKey": "GENESIS_ADMIN"
}'

# Approve Citizen 3
curl -X POST http://localhost:3001/citizens/approve \
-H "Content-Type: application/json" \
-d '{
    "citizenId": "CITIZEN3_ID_FROM_REGISTRATION",
    "approverKey": "GENESIS_ADMIN"
}'
```

### 3. Start an Election

```bash
curl -X POST http://localhost:3001/elections/start \
-H "Content-Type: application/json" \
-d '{
    "name": "Presidential Election 2024",
    "durationDays": 30
}'
```

### 4. Register Two Candidates

```bash
# Register Candidate 1
curl -X POST http://localhost:3001/elections/candidates \
-H "Content-Type: application/json" \
-d '{
    "name": "Alice Brown",
    "publicKey": "citizen1_key",
    "platform": "Innovation and Growth"
}'

# Register Candidate 2
curl -X POST http://localhost:3001/elections/candidates \
-H "Content-Type: application/json" \
-d '{
    "name": "Charlie Davis",
    "publicKey": "citizen2_key",
    "platform": "Sustainability and Education"
}'

# Check registered candidates
curl http://localhost:3001/elections/current
```

### 5. Cast Votes
Note: Use the actual candidate IDs from the candidate registration responses

```bash
# Vote from Citizen 1
curl -X POST http://localhost:3001/elections/vote \
-H "Content-Type: application/json" \
-d '{
    "citizenPublicKey": "citizen1_key",
    "candidateId": "CANDIDATE2_ID"
}'

# Vote from Citizen 2
curl -X POST http://localhost:3001/elections/vote \
-H "Content-Type: application/json" \
-d '{
    "citizenPublicKey": "citizen2_key",
    "candidateId": "CANDIDATE1_ID"
}'

# Vote from Citizen 3
curl -X POST http://localhost:3001/elections/vote \
-H "Content-Type: application/json" \
-d '{
    "citizenPublicKey": "citizen3_key",
    "candidateId": "CANDIDATE2_ID"
}'
```

### 6. End Election and Check Results

```bash
# End the election
curl -X POST http://localhost:3001/elections/end

# Check final results
curl http://localhost:3001/elections/current
```

## Monitoring

- Access Grafana dashboard: http://localhost:3000 (admin/admin)
- Access Prometheus metrics: http://localhost:9090

## Shutdown

To stop all nodes:
```bash
docker-compose -f docker/docker-compose.yml down
```

## Common Issues

1. Port conflicts: Make sure ports 3000-3003, 30301-30303, and 9090 are not in use
2. Docker errors: Ensure Docker Desktop is running (Windows/Mac)
3. Connection refused: Wait a few seconds after startup for all services to initialize

## Notes

- The prototype uses in-memory storage; data will be lost when containers are stopped
- All API interactions are done through node1 (port 3001) but you can use other nodes (3002, 3003) as well
- The GENESIS_ADMIN key is pre-configured for citizen approval
