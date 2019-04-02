#  version

record version information to Go program. include:

* git tag
* build date
* git commitment
* git tree state

## Usage

edit follow go file:

```go
package main

import (
	"github.com/hzxiao/goutil/version"
)

var showVersion = true

func main() {
	if showVersion {
		err := version.Print()
		if err != nil {
			panic(err)
		}
		return
	}
}
```

provider follow ldflags when build

```shell
versionDir=github.com/hzxiao/goutil/version
PWD=`pwd`
if [[ ${PWD} = ${GOPATH}/* ]]; then
    if [[ -d vendor ]]; then
        gp=${GOPATH//\//\\\/}
        proj_path=`echo ${PWD} | sed "s/${gp}\/src\///g"`
        versionDir=${proj_path}/vendor/github.com/hzxiao/goutil/version
    fi
fi
echo "version dir: ${versionDir}"

gitTag=$(if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:''%h'' -n 1; fi)
buildDate=$(TZ=Asia/Shanghai date +%FT%T%z)
gitCommit=$(git log --pretty=format:''%H'' -n 1)
gitTreeState=$(if git status|grep -q ''clean'';then echo clean; else echo dirty; fi)

go build -v -ldflags "-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}" .
```