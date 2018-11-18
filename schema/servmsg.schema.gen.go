package schema

import (
	"io"
	"time"
	"unsafe"
)

var (
	_ = unsafe.Sizeof(0)
	_ = io.ReadFull
	_ = time.Now()
)

type ReqMsg struct {
	Id     string
	Params []byte
}

func (d *ReqMsg) Size() (s uint64) {

	{
		l := uint64(len(d.Id))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Params))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	return
}
func (d *ReqMsg) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Id))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Id)
		i += l
	}
	{
		l := uint64(len(d.Params))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Params)
		i += l
	}
	return buf[:i+0], nil
}

func (d *ReqMsg) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Id = string(buf[i+0 : i+0+l])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Params)) >= l {
			d.Params = d.Params[:l]
		} else {
			d.Params = make([]byte, l)
		}
		copy(d.Params, buf[i+0:])
		i += l
	}
	return i + 0, nil
}

type CodeError struct {
	Code    string
	MsgBody string
	Prefix  string
}

func (d *CodeError) Size() (s uint64) {

	{
		l := uint64(len(d.Code))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.MsgBody))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Prefix))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	return
}
func (d *CodeError) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Code))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Code)
		i += l
	}
	{
		l := uint64(len(d.MsgBody))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.MsgBody)
		i += l
	}
	{
		l := uint64(len(d.Prefix))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Prefix)
		i += l
	}
	return buf[:i+0], nil
}

func (d *CodeError) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Code = string(buf[i+0 : i+0+l])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.MsgBody = string(buf[i+0 : i+0+l])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Prefix = string(buf[i+0 : i+0+l])
		i += l
	}
	return i + 0, nil
}

type Result struct {
	ResultRef []byte
	Err       *CodeError
}

func (d *Result) Size() (s uint64) {

	{
		l := uint64(len(d.ResultRef))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		if d.Err != nil {

			{
				s += (*d.Err).Size()
			}
			s += 0
		}
	}
	s += 1
	return
}
func (d *Result) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.ResultRef))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.ResultRef)
		i += l
	}
	{
		if d.Err == nil {
			buf[i+0] = 0
		} else {
			buf[i+0] = 1

			{
				nbuf, err := (*d.Err).Marshal(buf[i+1:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}
			i += 0
		}
	}
	return buf[:i+1], nil
}

func (d *Result) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.ResultRef)) >= l {
			d.ResultRef = d.ResultRef[:l]
		} else {
			d.ResultRef = make([]byte, l)
		}
		copy(d.ResultRef, buf[i+0:])
		i += l
	}
	{
		if buf[i+0] == 1 {
			if d.Err == nil {
				d.Err = new(CodeError)
			}

			{
				ni, err := (*d.Err).Unmarshal(buf[i+1:])
				if err != nil {
					return 0, err
				}
				i += ni
			}
			i += 0
		} else {
			d.Err = nil
		}
	}
	return i + 1, nil
}

type ResultItem struct {
	Key    string
	Result *Result
}

func (d *ResultItem) Size() (s uint64) {

	{
		l := uint64(len(d.Key))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		if d.Result != nil {

			{
				s += (*d.Result).Size()
			}
			s += 0
		}
	}
	s += 1
	return
}
func (d *ResultItem) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Key))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Key)
		i += l
	}
	{
		if d.Result == nil {
			buf[i+0] = 0
		} else {
			buf[i+0] = 1

			{
				nbuf, err := (*d.Result).Marshal(buf[i+1:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}
			i += 0
		}
	}
	return buf[:i+1], nil
}

func (d *ResultItem) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Key = string(buf[i+0 : i+0+l])
		i += l
	}
	{
		if buf[i+0] == 1 {
			if d.Result == nil {
				d.Result = new(Result)
			}

			{
				ni, err := (*d.Result).Unmarshal(buf[i+1:])
				if err != nil {
					return 0, err
				}
				i += ni
			}
			i += 0
		} else {
			d.Result = nil
		}
	}
	return i + 1, nil
}

type ResMsg struct {
	Id       string
	Response []*ResultItem
}

func (d *ResMsg) Size() (s uint64) {

	{
		l := uint64(len(d.Id))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Response))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}

		for k0 := range d.Response {

			{
				if d.Response[k0] != nil {

					{
						s += (*d.Response[k0]).Size()
					}
					s += 0
				}
			}

			s += 1

		}

	}
	return
}
func (d *ResMsg) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.Id))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.Id)
		i += l
	}
	{
		l := uint64(len(d.Response))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		for k0 := range d.Response {

			{
				if d.Response[k0] == nil {
					buf[i+0] = 0
				} else {
					buf[i+0] = 1

					{
						nbuf, err := (*d.Response[k0]).Marshal(buf[i+1:])
						if err != nil {
							return nil, err
						}
						i += uint64(len(nbuf))
					}
					i += 0
				}
			}

			i += 1

		}
	}
	return buf[:i+0], nil
}

func (d *ResMsg) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Id = string(buf[i+0 : i+0+l])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Response)) >= l {
			d.Response = d.Response[:l]
		} else {
			d.Response = make([]*ResultItem, l)
		}
		for k0 := range d.Response {

			{
				if buf[i+0] == 1 {
					if d.Response[k0] == nil {
						d.Response[k0] = new(ResultItem)
					}

					{
						ni, err := (*d.Response[k0]).Unmarshal(buf[i+1:])
						if err != nil {
							return 0, err
						}
						i += ni
					}
					i += 0
				} else {
					d.Response[k0] = nil
				}
			}

			i += 1

		}
	}
	return i + 0, nil
}
