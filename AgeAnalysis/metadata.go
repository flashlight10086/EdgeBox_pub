package sample

import "github.com/project-flogo/core/data/coerce"


type Input struct {
	Serial string `md:"serial,required"`
	AgeAttr string `md:"ageAttr,required"`
}

func (r *Input) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["serial"])
	r.Serial = strVal
	strVal, _ := coerce.ToString(values["ageAttr"])
	r.AgeAttr = strVal
	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"serial": r.Serial,
		"ageAttr": r.AgeAttr,
	}
}

type Output struct {
	AgeJson string `md:"ageJson"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["ageJson"])
	o.AgeJson = strVal
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"ageJson": o.AgeJson,
	}
}



