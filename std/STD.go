package std




// Returns the complete STD objects
// These items can be accessed with including `stdkbm`.
// An example would look like
//	include("stdkbm")
func STD() map[string]any {
	return map[string]any{
		"log":				Log,
		"logf":				Logf,

		"type":				Type,
		"print":		 	Print,
		"printf":		 	Printf,
		"sprint":		 	Sprint,
		"sprintf":		 	Sprintf,

		"len":			 	Len,
		"StringToInt":	 	StringToInt,
		"NumberToString":	NumberToString,
		"BooleanToString":	BooleanToString,
		"BooleanToInt":		BooleanToInt,
		"IntToBoolean":		IntToBoolean,
		"StringToBoolean":	StringToBoolean,
	}
}