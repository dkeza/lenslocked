#!/bin/bash

# Change to the directory with our code that we plan to work from
cd "$GOPATH/src/simplegallery"

echo "==== Releasing simplegallery ===="
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm simplegallery
echo "  Done!"

echo "  Deleting existing code..."
ssh root@simplegallery.kezic.net "rm -rf /root/go/src/simplegallery"
echo "  Code deleted successfully!"

echo "  Uploading code..."
rsync -avr --exclude '.git/*' --exclude 'tmp/*' --exclude 'images/*' ./ root@simplegallery.kezic.net:/root/go/src/simplegallery/
echo "  Code uploaded successfully!"

echo "  Go getting deps..."
ssh root@simplegallery.kezic.net "export GOPATH=/root/go; /usr/local/go/bin/go get golang.org/x/crypto/bcrypt"
ssh root@simplegallery.kezic.net "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/mux"
ssh root@simplegallery.kezic.net "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/schema"
ssh root@simplegallery.kezic.net "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/lib/pq"
ssh root@simplegallery.kezic.net "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/jinzhu/gorm"
ssh root@simplegallery.kezic.net "export GOPATH=/root/go; /usr/local/go/bin/go get github.com/gorilla/csrf"

echo "  Building the code on remote server..."
ssh root@simplegallery.kezic.net 'export GOPATH=/root/go; cd /root/app; /usr/local/go/bin/go build -o ./server $GOPATH/src/simplegallery/*.go'
echo "  Code built successfully!"

echo "  Moving assets..."
ssh root@simplegallery.kezic.net "cd /root/app; cp -R /root/go/src/simplegallery/assets ."
echo "  Assets moved successfully!"

echo "  Moving views..."
ssh root@simplegallery.kezic.net "cd /root/app; cp -R /root/go/src/simplegallery/views ."
echo "  Views moved successfully!"

echo "  Moving Caddyfile..."
ssh root@simplegallery.kezic.net "cd /root/app; cp /root/go/src/simplegallery/Caddyfile ."
echo "  Views moved successfully!"

echo "  Restarting the server..."
ssh root@simplegallery.kezic.net "sudo service simplegallery restart"
echo "  Server restarted successfully!"

echo "  Restarting Caddy server..."
ssh root@simplegallery.kezic.net "sudo service caddy restart"
echo "  Caddy restarted successfully!"

echo "==== Done releasing simplegallery ===="
