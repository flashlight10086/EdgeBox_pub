package sample

import "github.com/project-flogo/core/data/coerce"


type Input struct {
	Serial string `md:"serial,required"`
}

func (r *Input) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["serial"])
	r.Serial = strVal
	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"serial": r.Serial,
	}
}

type Output struct {
	AgeJson string `md:"genderJson"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["genderJson"])
	o.GenderJson = strVal
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"genderJson": o.GenderJson,
	}
}



