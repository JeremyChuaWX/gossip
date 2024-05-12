package urlquery

import (
	"net/url"
	"testing"
)

func TestParseStruct(t *testing.T) {
	type querystruct struct {
		Param1 string `query:"param1"`
		Param2 string `query:"param2"`
	}
	url, _ := url.Parse("https://google.com?param1=var1&param2=var2")
	res, err := ParseStruct[querystruct](url)
	if err != nil {
		t.Fatal("failed to parse query struct", err)
	}
	if res.Param1 != "var1" || res.Param2 != "var2" {
		t.Fatal("failed to parse query struct", err)
	}
}
