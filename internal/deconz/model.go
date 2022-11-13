package deconz

type Sensor struct {
	Config           map[string]interface{} `json:"config"`
	Ep               int                    `json:"ep"`
	Etag             string                 `json:"etag"`
	Lastseen         string                 `json:"lastseen"`
	Manufacturername string                 `json:"manufacturername"`
	Modelid          string                 `json:"modelid"`
	Name             string                 `json:"name"`
	State            map[string]interface{} `json:"state"`
	Swversion        string                 `json:"swversion"`
	Type             string                 `json:"type"`
	Uniqueid         string                 `json:"uniqueid"`
}

type GetSensorsResponse map[string]Sensor
