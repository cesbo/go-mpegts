package mpegts

import "fmt"

// Desc_0A ISO 639 language descriptor

type Desc_0A_Info struct {
	LanguageCode [3]byte
	AudioType    byte
}
type Desc_0A struct {
	Info []Desc_0A_Info
}

func (d *Desc_0A) String() string {
	return "0x0A ISO 639 language descriptor"
}
func (d *Desc_0A) InfoString() string {
	str := ""
	for _, info := range d.Info {
		str += fmt.Sprintf("LanguageCode: %s, AudioType: %d\n", string(info.LanguageCode[:]), info.AudioType)
	}
	return str
}

func (d *Desc_0A) Encode() (desc Descriptors) {
	lenght := 2 + len(d.Info)*4 // each info is 4 bytes
	desc = make(Descriptors, lenght)
	desc[0] = 0x0A
	desc[1] = byte(lenght - 2)
	for i := 0; i < len(d.Info); i++ {
		copy(desc[2+i*4:], d.Info[i].LanguageCode[:])
		desc[2+i*4+3] = d.Info[i].AudioType
	}
	return desc
}

func (d *Desc_0A) Decode(desc Descriptors) error {
	length := 0
	if len(desc) >= 2 && desc[0] == 0x0A {
		length = int(desc[1])
	} else {
		return ErrDescriptorFormat
	}
	if length != 0 && length%4 == 0 && len(desc) >= length+2 {
		num_info := length / 4
		d.Info = make([]Desc_0A_Info, num_info)
		for i := 0; i < num_info; i++ {
			copy(d.Info[i].LanguageCode[:], desc[2+i*4:2+i*4+3])
			d.Info[i].AudioType = desc[2+i*4+3]
		}
		return nil
	} else {
		return ErrDescriptorFormat
	}
}
