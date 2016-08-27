export GOPATH=$(pwd)
export PATH=$PATH:/usr/local/go/bin
rev=$(git log --pretty=format:'' | wc -l)
go build -ldflags "-s -X main.VERSION=0.1.$rev" -o ./gideon ./*.go;
echo "Build completed"
