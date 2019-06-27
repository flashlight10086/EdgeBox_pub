package sample

import (
	"github.com/project-flogo/core/activity"
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
		boxid int `json:"boxid"`
		x1 int `json:"x1"`
		y1 int `json:"y1"`
		x2 int `json:"x2"`
		y2 int `json:"y2"`
	}
//json format of person recognition
type imgJson struct {
	imgid   int    `json:"imgid"`
	imgpath string `json:"imgpath"`
	bboxs    []Bbox `json:"bbox"`
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
