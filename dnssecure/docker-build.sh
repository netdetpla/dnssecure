docker build -t dnssecure-builder . --network host --build-arg HTTP_PROXY=http://127.0.0.1:1080 --build-arg HTTPS_PROXY=http://127.0.0.1:1080
