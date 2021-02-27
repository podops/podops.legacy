
1) create a new show

po new-show minimal
-> show-91804b93b76b.yaml

2) Edit show-91804b93b76b.yaml:
-> paste .yaml

3) Update the service
po update show-91804b93b76b.yaml

-> Can not find 'cover.png'

4) Upload cover.png first !
po upload cover.png

5) Try again:
po update show-91804b93b76b.yaml

-> success !

6) List all resources/assets

po get

7) Inspect resources / assets

Show the show:
po get show 91804b93b76b

Show the assets:
po get asset 6cec0a615a8a8e3bc7cfbf146a419333

8) Create a first episode

po build
-> error need min. 1 episode for a feed

parent GUID = 91804b93b76b
parent name = minimal

po template --name drums --parent minimal --parentid 91804b93b76b episode

-> show template

9) Edit episode

enclosure:
    uri: 91804b93b76b/drums.mp3
    rel: local
    type: audio/mpeg
    size: 1

enclosure:
    uri: https://cdn.podops.dev/c/default/sample.mp3
    rel: import
    type: audio/mpeg
    size: 1

10) Update the episode

po create ...

Explain create vs update and -f flags

11) Build 

po build


### build a static site using gridsome

https://levelup.gitconnected.com/how-to-build-an-awesome-website-with-gridsome-and-markdown-files-8f5422bb0183

