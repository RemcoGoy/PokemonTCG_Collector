#!/bin/bash

cd backend && gcloud run deploy backend --source . --region europe-west1 --platform managed --allow-unauthenticated
