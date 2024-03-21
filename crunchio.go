package crunchio

import (
	"github.com/superwhiskers/crunch/v3"
)

//Bytes requires a type to be able to represent itself as a byte slice
type Bytes interface {
	Bytes() []byte
}

type Buffer struct {
	crunch.Buffer
}

func NewBuffer(slices ...[]byte) (*Buffer) {
	return &Buffer{crunch.NewBuffer(slices...)}
}

//SetValue changes the data of the value using a supported input type
//Data must be one of the following types:
//- byte, bool, int, uint
//- []byte
//- string
//- int16, int32, int64
//- uint16, uint32, uint64
//- float32, float64
//Alternatively, data may satisfy the Bytes interface by providing a `Bytes() []byte` (such as another Buffer)
func (buf *Buffer) WriteAbstract(data interface{}) error {
	switch data.(type) {
	case Bytes:
		buf.Write(data.(Bytes).Bytes())
	case *Bytes:
		buf.Write(data.(*Bytes).Bytes())
	case byte, bool, int, uint:
		buf.Write([]byte{data.(byte)})
	case []byte:
		buf.Write(data.([]byte))
	case string:
		buf.Write([]byte(data.(string)))
	case int16:
		buf.WriteI16LENext([]int16{data.(int16)})
	case int32:
		buf.WriteI32LENext([]int32{data.(int32)})
	case int64:
		buf.WriteI64LENext([]int64{data.(int64)})
	case uint16:
		buf.WriteU16LENext([]uint16{data.(uint16)})
	case uint32:
		buf.WriteU32LENext([]uint32{data.(uint32)})
	case uint64:
		buf.WriteU64LENext([]uint64{data.(uint64)})
	case float32:
		buf.WriteF32LENext([]float32{data.(float32)})
	case float64:
		buf.WriteF64LENext([]float64{data.(float64)})
	default:
		return fmt.Errorf("crunchio: unsupported abstract type")
	}
	return nil
}

//Write implements io.Writer
func (buf *Buffer) Write(p []byte) (int, error) {
	if p == nil || len(p) <= 0 {
		return 0, nil
	}
	buf.Grow(int64(len(p)))
	buf.WriteBytesNext(p)
	return len(p), buf.sync()
}

//Read implements io.Reader
func (buf *Buffer) Read(p []byte) (int, error) {
	if p == nil || len(p) <= 0 {
		return 0, nil
	}
	bytes := buf.ReadBytesNext(int64(len(p)))
	for i := 0; i < len(bytes); i++ {
		p[i] = bytes[i]
	}
	return len(bytes), nil
}

//Seek implements io.Seeker
func (buf *Buffer) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		buf.SeekByte(offset, false)
	case io.SeekCurrent:
		buf.SeekByte(offset, true)
	case io.SeekEnd:
		buf.SeekByte(buf.ByteCapacity()-offset, false)
	default:
		return 0, ERROR_SEEK_WHENCE
	}
	return buf.ByteOffset(), nil
}

//Size returns the byte size of the buffer
func (buf *Buffer) Size() int {
	return int(buf.ByteCapacity())
}

//Reset wipes the buffer
func (buf *Buffer) Reset() {
	buf.Reset()
}

//String returns a string representation of the buffer
func (buf *Buffer) String() string {
	return string(buf.Bytes())
}
