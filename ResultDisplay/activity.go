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

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"

	"image"
	"image/color"
	"strings"

	"gocv.io/x/gocv"
)

var model *tf.SavedModel
var activityMd = activity.ToMetadata(&Input{})
var ageStage [3]string
var maxValueIndex int
var age string
var window = gocv.NewWindow("Age")
var textColor = color.RGBA{0, 255, 0, 0}
var pt = image.Pt(20, 20)
var left, top, right, bottom int

func init() {
	_ = activity.Register(&Activity{})
	//activity.Register(&Activity{}, New) to create instances using factory method 'New'
	// add by Yongtao
	// ageStage = [...]string{"0-6", "8-20", "25-32", "38-53", "60-"}
	ageStage = [...]string{"0-13", "25-40", "60-"}
	var err error
	model, err = tf.LoadSavedModel("resource/ageModel2", []string{"serve"}, nil)
	if err != nil {
		log.Fatal(err)
	}
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
	//recognition done here, dummy now

	// *************************
	// imgName := "tmpAge.jpg"
	receiveString := input.Serial
	faceArr := strings.Split(receiveString, ";")
	framePath := faceArr[0]

	// ***************************************
	if exists(framePath) {
		for faceIndex := 1; faceIndex < len(faceArr); faceIndex++ {
			img := gocv.IMRead(framePath, gocv.IMReadColor)
			rectString := strings.Replace(faceArr[faceIndex], "(", "", -1)
			rectString = strings.Replace(rectString, ")", "", -1)
			rectString = strings.Replace(rectString, "-", ",", -1)
			rectArr := strings.Split(rectString, ",")
			// fmt.Println("***************************************************")
			// fmt.Println(rectArr)
			left, err = strconv.Atoi(rectArr[0])
			if err != nil {
				return true, err
			}
			top, err = strconv.Atoi(rectArr[1])
			if err != nil {
				return true, err
			}
			right, err = strconv.Atoi(rectArr[2])
			if err != nil {
				return true, err
			}
			bottom, err = strconv.Atoi(rectArr[3])
			if err != nil {
				return true, err
			}
			left -= 20
			top -= 60
			right += 20
			bottom += 20
			if right > 1280 {
				right = 1280
			}
			if bottom > 720 {
				bottom = 720
			}

			rect := image.Rect(left, top, right, bottom)
			imgFace := img.Region(rect)
			gocv.IMWrite("resource/temp/tmpAge.jpg", imgFace)
			imgName := "resource/temp/tmpAge.jpg"

			imageFile, err := os.Open(imgName)
			if err != nil {
				log.Fatal(err)
			}
			var imgBuffer bytes.Buffer
			io.Copy(&imgBuffer, imageFile)
			imgtf, err := readImage(&imgBuffer, "jpg")
			if err != nil {
				log.Fatal("error making tensor: ", err)
			}

			result, err := model.Session.Run(
				map[tf.Output]*tf.Tensor{
					model.Graph.Operation("input_1").Output(0): imgtf,
				},
				[]tf.Output{
					model.Graph.Operation("dense_2/Softmax").Output(0),
				},
				nil,
			)

			if err != nil {
				log.Fatal(err)
			}

			if preds, ok := result[0].Value().([][]float32); ok {
				// 		fmt.Println(preds[0])
				// 		fmt.Println(reflect.TypeOf(preds[0]))
				maxValueIndex = indexOfMax(preds[0])
				age = ageStage[maxValueIndex]
				// fmt.Println("Age: ", age)
				fmt.Printf("\n %c[%d;%d;%dm%s%c[0m\n", 0x1B, 0, 0, 32, age, 0x1B)

				imgFace := gocv.IMRead(imgName, gocv.IMReadColor)
				gocv.PutText(&imgFace, age, pt, gocv.FontHersheyPlain, 1.2, textColor, 2)
				window.IMShow(imgFace)
				window.WaitKey(1)
			}

		}
	}

	// *******************************
	// fmt.Printf("Input serial: %s\n", input.Serial)
	fmt.Printf("\n %c[%d;%d;%dmInput serial: %s%c[0m\n", 0x1B, 0, 0, 31, input.Serial, 0x1B)

	ctx.Logger().Debugf("Input serial: %s", input.Serial)
	// 	ctx.Logger().Debugf("Age: %s", age)

	return true, nil
}

// add by Yongtao
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
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

func readImage(imageBuffer *bytes.Buffer, imageFormat string) (*tf.Tensor, error) {
	tensor, err := tf.NewTensor(imageBuffer.String())
	if err != nil {
		return nil, err
	}
	graph, input, output, err := transformGraph(imageFormat)
	if err != nil {
		return nil, err
	}
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, err
	}
	return normalized[0], nil
}

func transformGraph(imageFormat string) (graph *tf.Graph, input,
	output tf.Output, err error) {
	const (
		// H, W  = 224, 224
		H, W  = 227, 227
		Mean  = float32(117)
		Scale = float32(1)
	)
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)

	var decode tf.Output
	switch imageFormat {
	case "png":
		decode = op.DecodePng(s, input, op.DecodePngChannels(3))
	case "jpg",
		"jpeg":
		decode = op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))
	default:
		return nil, tf.Output{}, tf.Output{},
			fmt.Errorf("imageFormat not supported: %s", imageFormat)
	}

	output = op.Div(s,
		op.Sub(s,
			op.ResizeBilinear(s,
				op.ExpandDims(s,
					op.Cast(s, decode, tf.Float),
					op.Const(s.SubScope("make_batch"), int32(0))),
				op.Const(s.SubScope("size"), []int32{H, W})),
			op.Const(s.SubScope("mean"), Mean)),
		op.Const(s.SubScope("scale"), Scale))
	graph, err = s.Finalize()
	return graph, input, output, err
}

func indexOfMax(arr []float32) int {

	//Get the maximum value in an array and get the index

	//Declare an array
	// var arr [5]int = [...]int{6, 45, 63, 16, 86}
	//Suppose the first element is the maximum value and the index is 0.
	maxVal := arr[0]
	maxIndex := 0

	for i := 1; i < len(arr); i++ {
		//Cycle comparisons from the second element, exchange if found to be larger
		if maxVal < arr[i] {
			maxVal = arr[i]
			maxIndex = i
		}
	}

	return maxIndex
}
