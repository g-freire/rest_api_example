IMAGE_NAME="gym"
VERSION="v1"

echo $'\n##############################'
echo " DEPLOYING IMAGE TO GITLAB REGISTRY "
echo $'##############################\n'
docker login registry.gitlab.com
docker build -f "Dockerfile-EKS" -t $IMAGE_NAME:$VERSION .
docker tag $IMAGE_NAME:$VERSION registry.gitlab.com/prettytechnical/domino/extractor/$IMAGE_NAME:$VERSION
docker push registry.gitlab.com/prettytechnical/domino/extractor/$IMAGE_NAME:$VERSION
