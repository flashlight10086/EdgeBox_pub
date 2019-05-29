package sample

import "github.com/project-flogo/core/data/coerce"


type Output struct {
	serial string `md:"serial"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["serial"])
	o.serial = strVal
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"serial": o.serial,
	}
}
