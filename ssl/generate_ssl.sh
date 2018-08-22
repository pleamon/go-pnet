#!/usr/bin/env bash

# clear
rm -rf ca index.txt serial private server client

mkdir -p ca/private server client

# 准备生成ca证书
mkdir private
touch serial index.txt
echo 01 > serial

# ca
CAKEY=ca/private/cakey.pem
CACERT=ca/cacert.pem
# server
SERVER_KEY=server/server_key.pem
SERVER_CSR=server/server_csr.csr
SERVER_CERT_CRT=server/server_cert.crt
SERVER_CERT_PEM=server/server_cert.pem
SERVER_CACERT=server/cacert.pem

# client
CLIENT_KEY=client/client_key.pem
CLIENT_CSR=client/client_csr.csr
CLIENT_CERT_CRT=client/client_cert.crt
CLIENT_CERT_PEM=client/client_cert.pem
CLIENT_CACERT=client/cacert.pem

# 请修改为自己的ca信息
CA_COUNTRY_NAME=CN
CA_STATE_OR_PROVINCE_NAME=Shanghai
CA_LOCATION_NAME=Shanghai
CA_ORGANIZATIONAL_NAME="Qiaoting Inc."
CA_ORGANIZATIONAL_UNIT_NAME="Qiaoting SSL Security"
CA_COMMON_NAME=qiaoting

# 修改为自己的server证书信息，SERVER_COMMON_NAME一定要是需要使用的域名地址
SERVER_COUNTRY_NAME=CN
SERVER_STATE_OR_PROVINCE_NAME=Shanghai
SERVER_LOCATION_NAME=Shanghai
SERVER_ORGANIZATIONAL_NAME="Qiaoting Inc."
SERVER_ORGANIZATIONAL_UNIT_NAME="Qiaoting Server SSL Security"
SERVER_COMMON_NAME=qiaoting

# 修改为自己的client证书信息，CLIENT_COMMON_NAME一定要是需要使用的域名地址
CLIENT_COUNTRY_NAME=CN
CLIENT_STATE_OR_PROVINCE_NAME=Shanghai
CLIENT_LOCATION_NAME=Shanghai
CLIENT_ORGANIZATIONAL_NAME="Qiaoting Inc."
CLIENT_ORGANIZATIONAL_UNIT_NAME="Qiaoting Client SSL Security"
CLIENT_COMMON_NAME=qiaoting

# 生成ca证书
openssl genrsa -out ${CAKEY} 2048
openssl req -new -x509 -key ${CAKEY} -out ${CACERT} -subj "/C=${CA_COUNTRY_NAME}/ST=${CA_STATE_OR_PROVINCE_NAME}/L=${CA_LOCATION_NAME}/O=${CA_ORGANIZATIONAL_NAME}/OU=${CA_ORGANIZATIONAL_UNIT_NAME}/CN=${CA_COMMON_NAME}"

# 生成server证书
openssl genrsa -out ${SERVER_KEY} 2048
openssl req -new -key ${SERVER_KEY} -out ${SERVER_CSR} -subj "/C=${SERVER_COUNTRY_NAME}/ST=${SERVER_STATE_OR_PROVINCE_NAME}/L=${SERVER_LOCATION_NAME}/O=${SERVER_ORGANIZATIONAL_NAME}/OU=${SERVER_ORGANIZATIONAL_UNIT_NAME}/CN=${SERVER_COMMON_NAME}"

# 使用ca为server颁发
openssl x509 -req -in ${SERVER_CSR} -CA ${CACERT} -CAkey ${CAKEY} -CAcreateserial -out ${SERVER_CERT_CRT}
# crt 转 pem
openssl x509 -inform PEM -in ${SERVER_CERT_CRT} -out ${SERVER_CERT_PEM}

# 提取证书
cp ${CACERT} ./server/

# 生成客户端证书

# 生成server证书
openssl genrsa -out ${CLIENT_KEY} 2048
openssl req -new -key ${CLIENT_KEY} -out ${CLIENT_CSR} -subj "/C=${CLIENT_COUNTRY_NAME}/ST=${CLIENT_STATE_OR_PROVINCE_NAME}/L=${CLIENT_LOCATION_NAME}/O=${CLIENT_ORGANIZATIONAL_NAME}/OU=${CLIENT_ORGANIZATIONAL_UNIT_NAME}/CN=${CLIENT_COMMON_NAME}"

# 使用ca为client颁发
openssl x509 -req -in ${CLIENT_CSR} -CA ${CACERT} -CAkey ${CAKEY} -CAcreateserial -out ${CLIENT_CERT_CRT}
# crt 转 pem
openssl x509 -inform PEM -in ${CLIENT_CERT_CRT} -out ${CLIENT_CERT_PEM}

# 提取证书
cp ${CACERT} ./client/