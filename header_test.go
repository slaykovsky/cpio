package cpio

import (
	//"fmt"
	"io"
	"reflect"
	"testing"
	"io/ioutil"
	"os"
	"github.com/slaykovsky/fetcher"
)

const (
	Location string = "http://mirror.yandex.ru"
	ScratchDir string = "/tmp"
	FileName string = "centos/7/os/x86_64/isolinux/initrd.img"
)
func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}
func TestReflection(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "header_test")
	if err != nil {
		t.Fatal(err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	f := fetcher.HTTPFetcher{
		Location: Location,
		ScratchDir: ScratchDir,
	}
	cpioArchivePath, err := f.AcquireFile(FileName)
	if err != nil {
		t.Fatal(err.Error())
	}

	cpioArchive, err := os.Open(cpioArchivePath)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer cpioArchive.Close()
	//defer os.Remove(cpioArchivePath)

	header := &CpioHeader{}
	headerReflection := reflect.ValueOf(header).Elem()
	fields := reflect.Indirect(reflect.ValueOf(header))

	reader := io.Reader(cpioArchive)
	//binaryReader := BinaryReader{reader:reader}

	for i := 0; i < headerReflection.NumField(); i++ {
		t.Log(fields.Type().Field(i).Name)
		if i == 0 {
			magic := make([]byte, 6)
			v, err := reader.Read(magic)
			if err != nil {
				t.Fatal(err.Error())
			}
			if v == 0 {
				t.Fatal(err.Error())
			}
			t.Log("MAGIC: ", CToGoString(magic[:]))

		}
		//res, err := binaryReader.ReadField()
		//if err != nil {
		//	t.Fatal(err.Error())
		//}
		//t.Log("Next: ", res)
		//fmt.Println(headerReflection.Field(i))
		b := make([]byte, 8)
		v, err := reader.Read(b)
		if err != nil {
			t.Fatal(err.Error())
		}
		if v == 0 {
			t.Fatal(err.Error())
		}
		t.Log("FIELD: ", b)
	}
}
