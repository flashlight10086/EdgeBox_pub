package sample

import (
	"github.com/project-flogo/core/activity"
	"encoding/json"
//	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Output{})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {

//	s := &Settings{}
//	err := metadata.MapToStruct(ctx.Settings(), s, true)
//	if err != nil {
//		return nil, err
//	}

//	ctx.Logger().Debugf("Setting: %s", s.ASetting)

	act := &Activity{} //add aSetting to instance//nothing to add now

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}
//bounding box by form of x1,y1,x2,y2
type Bbox struct {
		Boxid int `json:"boxid"`
		X1 int `json:"x1"`
		Y1 int `json:"y1"`
		X2 int `json:"x2"`
		Y2 int `json:"y2"`
	}
//json format of person recognition
type imgJson struct {
	Imgid   int    `json:"imgid"`
	Imgpath string `json:"imgpath"`
	Bboxs    []Bbox `json:"bbox"`
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	//call neural network here
        ctx.Logger().Debugf("result of picking out a person: %s", "found") //log is also dummy here
	err = nil //set if neural network go wrong
	if err != nil {
		return true, err
	}
//dummy json generation here
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
		Bboxs:[
			Bbox{
				Boxid:bboxid,
				X1:x1,
				Y1:y1,
				X2:x2,
				Y2:y2,
			}
		]	   
	}
	if jsonString, err := json.Marshal(config); err == nil {
           fmt.Println("================struct to json str==")
           fmt.Println(string(jsonString))
		output := &Output{Serial: string(jsonString)}
        }
	//output := &Output{Serial: "1"}//should be serial of the record in the database
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
