export SERVER_PORT=8080
export ETCD_URL="http://10.42.42.1:4001" #http://$(ifconfig docker0 | awk '/\<inet\>/ { print $2}'):4001"
export SECURITY_KEY=goalbook123
export GITHUB_API_TOKEN=goalbook123
export DOCKER_REGISTRY_AUTH=goalbookproduct:goalbook123
