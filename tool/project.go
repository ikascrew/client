package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ikascrew/ikasbox/handler"

	"github.com/nfnt/resize"
	"golang.org/x/xerrors"
	"gopkg.in/cheggaaa/pb.v1"
)

const CLIENT = ".client"
const IMAGE = "images"

func GetClientDir() string {
	return CLIENT
}

const Ikasbox = "http://localhost:5555"
const ListAPI = "/project/content/list/%d"

//const Base = "http://10.0.0.1:5555/static/images/thumb/"
const Base = Ikasbox + "/static/images/thumb/"

func CreateProject(id string) error {

	p, err := strconv.Atoi(id)
	if err != nil {
		return xerrors.Errorf("Project id: %w", err)
	}

	images := filepath.Join(GetClientDir(), IMAGE)
	os.RemoveAll(images)

	err = Mkdir(images)
	if err != nil {
		return xerrors.Errorf("Error make directory: %w", err)
	}

	url := fmt.Sprintf("%s%s", Ikasbox, fmt.Sprintf(ListAPI, p))
	contents, err := getContentList(url)
	if err != nil {
		return xerrors.Errorf("Error get content list[%s]: %w", url, err)
	}

	bar := pb.StartNew(len(contents)).Prefix("Create Thumbnail")

	for _, elm := range contents {

		content_id := elm.ContentID
		//3つのファイルにアクセスして保存
		err = create(content_id)
		if err != nil {
			return xerrors.Errorf("Error Create : %w", err)
		}

		bar.Increment()
	}
	bar.FinishPrint("Thumbnail Completion")

	return nil
}

func create(id int) error {

	imageDir := filepath.Join(GetClientDir(), IMAGE)
	url := fmt.Sprintf(Base+"%d.jpg", id)

	img, err := downloadImage(url)
	if err != nil {
		return xerrors.Errorf("downloadImage : %w", err)
	}

	out := fmt.Sprintf(imageDir+"/%d.jpg", id)

	new_image := resize.Resize(320, 0, img, resize.Lanczos3)

	err = writeImage(new_image, out)
	if err != nil {
		return xerrors.Errorf("writeImage : %w", err)
	}

	url = fmt.Sprintf(Base+"/%d_4.jpg", id)

	img1, err := downloadImage(url)
	if err != nil {
		return xerrors.Errorf("downloadImage : %w", err)
	}

	url = fmt.Sprintf(Base+"/%d_12.jpg", id)

	img3, err := downloadImage(url)
	if err != nil {
		return xerrors.Errorf("downloadImage : %w", err)
	}

	cut := 3
	images := make([]image.Image, cut)

	images[0] = img1
	images[1] = img
	images[2] = img3

	thumbRect := image.Rect(0, 0, 1280*cut, 720)
	thumb := image.NewRGBA(thumbRect)
	left := 0
	for _, wk := range images {

		rect := image.Rect(left, 0, 1280+left, 720)
		draw.Draw(thumb, rect, wk, image.Pt(0, 0), draw.Over)

		left += 1280
	}

	new_thumb := resize.Resize(320, 0, thumb, resize.Lanczos3)

	out = fmt.Sprintf(imageDir+"/%d_thumb.jpg", id)

	err = writeImage(new_thumb, out)
	if err != nil {
		return xerrors.Errorf("writeImage : %w", err)
	}

	//mat, err := gocv.ImageToMatRGB(thumb)
	//gocv.Resize(mat, &resize, image.Point{}, 0.166, 0.166, gocv.InterpolationDefault)
	//gocv.IMWrite(out, resize)
	return nil
}

func writeImage(img image.Image, out string) error {
	of, err := os.Create(out)
	if err != nil {
		return err
	}
	defer of.Close()

	err = jpeg.Encode(of, img, nil)
	if err != nil {
		return err
	}
	return nil
}

func downloadImage(url string) (image.Image, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, xerrors.Errorf("http get(%s): %w", url, err)
	}
	defer resp.Body.Close()

	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("read all(%s): %w", url, err)
	}

	img, err := jpeg.Decode(bytes.NewReader(byt))
	if err != nil {
		return nil, xerrors.Errorf("decode(%s) : %w", url, err)
	}

	return img, nil
}

func getContentList(url string) ([]*handler.ProjectContent, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dst := handler.ProjectResponse{}
	err = json.Unmarshal(byt, &dst)
	if err != nil {
		return nil, err
	}

	return dst.Contents, nil
}
