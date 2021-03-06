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
		Boxid int `json:"boxid"`
		X1 int `json:"x1"`
		Y1 int `json:"y1"`
		X2 int `json:"x2"`
		Y2 int `json:"y2"`
		Result  string `json:"result"`
	}
//json format of person recognition
type imgJson struct {
	Imgid   int    `json:"imgid"`
	Imgpath string `json:"imgpath"`
	Bboxs    []Bbox `json:"bbox"`

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
	bboxid:=0
	x1:=1
	y1:=1
	x2:=3
	y2:=3
	imgjson:=imgJson{
		Imgid: imgId,
		Imgpath: imgPath,
		Bboxs:[]Bbox{
			      Bbox{
				Boxid:bboxid,
				X1:x1,
				Y1:y1,
				X2:x2,
				Y2:y2,
			        Result:"dummy-result",
			      },
			},	   
	}
	if jsonString, err := json.Marshal(imgjson); err == nil {
		output := &Output{AgeJson: string(jsonString)}
		err = ctx.SetOutputObject(output)
		if err != nil {
		   return true, err
	        }
	}


	return true, nil
}
