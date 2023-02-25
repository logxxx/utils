package screen

import (
	"bytes"
	"github.com/kbinani/screenshot"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
)

func SaveToLocal(img *image.RGBA, filename string) error {
	data, err := encode(img)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		return err
	}

	return nil
}

func ShotRectOrig(x1, y1, x2, y2 int) (*image.RGBA, error) {
	w := x2 - x1
	h := y2 - y1
	//log.Printf("ShotRect x:%v y:%v w:%v h:%v", x1, y1, w, h)
	shotImg, err := screenshot.Capture(x1, y1, w, h)
	if err != nil {
		log.Printf("ShotRect CaptureDisplay err:%v x1:%v y1:%v w:%v h:%v",
			err, x1, y1, w, h)
		return nil, err
	}
	return shotImg, nil
}

func ShotRect(x1, y1, x2, y2 int) ([]byte, error) {

	shotImg, err := ShotRectOrig(x1, y1, x2, y2)
	if err != nil {
		return nil, err
	}

	return encode(shotImg)
}

func Shot() ([]byte, error) {

	shotImg, err := screenshot.CaptureDisplay(0)
	if err != nil {
		log.Printf("Shot CaptureDisplay err:%v", err)
		return nil, err
	}

	return encode(shotImg)
}

func ShotAreaAndResize(x1, y1, x2, y2, w, h int) ([]byte, error) {
	shotImg, err := screenshot.Capture(x1, y1, x2-x1, y2-y1)
	if err != nil {
		log.Printf("Shot CaptureDisplay err:%v", err)
		return nil, err
	}

	resized := resize.Resize(uint(w), uint(h), shotImg, resize.NearestNeighbor)

	return encode(resized)
}

func ShotAndResize(w, h int) ([]byte, error) {

	shotImg, err := screenshot.CaptureDisplay(0)
	if err != nil {
		log.Printf("Shot CaptureDisplay err:%v", err)
		return nil, err
	}

	resized := resize.Resize(uint(w), uint(h), shotImg, resize.NearestNeighbor)

	return encode(resized)
}

func encode(img image.Image) ([]byte, error) {

	buf := bytes.NewBufferString("")

	err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Printf("GetScreenshotWithMouse Encode err:%v", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func DrawScreenOnly() ([]byte, error) {

	//1.截屏
	shotImg, err := screenshot.CaptureDisplay(0)
	if err != nil {
		log.Printf("GetScreenshotWithMouse CaptureDisplay err:%v", err)
		return nil, err
	}

	buf := bytes.NewBufferString("")

	err = jpeg.Encode(buf, shotImg, &jpeg.Options{Quality: 30})
	if err != nil {
		log.Printf("GetScreenshotWithMouse Encode err:%v", err)
		return nil, err
	}

	return buf.Bytes(), nil

}
