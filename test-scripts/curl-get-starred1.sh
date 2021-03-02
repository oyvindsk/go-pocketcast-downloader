
# TODO : 
# - replace token
# - cleanup headers

# This is a POST, becasue of --data-raw

curl 'https://api.pocketcasts.com/user/starred' \
    -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:86.0) Gecko/20100101 Firefox/86.0' \
    -H 'Accept: */*' \
    -H 'Accept-Language: en-US,en;q=0.5' \
    --compressed \
    -H 'Referer: https://play.pocketcasts.com/starred' \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer JWT-TOKEN-FIXME..' \
    -H 'Origin: https://play.pocketcasts.com' \
    -H 'DNT: 1' \
    -H 'Connection: keep-alive' \
    -H 'Sec-GPC: 1' \
    -H 'Pragma: no-cache' \
    -H 'Cache-Control: no-cache' \
    -H 'TE: Trailers' \
    --data-raw '{}'
