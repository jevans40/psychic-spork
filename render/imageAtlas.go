package render

import "image"

type ImageAtlas struct {
	imageSize int32
	imagesNum int32
	images    []image.RGBA
	imagepath string
	Allocator imageAllocator
}

type imageAllocator struct {
	size         int32
	allocated    [4]bool
	childPointer [4]*imageAllocator
}

func ImageAtlasFactory(ImageSize int32, Images int32, ImageFolder string) ImageAtlas {

	//Allocator := new(imageAllocator)
	return *new(ImageAtlas)
}

func (i *imageAllocator) Allocate(ToAllocate int32) {

}
