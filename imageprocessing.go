package crawler

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"

	"gocv.io/x/gocv"
)

// CheckGraphImage judges wheter the image contains graph or not
//
// Parameters:
//     urlString - crawled image url source
//     Returns bool value
func CheckGraphImage(urlString string) bool {
	img := UrlToImage(urlString)
	if img.Empty() {
		return false
	}

	height := img.Size()[0]
	width := img.Size()[1]
	var interpolation gocv.InterpolationFlags
	if height < 640 || width < 480 {
		interpolation = gocv.InterpolationCubic
	} else {
		interpolation = gocv.InterpolationArea
	}

	gocv.Resize(img, &img, image.Point{640, 480}, 0, 0, interpolation)
	houghImg := HoughTransform(img)
	entropy := CalcImageEntropy(img, 32)
	hist := CalcImageHistogram(img)
	_, maxVal, _, maxIdx := gocv.MinMaxLoc(hist)

	if houghImg.Empty() {
		return false
	}

	if (maxIdx.Y > 200 || maxIdx.Y < 50) && maxVal > 120000 {
		if entropy < 2.8 {
			return true
		}
	}
	return false
}

// CalcImageHistogram calculate histogram of image(gray-scale)
//
// Parameters:
//     src - source image matrix
//     Returns histogram matrix
func CalcImageHistogram(src gocv.Mat) gocv.Mat {
	img := src.Clone()
	defer img.Close()

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	hist := gocv.NewMat()
	mask := gocv.NewMat()
	defer mask.Close()

	channels := []int{0}
	size := []int{256}
	ranges := []float64{0, 256}
	acc := false

	gocv.CalcHist([]gocv.Mat{gray}, channels, mask, &hist, size, ranges, acc)

	return hist
}

// DrawHistogram process image to show easily in NewWindow
//
// Parameters:
//     hist - source image matrix
//     Returns image matrix
func DrawHistogram(hist gocv.Mat) gocv.Mat {

	dHist := gocv.NewMat()
	defer dHist.Close()

	//set matrix size to be shown
	histW := 512
	histH := 400
	size := hist.Size()[0]
	binW := int(float64(histW) / float64(size))
	histImage := gocv.NewMatWithSize(512, 400, gocv.MatTypeCV8U)

	gocv.Normalize(hist, &dHist, 0, float64(histImage.Rows()), gocv.NormMinMax) // normalize to show easily

	for idx := 1; idx < size; idx++ {
		gocv.Line(&histImage, image.Point{binW * (idx - 1), histH - int(dHist.GetFloatAt(idx-1, 0))}, image.Point{binW * (idx), histH - int(dHist.GetFloatAt(idx, 0))}, color.RGBA{255, 255, 255, 0}, 2)
	}

	return histImage
}

// HoughTransform extract edges using canny edge detector and find lines from image(gray-scale)
//
// Parameters:
//     src - source image matrix
//     Returns image matrix containing lines
func HoughTransform(src gocv.Mat) gocv.Mat {
	img := gocv.NewMat()
	img = src.Clone()

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	edges := gocv.NewMat()
	defer edges.Close()

	var weakThreshold float32 = 300
	var strongThreshold float32 = 400
	gocv.Canny(gray, &edges, weakThreshold, strongThreshold)

	lines := gocv.NewMat()
	defer lines.Close()

	var rho float32 = 1
	var theta float32 = math.Pi / 180
	var threshold int = 100
	var minLineLength float32 = 100
	var maxLineGap float32 = 5
	if !edges.Empty() {
		gocv.HoughLinesPWithParams(edges, &lines, rho, theta, threshold, minLineLength, maxLineGap)
	}

	if !lines.Empty() {
		for idx := 0; idx < lines.Rows(); idx++ {
			line := lines.GetVeciAt(idx, 0)
			gocv.Line(&img, image.Point{int(line[0]), int(line[1])}, image.Point{int(line[2]), int(line[3])}, color.RGBA{0, 255, 0, 0}, 2)
		}
	} else {
		//if there are no lines, return Empty Mat
		img = gocv.NewMat()
	}

	return img
}

// CalcImageEntropy calculates image(gray-scale) entropy
//
// Parameters:
//     src - source image matrix
//     stride - window size which makes sub image from source
//     Returns entropy value by float64
func CalcImageEntropy(src gocv.Mat, stride int) float64 {

	gray := gocv.NewMat()
	if src.Channels() != 1 {
		gocv.CvtColor(src, &gray, gocv.ColorBGRToGray)
	} else {
		gray = src.Clone()
	}

	shape := gray.Size()
	row := shape[0]
	col := shape[1]
	numBlock := int(row * col / (stride * stride))
	size := float64(stride * stride)

	var entropy float64
	entropy = 0

	for i := 0; i < row; i += stride {
		for j := 0; j < col; j += stride {
			//image 에서 x는 col, y는 row!
			croppedMat := gray.Region(image.Rect(j, i, j+stride, i+stride))
			subImg := croppedMat.Clone()
			hist := gocv.NewMat()
			defer croppedMat.Close()
			defer subImg.Close()
			defer hist.Close()
			gocv.CalcHist([]gocv.Mat{subImg}, []int{0}, gocv.NewMat(), &hist, []int{256}, []float64{0, 256}, false)

			for k := 0; k < 256; k++ {
				bin := float64(hist.GetFloatAt(k, 0))
				if bin > 0 {
					entropy += bin / size * math.Log2(size/bin)
				}
			}
		}
	}
	return entropy / float64(numBlock)
}

// UrlToImage read url string and convert to image based on MAT which is basic type on OpenCV
//
// Parameters:
//     urlString - image link
//     Returns the image matrix
func UrlToImage(urlString string) gocv.Mat {

	//url check whether it is valid or not
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		log.Println(err)
		return gocv.NewMat()
	}

	u, err := url.Parse(urlString)
	if err != nil || u.Scheme == "" || u.Host == "" {
		log.Println(err)
		return gocv.NewMat()
	}

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		log.Println(err)
		return gocv.NewMat()
	}
	//HTTP Error 403: Forbidden problem -> set header!
	req.Header.Add("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return gocv.NewMat()
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return gocv.NewMat()
	}

	img, err := gocv.IMDecode(bodyBytes, gocv.IMReadColor)
	if img.Empty() {
		log.Println(err)
		return gocv.NewMat()
	}

	return img
}
