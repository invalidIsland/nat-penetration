package test

import (
	"fmt"
	"nat-penetration/helper"
	"testing"
)

func TestAnalyseToken(t *testing.T) {
	claims, err := helper.AnalyseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IkplcnJ5R3VvIiwiZXhwIjoxNjcwMDcxNjQ0fQ.r_6WKHyIoZX2RoRQSEOL5k8-1i1W0Eot_LojG_5d0IU")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf(">>>> %s\n", claims.Username)
}
