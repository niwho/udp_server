package message

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func Test_dingding(t *testing.T) {
	ddmsg := "# 啦啦啦\n\n##### gogogo\n![screenshot](http://pic27.nipic.com/20130317/6608733_093917444000_2.jpg)"
	ddmsg = "# 测试测试\n#####  gggggg"
	SendDD(ddmsg, []string{"17388935273"})
}
