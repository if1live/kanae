package kanaelib

import "github.com/kardianos/osext"

func Check(e error) {
	if e != nil {
		//raven.CaptureErrorAndWait(e, nil)
		panic(e)
	}
}

func GetExecutablePath() string {
	path, err := osext.ExecutableFolder()
	Check(err)
	return path
}
