# PingMe API - Deployment Guide

This guide covers deploying PingMe API to various platforms.

## 📦 General Deployment Checklist

Before deploying to production:

- [ ] All unit tests pass locally (`go test -v`)
- [ ] Test coverage is above 80% (`go test -cover`)
- [ ] Update `go.mod` with your actual module path
- [ ] Set appropriate timeouts in `main.go`
- [ ] Add environment variable support for configuration
- [ ] Implement proper logging (consider structured logging)
- [ ] Add authentication/authorization if needed
- [ ] Enable CORS if serving a frontend
- [ ] Set up monitoring and alerts
- [ ] Configure HTTPS/TLS
- [ ] Review and test error handling
- [ ] CI/CD pipeline configured and passing ✅

---

## 🤖 CI/CD Pipeline

PingMe API uses GitHub Actions for automated testing and deployment.

### How It Works

Every time you push code to GitHub:

```
git push origin main
        ↓
GitHub Actions triggers automatically
        ↓
┌─────────────────────────────┐
│  1. Run unit tests          │
│  2. Check test coverage     │
│  3. Security scan           │
│  4. Lint code quality       │
│  5. Build Docker image      │
│  6. Test Docker container   │
└─────────────────────────────┘
        ↓
All pass? → Deploy automatically!
```

### Workflow Files

| File | Purpose | Triggers |
|------|---------|----------|
| `.github/workflows/ci.yml` | Run tests and build checks | Every push |
| `.github/workflows/deploy.yml` | Deploy to production | Push to main only |

### Required GitHub Secrets

Go to: **GitHub repo → Settings → Secrets and variables → Actions**

Add these secrets:

| Secret Name | Value | Where to Get It |
|------------|-------|-----------------|
| `DOCKER_USERNAME` | Your Docker Hub username | hub.docker.com |
| `DOCKER_PASSWORD` | Your Docker Hub password | hub.docker.com |
| `FLY_API_TOKEN` | Your Fly.io token | `flyctl auth token` |

### Viewing Pipeline Results

1. Go to your GitHub repository
2. Click the **Actions** tab
3. Click on any workflow run to see details
4. Green ✅ = passed, Red ❌ = failed

---

## 🐳 Docker Deployment

### Build and Push to Docker Hub

1. **Build the image:**
```bash
docker build -t yourusername/pingme-api:latest .
```

2. **Test locally:**
```bash
docker run -p 8080:8080 yourusername/pingme-api:latest
```

3. **Verify it works:**
```bash
curl http://localhost:8080/healthz
```

4. **Push to Docker Hub:**
```bash
docker login
docker push yourusername/pingme-api:latest
```

5. **Deploy anywhere Docker runs:**
```bash
docker pull yourusername/pingme-api:latest
docker run -d -p 8080:8080 --name pingme-api yourusername/pingme-api:latest
```

---

## ☁️ Cloud Platforms

### 1. Fly.io (Recommended)

**Why Fly.io?**
- Free tier available
- Simple deployment
- Automatic HTTPS
- Works great with Docker

**Steps:**

```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Login
flyctl auth login

# Launch app (run from project directory)
flyctl launch

# Deploy
flyctl deploy

# Test your live API
flyctl open
```

Your API will be at: `https://your-app.fly.dev`

**Get API token for GitHub Actions:**
```bash
flyctl auth token
# Copy this token → add to GitHub secrets as FLY_API_TOKEN
```

**Cost:** Free tier available

---

### 2. Railway

**One of the easiest deployments:**

1. **Go to [railway.app](https://railway.app)**
2. **Click "Start a New Project"**
3. **Connect GitHub repo**
4. **Railway auto-detects Dockerfile**
5. **Deploy!**

```bash
# Or use Railway CLI
npm install -g @railway/cli
railway login
railway init
railway up
```

Your API gets a URL like: `https://pingme-api.up.railway.app`

**Cost:** Free tier with 500 hours/month

---

### 3. Google Cloud Run

**Why Cloud Run?**
- Serverless, auto-scaling
- Pay per request
- Easy deployment

**Steps:**

```bash
# Install gcloud CLI
# https://cloud.google.com/sdk/docs/install

# Set your project
gcloud config set project YOUR_PROJECT_ID

# Build and deploy
gcloud run deploy pingme-api \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

Your API will be available at: `https://pingme-api-[hash]-uc.a.run.app`

**Cost:** Free tier includes 2 million requests/month

---

### 4. AWS (Elastic Container Service)

**Using AWS ECS with Fargate:**

1. **Push to AWS ECR:**
```bash
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin YOUR_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com

aws ecr create-repository --repository-name pingme-api

docker tag pingme-api:latest YOUR_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com/pingme-api:latest
docker push YOUR_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com/pingme-api:latest
```

2. **Create ECS service via AWS Console:**
   - Create Fargate task definition with your image
   - Set up load balancer
   - Deploy!

---

### 5. DigitalOcean App Platform

**Steps:**

1. Push code to GitHub
2. In DigitalOcean Dashboard: Create → Apps
3. Connect your GitHub repo
4. Choose "Dockerfile" as build method
5. Set HTTP port to 8080
6. Deploy!

**Cost:** Starting at $5/month

---

### 6. Heroku

```bash
# Create Procfile
echo "web: ./pingme-api" > Procfile

# Create and deploy
heroku create your-pingme-api
heroku buildpacks:set heroku/go
git push heroku main
```

Your API will be at: `https://your-pingme-api.herokuapp.com`

---

## 🖥️ VPS Deployment

### Using Docker on a VPS

```bash
# SSH into server
ssh user@your-server-ip

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Clone and run
git clone https://github.com/yourusername/pingme-api.git
cd pingme-api
docker build -t pingme-api .
docker run -d -p 80:8080 --name pingme-api --restart unless-stopped pingme-api
```

**Set up Nginx reverse proxy:**

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

**Set up SSL:**
```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

---

## 🔄 Kubernetes Deployment

**deployment.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pingme-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: pingme-api
  template:
    metadata:
      labels:
        app: pingme-api
    spec:
      containers:
      - name: pingme-api
        image: yourusername/pingme-api:latest
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
        readinessProbe:
          httpGet:
            path: /healthz
            port: 8080
```

```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

---

## 📊 Monitoring & Logging

Point health check monitors to: `/healthz`

For production logging consider:
- Structured JSON logging
- Log aggregation (ELK Stack, Datadog, CloudWatch)
- Error tracking (Sentry)
- Uptime monitoring (UptimeRobot - free tier available)

---

## 🔒 Security Considerations

1. **HTTPS only** - All recommended platforms provide this automatically
2. **Rate limiting** - Prevent abuse
3. **Input validation** - Already implemented in `main.go`
4. **Secrets management** - Use environment variables, never commit secrets to Git
5. **Regular updates** - Keep Go version and dependencies current

---

## 🎯 Environment Variables

```go
// In main.go - support configurable port
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
```

Set environment variables per platform:
- **Docker:** `docker run -e PORT=8080 ...`
- **Fly.io:** `flyctl secrets set PORT=8080`
- **Railway:** Dashboard → Variables
- **GitHub Actions:** Repository secrets

---

## 🆘 Troubleshooting

### CI/CD Pipeline Fails
```bash
# Check locally first
go test -v
go test -v -race

# Ensure these files are committed
git add go.mod go.sum
git commit -m "Add Go dependencies"
git push
```

### Container exits immediately
```bash
docker logs container-name
```

### Can't connect to API
```bash
# Test health endpoint
curl https://your-app.fly.dev/healthz

# Check Fly.io logs
flyctl logs
```

---

## 📚 Additional Resources

- [Go Web Applications Guide](https://golang.org/doc/articles/wiki/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Fly.io Documentation](https://fly.io/docs/)
- [12 Factor App](https://12factor.net/)

---

**Good luck with your deployment! 🚀**

If you run into issues, please create an issue on GitHub.