export GOPATH=$(pwd)
export PATH=$PATH:/usr/local/go/bin
cd $GOPATH/src/github.com/averrin/shodan
rev_s=$(git log --pretty=format:'' | wc -l)
cd -
rev=$(git log --pretty=format:'' | wc -l)
go build -ldflags "-s -X main.VERSION=0.1.$rev -X main.SHODAN_VERSION=$rev_s" -o ./gideon ./*.go;
echo "Build completed"
