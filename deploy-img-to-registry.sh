IMAGE_NAME="gym"
VERSION="v1"

echo $'\n##############################'
echo " DEPLOYING IMAGE TO GITLAB REGISTRY "
echo $'##############################\n'
docker login registry.gitlab.com
docker build -f "Dockerfile" -t $IMAGE_NAME:$VERSION .
docker tag $IMAGE_NAME:$VERSION registry.gitlab.com/gym-global/backend/$IMAGE_NAME:$VERSION
docker push registry.gitlab.com/gym-global/backend/$IMAGE_NAME:$VERSION