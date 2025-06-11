# GDrive API (Minimal GDrive-like File Storage Service)

A simple file storage API written in Go with user authentication and file upload support, using SQLite. This project follows a `store → service → handler` architecture and is containerized using Docker. It's deployed to AWS EKS via Kubernetes.

---

## Project Structure

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

## Basic Features

- JWT-based user authentication
- SQLite for lightweight storage
- Upload and list user files
---

## Docker: Build & Run

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

### 1. Update Kubeconfig
```bash
aws eks update-kubeconfig --name <cluster-name> --region <region>
```

### 2. Create Fargate Profile (Optional)
```bash
eksctl create fargateprofile \
--cluster <cluster-name> \
--region <region> \
--name alb-sample-app \
--namespace game-2048
```

### 3. Build & Push Docker Image to Amazon ECR
```bash
# Authenticate to ECR
aws ecr get-login-password | docker login --username AWS --password-stdin <aws_account_id>.dkr.ecr.<region>.amazonaws.com

# Tag your image
docker tag gdrive-api:latest <aws_account_id>.dkr.ecr.<region>.amazonaws.com/gdrive:latest

# Push
docker push <aws_account_id>.dkr.ecr.<region>.amazonaws.com/gdrive:latest
```

### 4. Update Image URL in `aws-kube.yaml`
```yaml
image: <aws_account_id>.dkr.ecr.<region>.amazonaws.com/gdrive:latest
```

### 5. Deploy to Cluster
```bash
kubectl apply -f aws-kube.yaml
```

---

## Setup ALB Ingress Controller

### 1. Download IAM Policy
```bash
curl -O https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/v2.11.0/docs/install/iam_policy.json
```

### 2. Create IAM Policy
```bash
aws iam create-policy \
--policy-name AWSLoadBalancerControllerIAMPolicy \
--policy-document file://iam_policy.json
```

### 3. Create IAM Role and Service Account
```bash
eksctl create iamserviceaccount \
--cluster=<your-cluster-name> \
--namespace=kube-system \
--name=aws-load-balancer-controller \
--role-name AmazonEKSLoadBalancerControllerRole \
--attach-policy-arn=arn:aws:iam::<your-aws-account-id>:policy/AWSLoadBalancerControllerIAMPolicy \
--approve
```

### 4. Add ALB Helm Repo
```bash
helm repo add eks https://aws.github.io/eks-charts
```

### 5. Install ALB Controller via Helm
```bash
helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
  -n kube-system \
  --set clusterName=<your-cluster-name> \
  --set serviceAccount.create=false \
  --set serviceAccount.name=aws-load-balancer-controller \
  --set region=<region> \
  --set vpcId=<your-vpc-id>
```

---

## API Endpoints (with `curl`)

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

## Environment & Notes

- Uses `sqlite3` installed in Docker image
- No external DB needed
- `uploads/` directory inside container is used to store files
- Kubernetes uses `NodePort` for internal testing, prefer `LoadBalancer` or ALB in production
- JWT signing logic is handled internally in the code