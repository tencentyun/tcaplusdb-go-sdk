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

PKG=TcaplusGoTdrApi_${VERSION}

rm -rf ${PKG}
mkdir -p ${PKG}/src

cd ..
git submodule init
git submodule update


cp ./vendor ./pack/${PKG}/src/ -rf
mkdir -p ./pack/${PKG}/src/vendor/git.code.oa.com/gcloud_storage_group/tcaplus-go-api
rsync -av --exclude vendor ./ ./pack/${PKG}/src/vendor/git.code.oa.com/gcloud_storage_group/tcaplus-go-api
rm -rf ./pack/${PKG}/src/vendor/git.code.oa.com/gcloud_storage_group/tcaplus-go-api/.*
cd -

cp ../example ${PKG}/src/ -rf
cp ../README.md ${PKG}/src/ -rf
cp ../autotest ${PKG}/src/ -rf
VERSION_FILE=${PKG}/src/vendor/git.code.oa.com/gcloud_storage_group/tcaplus-go-api/protocol/version/version.go
GIT_COMMIT_ID_GO=`echo -e "\t"GitCommitId = \"${GIT_COMMIT_ID}\"`
VERSION_GO=`echo -e "\t"Version = \"${VERSION_SVN}\"`
sed -i "/GitCommitId/c\ ${GIT_COMMIT_ID_GO}" ${VERSION_FILE}
sed -i "/Version/c\ ${VERSION_GO}" ${VERSION_FILE}

//makefile
sed -i "/GO111MODULE=on/c export GO111MODULE=off" ${PKG}/src/example/*/Makefile
sed -i "/generic_table\/service_info/c \"example\/generic_table\/service_info\"" ${PKG}/src/example/generic_table/main.go
gofmt -w ${PKG}/src/example/generic_table/main.go
sed -i "/syncrequest\/service_info/c \"example\/syncrequest\/service_info\"" ${PKG}/src/example/syncrequest/syncrequest.go
gofmt -w ${PKG}/src/example/syncrequest/syncrequest.go

tar -zcvf ${PKG}.tar.gz ${PKG}
