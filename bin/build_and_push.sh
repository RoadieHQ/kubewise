set -e

# Make sure we are authed with GCR. This step must only be run once.
# gcloud auth configure-docker

docker_state=$(docker info >/dev/null 2>&1)
if [[ $? -ne 0 ]]; then
    echo "Docker does not seem to be running, run it first and retry"
    exit 1
fi

git tag -a "v$KW_APP_VERSION" -m "AppVersion v$KW_APP_VERSION"
docker build -t us.gcr.io/larder-prod/kubewise:$KW_APP_VERSION .
docker push us.gcr.io/larder-prod/kubewise:$KW_APP_VERSION
