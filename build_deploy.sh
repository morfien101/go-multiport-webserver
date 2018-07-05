if [ -f ./dual_port ]; then
  echo "Delete dual_port"
  rm -f ./dual_port
fi

echo "Building"
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dual_port . \
&& chmod 550 dual_port \
&& docker build -t morfien101/multi_port:$(./dual_port -v) .

if [ $# -ne 0 ]; then
  if [ $1 = "push" ]; then
    echo "Push to docker"
    VERSION_TAG=$(./dual_port -v)
    docker tag morfien101/multi_port:$VERSION_TAG morfien101/multi_port:latest
    docker push morfien101/multi_port:$VERSION_TAG
    docker push morfien101/multi_port:latest
  else
    echo "Not pushing to Dockerhub"
  fi
else
  echo "Not pushing to Dockerhub"
fi

