package sample

import (
	"strconv"

	"github.com/project-flogo/core/activity"
	//"github.com/project-flogo/core/data/metadata"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	// 	"reflect"

	

	"image"
	"image/color"
	"strings"
//	simplejson "github.com/bitly/go-simplejson"

	"gocv.io/x/gocv"
)

var activityMd = activity.ToMetadata(&Input{})
var imgpath string = ""
var imgid = -1
//var content string = ""
var window = gocv.NewWindow("Output")
var content  []string
var textColor = color.RGBA{0, 255, 0, 0}
//var pt = image.Pt(20, 20)
//var left, top, right, bottom int
type Bbox struct {
	Boxid int `json:"boxid"`
	X1    int `json:"x1"`
	Y1    int `json:"y1"`
	X2    int `json:"x2"`
	Y2    int `json:"y2"`
}

//json format of person recognition
type imgJson struct {
	Imgid   int    `json:"imgid"`
	Imgpath string `json:"imgpath"`
	Bboxes  []Bbox `json:"bboxes"`
}

type imgJsonR struct {
	ImgJson imgJson  `json:"imgjson"`
	Content  []string `json:"content"`
}

func init() {
	_ = activity.Register(&Activity{})

}

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {

	act := &Activity{} //add aSetting to instance

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}
        imgjsonS := input.Serial
	imgjsonS = strings.Replace(imgjsonS, "\\\"", "\"", -1)
	imgjson := imgJsonR{}
	json.Unmarshal([]byte(imgjsonS), &imgjson)
	fmt.Println(imgjson)
	
	framePath := imgjson.Imgjson.Imgpath
	imgid_now := imgjson.Imgjson.Imgid
	if (imgid_now!=imgid){
		imgid=imgid_now
		imgPath=framePath
		content=imgjson.Content
		
	},
	else{
		if exists(framePath) {
			imgFace := gocv.IMRead(imgPath, gocv.IMReadColor)
			for faceIndex := 1; faceIndex < len(imgjson.Imgjson.Content); faceIndex++ {
				gocv.Rectangle(&imgFace,image.Pt(imgjson.Imgjson.BBoxes[faceIndex].X1,imgjson.Imgjson.BBoxes[faceIndex].Y1),image.Pt(imgjson.Imgjson.BBoxes[faceIndex].X2,imgjson.Imgjson.BBoxes[faceIndex].Y2),color.RGBA{R: 0, G: 255, B: 0, A: 100}, 1)
				gocv.PutText(&imgFace, content+","+imgjson.Content[faceIndex], image.Pt(imgjson.Imgjson.BBoxes[faceIndex].X1,imgjson.Imgjson.BBoxes[faceIndex].Y1+20), gocv.FontHersheyPlain, 1.2, textColor, 2)
				window.IMShow(img)
			        window.WaitKey(1)	
			}
		
	}
	
	return true, nil
}



// determine if the file/folder of the given path exists
func exists(path string) bool {

	_, err := os.Stat(path)
	//os.Stat get the file information
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}





