package utils

import (
    "bytes"
    "encoding/json"
    "encoding/xml"
    "golang.org/x/net/html/charset"
    "golang.org/x/text/encoding/simplifiedchinese"
    "golang.org/x/text/transform"
    "io/ioutil"
)

func ToJSONString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func ToPrettyString(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "    ")
	return string(b)
}

func GbkToUtf8(s []byte) ([]byte, error) {
    reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
    d, e := ioutil.ReadAll(reader)
    if e != nil {
        return s, e
    }
    return d, nil
}

func DecodeGbk(v interface{}, body []byte) error {
    bodyBytes, err := GbkToUtf8(body)
    if err != nil {
        return err
    }
    decoder := xml.NewDecoder(bytes.NewReader(bodyBytes))
    decoder.CharsetReader = charset.NewReaderLabel
    err = decoder.Decode(v)
    return err
}
