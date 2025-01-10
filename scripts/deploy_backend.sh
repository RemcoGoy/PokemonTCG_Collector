#!/bin/bash

cd backend
docker build --platform linux/amd64 -t tcg-backend:latest .
gcloud auth configure-docker europe-west1-docker.pkg.dev
docker tag tcg-backend:latest europe-west1-docker.pkg.dev/pokemontcg-446314/tcg-backend/tcg-backend:latest
docker push europe-west1-docker.pkg.dev/pokemontcg-446314/tcg-backend/tcg-backend:latest

gcloud run deploy backend --image europe-west1-docker.pkg.dev/pokemontcg-446314/tcg-backend/tcg-backend:latest --platform managed --region europe-west1 --allow-unauthenticated --port 3000 --env-vars-file .env.yaml
