set -e

# Make sure we are authed with GCR. This step must only be run once.
# gcloud auth configure-docker

git tag -a "v$KW_APP_VERSION" -m "AppVersion v$KW_APP_VERSION"
docker build -t us.gcr.io/larder-prod/kubewise:$KW_APP_VERSION .
docker push us.gcr.io/larder-prod/kubewise:$KW_APP_VERSION
