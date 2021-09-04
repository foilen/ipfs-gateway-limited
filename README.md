# About

To be able to share your website with IPFS while providing an HTTP fallback to your users that does not have IPFS installed locally on their machine, you need to use an IPFS Gateway. 

#1 You can use some of the public ones like:
- https://ipfs.io/ipfs/
- https://cloudflare-ipfs.com/ipfs/
- https://ipfs.github.io/public-gateway-checker/

but that means all the traffic goes through the selected one and you need to trust it. Also, if you want the Gateway to serve your website quickly, better have a dedicated one for your.

#2 You can host your own public Gateway
but:
- you most likely do not want to expose the full IPFS network ; only your specific website(s)
- you might not want to expose that this website is backed by IPFS (the paths that includes "/ipfs")

#3 You can front your local Gateway with this light "reverse proxy" that translates your host to IPFS paths and only exposes those you need.

Features:
- Reverse Proxy with a mapping from an hostname to an IPFS sub-part
- Provide X-Ipfs-Path Header for IPFS Companion (browser plugin) to use the local IPFS when available
- Resolve the mapping every minute to always use the direct object on the Gateway. Otherwise, periodically (around 10 times an hour), some requests could take more than 20 seconds to start responding due to a slow resolving step. By resolving in another thread, that leaves the web service keep serving the previously resolved objects.

## Example

Let say your files are on IPFS under:
- https://ipfs.io/ipns/k51qzi5uqu5dhmzpb6srlstaths6f2u7dpi1ru54anc1qfp9g5h4f3lnd8l6p1
- or https://ipfs.io/ipns/cdn.foilen.com

You can configure this light reverse proxy to serve
https://cdn.foilen.com/ and internaly, it will contact http://localhost:8080/ipns/cdn.foilen.com and stream the output.

Your different URLS won't contain any IPFS parts. For instance:
- requesting: https://cdn.foilen.com/foilen/Apache%20HTTP%20-%20Hotes%20virtuels.mp4
- will stream: http://localhost:8080/ipns/cdn.foilen.com/foilen/Apache%20HTTP%20-%20Hotes%20virtuels.mp4


# Local Usage

## Compile

`./create-local-release.sh`

The file is then in `build/bin/ipfs-gateway-limited`

## Config file content

```
cat > _config.json << _EOF
{
    "port" : 8888,
    "localGatewayUrl" : "http://127.0.0.1:8080",
    "localApiHostPort" : "127.0.0.1:5001",
    "mapping" : {
        "localhost.foilen-lab.com" : "/ipns/cdn.foilen.com",
        "localhost2.foilen-lab.com" : "/ipns/k51qzi5uqu5dhuj92m8egzrbx6e0apodpebvs7y4fqe1rc6rxb1hsiwgovl94o"
    }
}
_EOF
```

## Execute

To execute:
`./build/bin/ipfs-gateway-limited _config.json`

# Create release

`./create-public-release.sh`

That will show the latest created version. Then, you can choose one and execute:
`./create-public-release.sh X.X.X`

# Use with debian

Get the version you want from https://deploy.foilen.com/ipfs-gateway-limited/ .


```bash
wget https://deploy.foilen.com/ipfs-gateway-limited/ipfs-gateway-limited_X.X.X_amd64.deb
sudo dpkg -i ipfs-gateway-limited_X.X.X_amd64.deb
```
