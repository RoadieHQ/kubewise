#!/bin/bash
set -e

# brew install yq
# brew install goreleaser
# brew install helm

PATH_TO_HELM_REPO=~/dev/roadie/helm-repo

usage() {
  echo "Usage: ./bin/release.sh {major|minor|patch}"
  echo "Bumps the semantic version field by one for a project, tags git, releases the binaries and docker containers and builds a helm package."
  exit 1
}

if [ -z "$1" ] ; then
  usage
fi

BUMP_MODE="$1"

./bin/bump_chart.sh $BUMP_MODE
./bin/bump_app.sh $BUMP_MODE
BUMPED_CHART_VERSION=$(yq r helm_chart/Chart.yaml version)
BUMPED_APP_VERSION=$(yq r helm_chart/Chart.yaml appVersion)

git add --all
git commit -m "Bump: ChartVersion $BUMPED_CHART_VERSION. AppVersion $BUMPED_APP_VERSION"

# The tag cannot be the app version because there will be cases 
# where we want to make a new release even though the app has not
# changed. Perhaps just the chart values have.
git tag -a "$BUMPED_CHART_VERSION" -m "ChartVersion $BUMPED_CHART_VERSION"
git push
git push origin $BUMPED_CHART_VERSION

goreleaser --rm-dist

helm package ./helm_chart
mv kubewise-*.tgz $PATH_TO_HELM_REPO
