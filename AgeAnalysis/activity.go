package sample

import (
	"github.com/project-flogo/core/activity"
	"encoding/json"
//	"fmt"
//	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata( &Input{}, &Output{})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {


	act := &Activity{} //add aSetting to instance

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
}
type Bbox struct {
		boxid int `json:"boxid"`
		x1 int `json:"x1"`
		y1 int `json:"y1"`
		x2 int `json:"x2"`
		y2 int `json:"y2"`
		result  string `json:"result"`
	}
//json format of person recognition
type imgJson struct {
	imgid   int    `json:"imgid"`
	imgpath string `json:"imgpath"`
	bboxs    []Bbox `json:"bbox"`

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
//recognition done here, dummy now

	imgId:=215
	imgPath:="/home/test.jpg/"
	Bboxid:=0
	X1:=1
	Y1:=1
	X2:=3
	Y2:=3
	imgjson:=imgJson{
		imgid: imgId
		imgpath: imgPath
		bboxs:[
			Bbox{
				boxid:Bboxid
				x1:X1
				y1:Y1
				x2:X2
				y2:Y2
				result:"dummy-result"
			}
		]	   
	}
	if jsonString, err := json.Marshal(config); err == nil {
		output := &Output{AgeJson: string(jsonString)}
		err = ctx.SetOutputObject(output)
	}
	if err != nil {
		return true, err
	}

	return true, nil
}
