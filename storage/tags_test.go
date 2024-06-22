package storage

import (
	"RealWorldWeb/utils"
	"testing"
)

func TestListTags(t *testing.T) {
	res, err := ListPopularTags()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("tags: %v\n", utils.JsonMarshal(res))
}
