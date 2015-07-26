#!/bin/sh
git add -A
git commit -m "A"
git push origin master:master
read tags
git tag -a $tags -m $tags
git push origin --tag $tags
sleep 30
git push origin --tag :$tags
git tag -d $tags