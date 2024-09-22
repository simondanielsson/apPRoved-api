#!/bin/bash

check_variables() {
    if [[ -z "$LOCATION" || -z "$PROJECT_ID" || -z "$ARTIFACT_REGISTRY" ]]; then
        echo "Error: LOCATION, PROJECT_ID, and REPOSITORY must be set."
        exit 1
    fi
}

check_variables

SOURCE_IMAGE=$1
DEST_IMAGE=$2

if [[ -z "$SOURCE_IMAGE" || -z "$DEST_IMAGE" ]]; then
    echo "Usage: ./deploy_image.sh <source-image> <destination-image>"
    exit 1
fi


GIT_SHA=$(git rev-parse --short HEAD)
if [[ $? -ne 0 ]]; then
    echo "Error: Could not get the latest git commit SHA. Are you in a git repository?"
    exit 1
fi
echo "Latest git commit SHA: $GIT_SHA"


DOCKER_TAG="$LOCATION-docker.pkg.dev/$PROJECT_ID/$ARTIFACT_REGISTRY/$DEST_IMAGE:$GIT_SHA"
echo "(1/3) Tagging image: $SOURCE_IMAGE -> $DOCKER_TAG"
docker tag "$SOURCE_IMAGE" "$DOCKER_TAG"

echo "(2/3) Pushing image: $DOCKER_TAG"
docker push "$DOCKER_TAG"
if [[ $? -ne 0 ]]; then
    echo "Error: Failed to push the image to Artifact Registry."
    exit 1
fi
echo "Image pushed successfully: $DOCKER_TAG"


echo "(3/3) Verifying the image was pushed to Artifact Registry..."
gcloud artifacts docker images list "$LOCATION-docker.pkg.dev/$PROJECT_ID/$ARTIFACT_REGISTRY/$DEST_IMAGE" --include-tags --filter="tags:$GIT_SHA"
if [[ $? -ne 0 ]]; then
    echo "Error: Image verification failed."
    exit 1
fi
echo "Image verification successful. The image with the git SHA $GIT_SHA is present."
