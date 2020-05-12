package game

const (
	None       = 0
	Player1Won = 1
	Player2Won = 2
	Draw       = 3
)

//func (outcome *Outcome) UnmarshalJSON(b []byte) error {
//	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
//	type OC Outcome
//	var r *OC = (*OC)(outcome)
//	err := json.Unmarshal(b, &r)
//	if err != nil{
//		panic(err)
//	}
//	switch *outcome {
//	case None, Player1Won, Player2Won, Draw:
//		return nil
//	}
//	return errors.New("invalid outcome type")
//}
