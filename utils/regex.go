package utils

import(
	"log"
)

func SetRegex(ost, arch string) string{
	var regex string

	switch {
	case (ost == "darwin" && arch == "amd64"):
		// perl syntax
		// works with github.com/dlclark/regexp2
		//good
		regex = `(?i)(?=.*(?:apple|darwin|macos|mac))(?=.*(?:amd64|x86_64))(?!.*(?:freebsd|netbsd|openbsd|linux|windows|win64|.sha256sum|.sha256|.sbom|checksums|.txt))(?:.*(?:apple|darwin|macos|mac).*?(?:amd64|x86_64)|(?:amd64|x86_64).*?(?:apple|darwin|macos|mac))(?:[^a-z]|$)`
		// rescue = `^(?!.*(?:freebsd|openbsd|netbsd|windows|win64|appimage|sha256sum|linux|arm64|aarch64|amd64|x86_64|.deb|i686|armhf|.rpm|checksums|i386|arm|.sbom)).*?$`

	case (ost == "darwin" && arch == "arm64"):
		regex = `(?i)(?=.*(?:apple|darwin|macos|mac))(?=.*(?:arm64|aarch64))(?!.*(?:freebsd|netbsd|openbsd|linux|windows|win64))(?:.*(?:apple|darwin|macos|mac).*?(?:arm64|aarch64)|(?:arm64|aarch64).*?(?:apple|darwin|macos|mac))(?:[^a-z]|$)`
	
	case (ost == "linux" && arch == "amd64"):
		//good
		regex = `(?i)(?=.*(?:linux))(?=.*(?:amd64|x86_64|64))(?!.*(?:freebsd|netbsd|openbsd|windows|apple|darwin|macos|mac))(?:.*(?:linux).*?(?:amd64|x86_64)|(?:amd64|x86_64).*?(?:linux))(?:[^a-z]|$)`

	case (ost == "linux" && arch == "aarch64"):
		regex = `(?i)(?=.*(?:linux))(?=.*(?:arm64|aarch64))(?!.*(?:freebsd|netbsd|openbsd|windows|apple|darwin|macos|mac))(?:.*(?:linux).*?(?:arm64|aarch64)|(?:arm64|aarch64).*?(?:linux))(?:[^a-z]|$)`
	
	// case (ost == "linux" && arch == "amd64" && appimage == true):
	// 	regex =  something

	default:
		msg1 := "OS or Architecture is not supported or not found in the regex pattern"
		msg2 := "File an issue or make a pull request for your OS and Arch"
		msg3 := "Will only list/download for macOS/Linux for the following arch: x86_64/amd64 and arm64" 
		log.Fatalf("%v\n%v\n%v\n%v", msg1, msg2, msg3)
}
	return regex
	//return mainRegex, rescueRegex
}
