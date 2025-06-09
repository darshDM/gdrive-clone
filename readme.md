# 📁 GDrive API (Minimal GDrive-like File Storage Service)

A simple file storage API written in Go with user authentication and file upload support, using SQLite. This project follows a `store → service → handler` architecture and is containerized using Docker. It's deployed to AWS EKS via Kubernetes.

---

## 🧱 Project Structure

```
.
├── cmd/api             # Entry point for the API
├── db/                 # SQLite DB files
├── internal/           # Internal packages (user, store, storage services)
├── uploads/            # Uploaded files
├── aws-kube.yaml       # Kubernetes manifests (Namespace, Deployment, Service, Ingress)
├── Dockerfile          # Multi-stage build for Go app
├── go.mod / go.sum     # Go module files
```

---

## 🚀 Features

- JWT-based user authentication
- SQLite for lightweight storage
- Upload and list user files
- Storage quota tracking
- RESTful API
- Docker and Kubernetes ready for deployment

---

## 🐳 Docker: Build & Run

### Build Docker Image
```bash
docker build -t gdrive-api .
```

### Run Locally
```bash
docker run -p 8000:8000 gdrive-api
```

---

## ☸️ Kubernetes on AWS (EKS)

### 1. Build and Push Docker Image to Amazon ECR
```bash
# Authenticate to ECR
aws ecr get-login-password | docker login --username AWS --password-stdin <aws_account_id>.dkr.ecr.<region>.amazonaws.com

# Tag your image
docker tag gdrive-api:latest <aws_account_id>.dkr.ecr.<region>.amazonaws.com/gdrive:latest

# Push
docker push <aws_account_id>.dkr.ecr.<region>.amazonaws.com/gdrive:latest
```

### 2. Update Image URL in `aws-kube.yaml`
```yaml
image: <aws_account_id>.dkr.ecr.<region>.amazonaws.com/gdrive:latest
```

### 3. Apply Kubernetes Resources
```bash
kubectl apply -f aws-kube.yaml
```

This will create:
- Namespace: `gdrive`
- Deployment: `gdrive-api`
- Service: NodePort (8000)
- Ingress (for ALB)

Make sure your EKS cluster has:
- ALB Ingress Controller installed
- IAM permissions set for ingress controller

---

## 🔐 API Endpoints (with `curl`)

> Replace `<token>` with your JWT token and `<file.ext>` with your file.

```bash
# 1. Signup
curl -X POST -H "Content-Type: application/json" -d '{"username": "darsh", "password": "123"}' localhost:8000/v1/users/signup

# 2. Login
curl -X POST -H "Content-Type: application/json" -d '{"username": "darsh", "password": "123"}' localhost:8000/v1/users/login

# 3. Upload File
curl -X POST -H "Authorization: Bearer <token>" -F "file=@<file.ext>" localhost:8000/v1/storage/upload

# 4. Remaining Storage
curl -X GET -H "Authorization: Bearer <token>" localhost:8000/v1/storage/remaining

# 5. List Files
curl -X GET -H "Authorization: Bearer <token>" localhost:8000/v1/storage/files

# 6. Health Check
curl -X GET localhost:8000/v1/health
```

---

## ⚙️ Environment & Notes

- Uses `sqlite3` installed in Docker image
- No external DB needed
- `uploads/` directory inside container is used to store files
- Kubernetes uses `NodePort` for internal testing, prefer `LoadBalancer` or ALB in production
- JWT signing logic is handled internally in the code