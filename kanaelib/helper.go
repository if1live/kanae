package kanaelib

func check(e error) {
	if e != nil {
		//raven.CaptureErrorAndWait(e, nil)
		panic(e)
	}
}
