
export http_proxy=http://localhost:4444
http_proxy=http://localhost:4444

URLS=$(sed 's/=.*//' $HOME/i2p/hosts.txt)
for url in $URLS; do
    curl http://5fma2okrcondmxkf4j2ggwuazaoo5d3z5moyh6wurgob4nthe3oa.b32.i2p/$url
done