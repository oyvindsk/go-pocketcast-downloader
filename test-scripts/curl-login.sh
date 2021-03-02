
# Replace USERNAME and PASSWORD

curl 'https://api.pocketcasts.com/user/login' \
    -H 'Content-Type: application/json' \
    --data-raw '{"email":"USERNAME","password":"PASSWORD","scope":"webplayer"}'