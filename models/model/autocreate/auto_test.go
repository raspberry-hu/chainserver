package autocreate

import (
	"github.com/gohouse/converter"
	"testing"
)

func Test_Auto(t *testing.T) {
	converter.NewTable2Struct().
		SavePath("./model.go").
		Dsn("root:CaCa@0501@tcp(52.221.250.7:13306)/caca?charset=utf8").
		Run()
}
