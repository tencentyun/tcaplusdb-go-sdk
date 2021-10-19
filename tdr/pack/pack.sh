#!/bin/bash
GIT_COMMIT_ID=`git rev-parse --short HEAD`
CUR_DATE=`date +%Y%m%d`
MAJOR=`grep MAJOR ../protocol/version/version.go | sed s/[[:space:]]//g |awk -F '=' '{print $2}' `
MINOR=`grep MINOR ../protocol/version/version.go | sed s/[[:space:]]//g |awk -F '=' '{print $2}' `
#echo MAJOR
REV=`grep REV ../protocol/version/version.go | sed s/[[:space:]]//g |awk -F '=' '{print $2}' `
VERSION=${MAJOR}.${MINOR}.${REV}.${GIT_COMMIT_ID}.x86_64_release_${CUR_DATE}
VERSION_SVN=${MAJOR}.${MINOR}.${REV}.192853.x86_64_release_${CUR_DATE}
echo -e commit_id:${GIT_COMMIT_ID}
echo -e version:${VERSION}

#pack
CURRENTPATH=$(cd "$(dirname $0)";pwd)

PKG=TcaplusGoApi_${VERSION}

rm -rf ${PKG}
mkdir -p ${PKG}/src

cd ..
git submodule init
git submodule update


cp ./vendor ./pack/${PKG}/src/ -rf
mkdir -p ./pack/${PKG}/src/vendor/github.com/tencentyun/tcaplusdb-go-sdk/tdr
rsync -av --exclude vendor ./ ./pack/${PKG}/src/vendor/github.com/tencentyun/tcaplusdb-go-sdk/tdr
rm ./pack/${PKG}/src/vendor/github.com/tencentyun/tcaplusdb-go-sdk/tdr/.* -rf
rm ./pack/${PKG}/src/vendor/github.com/tencentyun/tcaplusdb-go-sdk/tdr/go.* -rf
rm ./pack/${PKG}/src/vendor/.git -rf

cd -

cp ../example ${PKG}/src/ -rf
cp ../README.md ${PKG}/src/ -rf
cp ../autotest ${PKG}/src/ -rf
VERSION_FILE=${PKG}/src/vendor/github.com/tencentyun/tcaplusdb-go-sdk/tdr/protocol/version/version.go
GIT_COMMIT_ID_GO=`echo -e "\t"GitCommitId = \"${GIT_COMMIT_ID}\"`
VERSION_GO=`echo -e "\t"Version = \"${VERSION_SVN}\"`
sed -i "/GitCommitId/c\ ${GIT_COMMIT_ID_GO}" ${VERSION_FILE}
sed -i "/Version/c\ ${VERSION_GO}" ${VERSION_FILE}

sed -i "/GO111MODULE=on/c export GO111MODULE=off" ${PKG}/src/example/*/*/Makefile

#mv ${PKG}/src/vendor/git.woa.com ${PKG}/src/vendor/git.code.com
sed -i "s#github.com/tencentyun/tcaplusdb-go-sdk/tdr#github.com/tencentyun/tcaplusdb-go-sdk/tdr#g" `grep -rl "git.woa.com" ./${PKG}`
sed -i "s#github.com/tencentyun/tsf4g/tdrcom#github.com/tencentyun/tsf4g/tdrcom#g" `grep -rl "github.com/tencentyun/tsf4g/tdrcom" ./${PKG}`
#sed -i "s:git.code.com/gcloud_storage_group/tcaplus-go-api/autotest:autotest:g" `grep -rl "git.code.com/gcloud_storage_group/tcaplus-go-api/autotest" ./${PKG}/src/autotest`
#sed -i "s:git.code.com/gcloud_storage_group/tcaplus-go-api/example:example:g" `grep -rl "git.code.com/gcloud_storage_group/tcaplus-go-api/example" ./${PKG}/src/example`
tar -zcvf ${PKG}.tar.gz ${PKG}
