#!/bin/bash
set -e

# brew install yq

if [ -z "$1" ] ; then
  echo "Bump mode required: major|minor|patch"
  exit 1
fi

CURRENT_VERSION=$(yq r helm_chart/Chart.yaml appVersion)
NEW_VERSION=$(./bin/bump.sh $CURRENT_VERSION $1)

echo "Bumping app version in $1 mode from $CURRENT_VERSION to $NEW_VERSION"
yq w -i helm_chart/Chart.yaml appVersion $NEW_VERSION
yq w -i helm_chart/values.yaml image.tag $NEW_VERSION

echo "Bumping chart version in $1 mode from $CURRENT_VERSION to $NEW_VERSION"
yq w -i helm_chart/Chart.yaml version $NEW_VERSION
