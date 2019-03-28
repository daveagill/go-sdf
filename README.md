# go-sdf - Signed Distance Fields for Go

A simple Go library for creating and manipulating Signed-Distance-Fields (and Displacement-Fields) with conversion utilities to and from images.

Included are two CLI commands (`png2sdf` and `gifanim`) that demonstrate some of the cool things you can do with Signed-Distance-Fields...

## Use `png2sdf` to generate Signed-Distance-Field textures:

The input image must be a PNG image with a **transparent** background.

    go run ./cmd/png2sdf twitter.png twitter-sdf.png
    go run ./cmd/png2sdf github.png github-sdf.png

|         twitter-sdf.png         |         github-sdf.png         |
|:-------------------------------:|:------------------------------:|
| ![](doc/images/twitter-sdf.png) | ![](doc/images/github-sdf.png) |

## Use `gifanim` to create morphing animations between images

The -from and -to images must be a PNG images with **transparent** backgrounds.

    go run ./cmd/gifanim -from=github.png -to=apple.png -out=github-to-apple.gif -frames=10
    go run ./cmd/gifanim -from=apple.png -to=twitter.png -out=apple-to-twitter.gif -frames=20
    go run ./cmd/gifanim -from=twitter.png -to=chrome.png -out=twitter-to-chrome.gif -frames=20

|          Github <-> Apple           |          Apple <-> Twitter           |          Twitter <-> Chrome           |
|:-----------------------------------:|:------------------------------------:|:-------------------------------------:|
| ![](doc/images/github-to-apple.gif) | ![](doc/images/apple-to-twitter.gif) | ![](doc/images/twitter-to-chrome.gif) |