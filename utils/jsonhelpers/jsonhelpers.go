package jsonhelpers

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
)

// Fetch: fetches json file containing users array.
// @args file{string}	users file
func Fetch(file string, self interface{}) error {
	b, err := ioutil.ReadFile(file)
	if nil != err {
		return err
	}
	return Unmarshal(string(b), self)
}

// Save: saves users to json file
func Save(file string, self interface{}) error {
	contents, err := Marshal(self)
	if nil != err {
		return err
	}
	return ioutil.WriteFile(file, []byte(contents), 0644)
}

// Unmarshal: json unmarshals string to struct
// @args string
// @return error
func Unmarshal(data string, self interface{}) error {
	return json.Unmarshal([]byte(data), self)
}

// Marshal: json marshals struct
// @return string
// @return error
func Marshal(self interface{}) (string, error) {
	b, err := json.Marshal(self)
	if nil != err {
		return "", err
	}
	return string(b), nil
}

func SaveGzip(file string, self interface{}) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	defer w.Close()
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	return enc.Encode(self)
}

func FetchGzip(file string, self interface{}) error {
	var err error
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return err
	}

	reader, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer reader.Close()

	err = json.NewDecoder(reader).Decode(self)
	if err != nil {
		return err
	}

	return nil
}
