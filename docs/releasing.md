# How to release new versions

This shows how to do it locally. This should all be wrapped up in a GitHub action so it can
happen automatically.

 1. Get a GitHub personal access token from GitHub. It needs the ability to write to
    the KubeWise GitHub repo.
 2. Log in to Docker so that the `docker` CLI has the permissions to publish images to
    Docker Hub. This may require getting a Docker Hub access token and manually logging
    in with `docker login --username`.
 3. Set the GitHub token in the environment with `export GITHUB_TOKEN=<my-token>`.
 4. Release with: `./bin/release.sh  {major|minor|patch}`
 5. Package the Helm chart: `helm package ./helm_chart`.
 6. Move it to the helm repos: `mv kubewise-*.tgz $PATH_TO_HELM_REPO`
 7. Add and commit that repo to publish the chart automatically.
