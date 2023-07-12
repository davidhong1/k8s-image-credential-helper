export INIT_CONFIG=environment
export IMAGE_PROVIDER=harbor
export IMAGE_HOST=myharbor.io
export IMAGE_USER=xxx
export IMAGE_PASSWORD=xxx

# build
make build

# run
./$1
