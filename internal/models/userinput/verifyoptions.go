package userinput

type VerifyOptions struct {
	// Can be:
	//  - empty: default to ./eva.json
	//  - a directory: <dir>/eva.json
	//  - a file: must be eva.json
	Path string
}
