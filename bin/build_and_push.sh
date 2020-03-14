set -e
# Make sure we are authed with GCR. This step must only be run once.
# gcloud auth configure-docker
docker build -t us.gcr.io/larder-prod/kubewise:1.0.0 .
docker push us.gcr.io/larder-prod/kubewise:1.0.0
