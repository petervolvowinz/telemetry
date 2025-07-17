package datapoint

type DataPoint struct {
	BikeID       string  `csv:"Bike_id" json:"Bike_id"`
	Timestamp    int64   `csv:"Timestamp" json:"Timestamp"`
	Latitude     float64 `csv:"Latitude" json:"Latitude"`
	Longitude    float64 `csv:"Longitude" json:"Longitude"`
	BatteryLevel float64 `csv:"Battery_level" json:"Battery_level"`
	Charging     bool    `csv:"Charging" json:"Charging"`
}

type getRequest struct {
	reply chan DataPoint
}

type setRequest struct {
	val DataPoint
}

type CurrentDataPoint struct {
	getChan chan getRequest
	setChan chan setRequest
}

func NewCurrentDataPoint() *CurrentDataPoint {
	c := &CurrentDataPoint{
		getChan: make(chan getRequest),
		setChan: make(chan setRequest),
	}
	go func() {
		var value DataPoint
		for {
			select {
			case set := <-c.setChan:
				value = set.val
			case get := <-c.getChan:
				get.reply <- value
			}
		}
	}()
	return c
}

func (c *CurrentDataPoint) SetCurrentDataPoint(val DataPoint) {
	c.setChan <- setRequest{val: val}
}

func (c *CurrentDataPoint) GetCurrentDataPoint() DataPoint {
	reply := make(chan DataPoint)
	c.getChan <- getRequest{reply: reply}
	return <-reply
}
