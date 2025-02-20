package rl

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

import (
	"image"
	"image/color"
	"unsafe"
)

// cptr returns C pointer
func (i *Image) cptr() *C.Image {
	return (*C.Image)(unsafe.Pointer(i))
}

// ToImage converts a Image to Go image.Image
func (i *Image) ToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(i.Width), int(i.Height)))

	// Get pixel data from image (RGBA 32bit)
	cimg := i.cptr()
	ret := C.LoadImageColors(*cimg)
	pixels := (*[1 << 24]uint8)(unsafe.Pointer(ret))[0 : i.Width*i.Height*4]

	img.Pix = pixels

	return img
}

// cptr returns C pointer
func (t *Texture2D) cptr() *C.Texture2D {
	return (*C.Texture2D)(unsafe.Pointer(t))
}

// cptr returns C pointer
func (r *RenderTexture2D) cptr() *C.RenderTexture2D {
	return (*C.RenderTexture2D)(unsafe.Pointer(r))
}

// LoadImage - Load an image into CPU memory (RAM)
func LoadImage(fileName string) *Image {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadImage(cfileName)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadImageRaw - Load image data from RAW file
func LoadImageRaw(fileName string, width, height int32, format PixelFormat, headerSize int32) *Image {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cformat := (C.int)(format)
	cheaderSize := (C.int)(headerSize)
	ret := C.LoadImageRaw(cfileName, cwidth, cheight, cformat, cheaderSize)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadImageAnim - Load image sequence from file (frames appended to image.data)
func LoadImageAnim(fileName string, frames *int32) *Image {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cframes := (*C.int)(frames)
	ret := C.LoadImageAnim(cfileName, cframes)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadImageFromMemory - Load image from memory buffer, fileType refers to extension: i.e. ".png"
func LoadImageFromMemory(fileType string, fileData []byte, dataSize int32) *Image {
	cfileType := C.CString(fileType)
	defer C.free(unsafe.Pointer(cfileType))
	cfileData := (*C.uchar)(unsafe.Pointer(&fileData[0]))
	cdataSize := (C.int)(dataSize)
	ret := C.LoadImageFromMemory(cfileType, cfileData, cdataSize)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadTexture - Load an image as texture into GPU memory
func LoadTexture(fileName string) Texture2D {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadTexture(cfileName)
	v := newTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadTextureFromImage - Load a texture from image data
func LoadTextureFromImage(image *Image) Texture2D {
	cimage := image.cptr()
	ret := C.LoadTextureFromImage(*cimage)
	v := newTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadRenderTexture - Load a texture to be used for rendering
func LoadRenderTexture(width, height int32) RenderTexture2D {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ret := C.LoadRenderTexture(cwidth, cheight)
	v := newRenderTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// UnloadImage - Unload image from CPU memory (RAM)
func UnloadImage(image *Image) {
	cimage := image.cptr()
	C.UnloadImage(*cimage)
}

// UnloadTexture - Unload texture from GPU memory
func UnloadTexture(texture Texture2D) {
	ctexture := texture.cptr()
	C.UnloadTexture(*ctexture)
}

// UnloadRenderTexture - Unload render texture from GPU memory
func UnloadRenderTexture(target RenderTexture2D) {
	ctarget := target.cptr()
	C.UnloadRenderTexture(*ctarget)
}

// LoadImageColors - Get pixel data from image as a Color slice
func LoadImageColors(img *Image) []color.RGBA {
	cimg := img.cptr()
	ret := C.LoadImageColors(*cimg)
	return (*[1 << 24]color.RGBA)(unsafe.Pointer(ret))[0 : img.Width*img.Height]
}

// LoadImageFromTexture - Get pixel data from GPU texture and return an Image
func LoadImageFromTexture(texture Texture2D) *Image {
	ctexture := texture.cptr()
	ret := C.LoadImageFromTexture(*ctexture)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// UpdateTexture - Update GPU texture with new data
func UpdateTexture(texture Texture2D, pixels []color.RGBA) {
	ctexture := texture.cptr()
	cpixels := unsafe.Pointer(&pixels[0])
	C.UpdateTexture(*ctexture, cpixels)
}

// UpdateTextureRec - Update GPU texture rectangle with new data
func UpdateTextureRec(texture Texture2D, rec Rectangle, pixels []color.RGBA) {
	ctexture := texture.cptr()
	cpixels := unsafe.Pointer(&pixels[0])
	crec := rec.cptr()
	C.UpdateTextureRec(*ctexture, *crec, cpixels)
}

// ExportImage - Export image as a PNG file
func ExportImage(image Image, name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cimage := image.cptr()

	C.ExportImage(*cimage, cname)
}

// ImageCopy - Create an image duplicate (useful for transformations)
func ImageCopy(image *Image) *Image {
	cimage := image.cptr()
	ret := C.ImageCopy(*cimage)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// ImageText - Create an image from text (default font)
func ImageText(text string, fontSize int32, col color.RGBA) *Image {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.int)(fontSize)
	ccolor := colorCptr(col)
	ret := C.ImageText(ctext, cfontSize, *ccolor)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// ImageTextEx - Create an image from text (custom sprite font)
func ImageTextEx(font Font, text string, fontSize, spacing float32, tint color.RGBA) *Image {
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.float)(fontSize)
	cspacing := (C.float)(spacing)
	ctint := colorCptr(tint)
	ret := C.ImageTextEx(*cfont, ctext, cfontSize, cspacing, *ctint)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// ImageFormat - Convert image data to desired format
func ImageFormat(image *Image, newFormat PixelFormat) {
	cimage := image.cptr()
	cnewFormat := (C.int)(newFormat)
	C.ImageFormat(cimage, cnewFormat)
}

// ImageToPOT - Convert image to POT (power-of-two)
func ImageToPOT(image *Image, fillColor color.RGBA) {
	cimage := image.cptr()
	cfillColor := colorCptr(fillColor)
	C.ImageToPOT(cimage, *cfillColor)
}

// ImageCrop - Crop an image to a defined rectangle
func ImageCrop(image *Image, crop Rectangle) {
	cimage := image.cptr()
	ccrop := crop.cptr()
	C.ImageCrop(cimage, *ccrop)
}

// ImageAlphaCrop - Crop image depending on alpha value
func ImageAlphaCrop(image *Image, threshold float32) {
	cimage := image.cptr()
	cthreshold := (C.float)(threshold)
	C.ImageAlphaCrop(cimage, cthreshold)
}

// ImageAlphaClear - Apply alpha mask to image
func ImageAlphaClear(image *Image, col color.RGBA, threshold float32) {
	cimage := image.cptr()
	ccolor := colorCptr(col)
	cthreshold := (C.float)(threshold)
	C.ImageAlphaClear(cimage, *ccolor, cthreshold)
}

// ImageAlphaMask - Apply alpha mask to image
func ImageAlphaMask(image, alphaMask *Image) {
	cimage := image.cptr()
	calphaMask := alphaMask.cptr()
	C.ImageAlphaMask(cimage, *calphaMask)
}

// ImageAlphaPremultiply - Premultiply alpha channel
func ImageAlphaPremultiply(image *Image) {
	cimage := image.cptr()
	C.ImageAlphaPremultiply(cimage)
}

// ImageResize - Resize an image (bilinear filtering)
func ImageResize(image *Image, newWidth, newHeight int32) {
	cimage := image.cptr()
	cnewWidth := (C.int)(newWidth)
	cnewHeight := (C.int)(newHeight)
	C.ImageResize(cimage, cnewWidth, cnewHeight)
}

// ImageResizeNN - Resize an image (Nearest-Neighbor scaling algorithm)
func ImageResizeNN(image *Image, newWidth, newHeight int32) {
	cimage := image.cptr()
	cnewWidth := (C.int)(newWidth)
	cnewHeight := (C.int)(newHeight)
	C.ImageResizeNN(cimage, cnewWidth, cnewHeight)
}

// ImageResizeCanvas - Resize canvas and fill with color
func ImageResizeCanvas(image *Image, newWidth, newHeight, offsetX, offsetY int32, col color.RGBA) {
	cimage := image.cptr()
	cnewWidth := (C.int)(newWidth)
	cnewHeight := (C.int)(newHeight)
	coffsetX := (C.int)(offsetX)
	coffsetY := (C.int)(offsetY)
	ccolor := colorCptr(col)
	C.ImageResizeCanvas(cimage, cnewWidth, cnewHeight, coffsetX, coffsetY, *ccolor)
}

// ImageMipmaps - Generate all mipmap levels for a provided image
func ImageMipmaps(image *Image) {
	cimage := image.cptr()
	C.ImageMipmaps(cimage)
}

// ImageDither - Dither image data to 16bpp or lower (Floyd-Steinberg dithering)
func ImageDither(image *Image, rBpp, gBpp, bBpp, aBpp int32) {
	cimage := image.cptr()
	crBpp := (C.int)(rBpp)
	cgBpp := (C.int)(gBpp)
	cbBpp := (C.int)(bBpp)
	caBpp := (C.int)(aBpp)
	C.ImageDither(cimage, crBpp, cgBpp, cbBpp, caBpp)
}

// ImageFlipVertical - Flip image vertically
func ImageFlipVertical(image *Image) {
	cimage := image.cptr()
	C.ImageFlipVertical(cimage)
}

// ImageFlipHorizontal - Flip image horizontally
func ImageFlipHorizontal(image *Image) {
	cimage := image.cptr()
	C.ImageFlipHorizontal(cimage)
}

// ImageRotateCW - Rotate image clockwise 90deg
func ImageRotateCW(image *Image) {
	cimage := image.cptr()
	C.ImageRotateCW(cimage)
}

// ImageRotateCCW - Rotate image counter-clockwise 90deg
func ImageRotateCCW(image *Image) {
	cimage := image.cptr()
	C.ImageRotateCCW(cimage)
}

// ImageColorTint - Modify image color: tint
func ImageColorTint(image *Image, col color.RGBA) {
	cimage := image.cptr()
	ccolor := colorCptr(col)
	C.ImageColorTint(cimage, *ccolor)
}

// ImageColorInvert - Modify image color: invert
func ImageColorInvert(image *Image) {
	cimage := image.cptr()
	C.ImageColorInvert(cimage)
}

// ImageColorGrayscale - Modify image color: grayscale
func ImageColorGrayscale(image *Image) {
	cimage := image.cptr()
	C.ImageColorGrayscale(cimage)
}

// ImageColorContrast - Modify image color: contrast (-100 to 100)
func ImageColorContrast(image *Image, contrast float32) {
	cimage := image.cptr()
	ccontrast := (C.float)(contrast)
	C.ImageColorContrast(cimage, ccontrast)
}

// ImageColorBrightness - Modify image color: brightness (-255 to 255)
func ImageColorBrightness(image *Image, brightness int32) {
	cimage := image.cptr()
	cbrightness := (C.int)(brightness)
	C.ImageColorBrightness(cimage, cbrightness)
}

// ImageColorReplace - Modify image color: replace color
func ImageColorReplace(image *Image, col, replace color.RGBA) {
	cimage := image.cptr()
	ccolor := colorCptr(col)
	creplace := colorCptr(replace)
	C.ImageColorReplace(cimage, *ccolor, *creplace)
}

// GetImageColor - Get image pixel color at (x, y) position
func GetImageColor(image Image, x, y int32) color.RGBA {
	cimage := image.cptr()
	cx := (C.int)(x)
	cy := (C.int)(y)

	ret := C.GetImageColor(*cimage, cx, cy)
	v := newColorFromPointer(unsafe.Pointer(&ret))
	return v
}

// ImageClearBackground - Clear image background with given color
func ImageClearBackground(dst *Image, col color.RGBA) {
	cdst := dst.cptr()
	ccolor := colorCptr(col)
	C.ImageClearBackground(cdst, *ccolor)
}

// ImageDraw - Draw a source image within a destination image
func ImageDraw(dst, src *Image, srcRec, dstRec Rectangle, tint color.RGBA) {
	cdst := dst.cptr()
	csrc := src.cptr()
	csrcRec := srcRec.cptr()
	cdstRec := dstRec.cptr()
	ctint := colorCptr(tint)
	C.ImageDraw(cdst, *csrc, *csrcRec, *cdstRec, *ctint)
}

// ImageDrawCircle - Draw circle within an image
func ImageDrawCircle(dst *Image, centerX, centerY, radius int32, col color.RGBA) {
	cdst := dst.cptr()
	ccenterX := (C.int)(centerX)
	ccenterY := (C.int)(centerY)
	cradius := (C.int)(radius)
	ccolor := colorCptr(col)
	C.ImageDrawCircle(cdst, ccenterX, ccenterY, cradius, *ccolor)
}

// ImageDrawCircleV - Draw circle within an image
func ImageDrawCircleV(dst *Image, center Vector2, radius int32, col color.RGBA) {
	cdst := dst.cptr()
	ccenter := center.cptr()
	cradius := (C.int)(radius)
	ccolor := colorCptr(col)
	C.ImageDrawCircleV(cdst, *ccenter, cradius, *ccolor)
}

// ImageDrawPixel - Draw pixel within an image
func ImageDrawPixel(dst *Image, posX, posY int32, col color.RGBA) {
	cdst := dst.cptr()
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	ccolor := colorCptr(col)
	C.ImageDrawPixel(cdst, cposX, cposY, *ccolor)
}

// ImageDrawPixelV - Draw pixel within an image (Vector version)
func ImageDrawPixelV(dst *Image, position Vector2, col color.RGBA) {
	cdst := dst.cptr()
	cposition := position.cptr()
	ccolor := colorCptr(col)
	C.ImageDrawPixelV(cdst, *cposition, *ccolor)
}

// ImageDrawRectangle - Draw rectangle within an image
func ImageDrawRectangle(dst *Image, x, y, width, height int32, col color.RGBA) {
	cdst := dst.cptr()
	cx := (C.int)(x)
	cy := (C.int)(y)
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ccolor := colorCptr(col)
	C.ImageDrawRectangle(cdst, cx, cy, cwidth, cheight, *ccolor)
}

// ImageDrawRectangleV - Draw rectangle within an image (Vector version)
func ImageDrawRectangleV(dst *Image, position, size Vector2, col color.RGBA) {
	cdst := dst.cptr()
	cposition := position.cptr()
	csize := size.cptr()
	ccolor := colorCptr(col)
	C.ImageDrawRectangleV(cdst, *cposition, *csize, *ccolor)
}

// ImageDrawRectangleLines - Draw rectangle lines within an image
func ImageDrawRectangleLines(dst *Image, rec Rectangle, thick int, col color.RGBA) {
	cdst := dst.cptr()
	crec := rec.cptr()
	cthick := (C.int)(thick)
	ccolor := colorCptr(col)
	C.ImageDrawRectangleLines(cdst, *crec, cthick, *ccolor)
}

// ImageDrawRectangleRec - Draw rectangle within an image
func ImageDrawRectangleRec(dst *Image, rec Rectangle, col color.RGBA) {
	cdst := dst.cptr()
	crec := rec.cptr()
	ccolor := colorCptr(col)
	C.ImageDrawRectangleRec(cdst, *crec, *ccolor)
}

// ImageDrawText - Draw text (default font) within an image (destination)
func ImageDrawText(dst *Image, posX, posY int32, text string, fontSize int32, col color.RGBA) {
	cdst := dst.cptr()
	posx := (C.int)(posX)
	posy := (C.int)(posY)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.int)(fontSize)
	ccolor := colorCptr(col)
	C.ImageDrawText(cdst, ctext, posx, posy, cfontSize, *ccolor)
}

// ImageDrawTextEx - Draw text (custom sprite font) within an image (destination)
func ImageDrawTextEx(dst *Image, position Vector2, font Font, text string, fontSize, spacing float32, col color.RGBA) {
	cdst := dst.cptr()
	cposition := position.cptr()
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.float)(fontSize)
	cspacing := (C.float)(spacing)
	ccolor := colorCptr(col)
	C.ImageDrawTextEx(cdst, *cfont, ctext, *cposition, cfontSize, cspacing, *ccolor)
}

// GenImageColor - Generate image: plain color
func GenImageColor(width, height int, col color.RGBA) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ccolor := colorCptr(col)

	ret := C.GenImageColor(cwidth, cheight, *ccolor)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageGradientV - Generate image: vertical gradient
func GenImageGradientV(width, height int, top, bottom color.RGBA) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ctop := colorCptr(top)
	cbottom := colorCptr(bottom)

	ret := C.GenImageGradientV(cwidth, cheight, *ctop, *cbottom)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageGradientH - Generate image: horizontal gradient
func GenImageGradientH(width, height int, left, right color.RGBA) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cleft := colorCptr(left)
	cright := colorCptr(right)

	ret := C.GenImageGradientH(cwidth, cheight, *cleft, *cright)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageGradientRadial - Generate image: radial gradient
func GenImageGradientRadial(width, height int, density float32, inner, outer color.RGBA) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cdensity := (C.float)(density)
	cinner := colorCptr(inner)
	couter := colorCptr(outer)

	ret := C.GenImageGradientRadial(cwidth, cheight, cdensity, *cinner, *couter)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageChecked - Generate image: checked
func GenImageChecked(width, height, checksX, checksY int, col1, col2 color.RGBA) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cchecksX := (C.int)(checksX)
	cchecksY := (C.int)(checksY)
	ccol1 := colorCptr(col1)
	ccol2 := colorCptr(col2)

	ret := C.GenImageChecked(cwidth, cheight, cchecksX, cchecksY, *ccol1, *ccol2)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageWhiteNoise - Generate image: white noise
func GenImageWhiteNoise(width, height int, factor float32) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cfactor := (C.float)(factor)

	ret := C.GenImageWhiteNoise(cwidth, cheight, cfactor)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageCellular - Generate image: cellular algorithm. Bigger tileSize means bigger cells
func GenImageCellular(width, height, tileSize int) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ctileSize := (C.int)(tileSize)

	ret := C.GenImageCellular(cwidth, cheight, ctileSize)
	v := newImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenTextureMipmaps - Generate GPU mipmaps for a texture
func GenTextureMipmaps(texture *Texture2D) {
	ctexture := texture.cptr()
	C.GenTextureMipmaps(ctexture)
}

// SetTextureFilter - Set texture scaling filter mode
func SetTextureFilter(texture Texture2D, filterMode TextureFilterMode) {
	ctexture := texture.cptr()
	cfilterMode := (C.int)(filterMode)
	C.SetTextureFilter(*ctexture, cfilterMode)
}

// SetTextureWrap - Set texture wrapping mode
func SetTextureWrap(texture Texture2D, wrapMode TextureWrapMode) {
	ctexture := texture.cptr()
	cwrapMode := (C.int)(wrapMode)
	C.SetTextureWrap(*ctexture, cwrapMode)
}

// DrawTexture - Draw a Texture2D
func DrawTexture(texture Texture2D, posX int32, posY int32, tint color.RGBA) {
	ctexture := texture.cptr()
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	ctint := colorCptr(tint)
	C.DrawTexture(*ctexture, cposX, cposY, *ctint)
}

// DrawTextureV - Draw a Texture2D with position defined as Vector2
func DrawTextureV(texture Texture2D, position Vector2, tint color.RGBA) {
	ctexture := texture.cptr()
	cposition := position.cptr()
	ctint := colorCptr(tint)
	C.DrawTextureV(*ctexture, *cposition, *ctint)
}

// DrawTextureEx - Draw a Texture2D with extended parameters
func DrawTextureEx(texture Texture2D, position Vector2, rotation, scale float32, tint color.RGBA) {
	ctexture := texture.cptr()
	cposition := position.cptr()
	crotation := (C.float)(rotation)
	cscale := (C.float)(scale)
	ctint := colorCptr(tint)
	C.DrawTextureEx(*ctexture, *cposition, crotation, cscale, *ctint)
}

// DrawTextureRec - Draw a part of a texture defined by a rectangle
func DrawTextureRec(texture Texture2D, sourceRec Rectangle, position Vector2, tint color.RGBA) {
	ctexture := texture.cptr()
	csourceRec := sourceRec.cptr()
	cposition := position.cptr()
	ctint := colorCptr(tint)
	C.DrawTextureRec(*ctexture, *csourceRec, *cposition, *ctint)
}

// DrawTexturePro - Draw a part of a texture defined by a rectangle with 'pro' parameters
func DrawTexturePro(texture Texture2D, sourceRec, destRec Rectangle, origin Vector2, rotation float32, tint color.RGBA) {
	ctexture := texture.cptr()
	csourceRec := sourceRec.cptr()
	cdestRec := destRec.cptr()
	corigin := origin.cptr()
	crotation := (C.float)(rotation)
	ctint := colorCptr(tint)
	C.DrawTexturePro(*ctexture, *csourceRec, *cdestRec, *corigin, crotation, *ctint)
}
