#!/bin/bash

gcloud run deploy --source . --region europe-west1 --platform managed --allow-unauthenticated
