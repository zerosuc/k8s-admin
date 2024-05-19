#!/bin/bash

dockerComposeFilePath="deployments/docker-compose"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

mkdir -p ${dockerComposeFilePath}/configs
if [ ! -f "${dockerComposeFilePath}/configs/admin.yml" ];then
  cp configs/admin.yml ${dockerComposeFilePath}/configs
fi

# shellcheck disable=SC2164
cd ${dockerComposeFilePath}

docker-compose down
checkResult $?

docker-compose up -d
checkResult $?

colorCyan='\033[1;36m'
highBright='\033[1m'
markEnd='\033[0m'

echo ""
echo -e "run service successfully, if you want to stop the service, go into the ${highBright}${dockerComposeFilePath}${markEnd} directory and execute the command ${colorCyan}docker-compose down${markEnd}."
echo ""
