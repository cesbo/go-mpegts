package mpegts

// Desc_7A enhanced_AC-3_descriptor

type Desc_7A struct {
	len byte

	data []byte

	component_flag       byte
	component_type       byte
	bsid                 byte
	mainid               byte
	asvc                 byte
	substream1           byte
	substream2           byte
	substream3           byte
	additional_info_byte []byte
}

func (d *Desc_7A) String() string {
	return "0x7A enhanced_AC-3_descriptor"
}

func (d *Desc_7A) Encode() (desc Descriptors) {
	desc = make(Descriptors, d.len+2)
	desc[0] = 0x7A
	desc[1] = d.len
	copy(desc[2:], d.data)
	return desc
}

func (d *Desc_7A) Decode(desc Descriptors) error {
	if len(desc) >= 2 && desc[0] == 0x7A {
		d.len = desc[1]
		d.data = make([]byte, d.len)
		copy(d.data, desc[2:])
		if d.len > 0 {
			component_pos := 1
			d.component_flag = d.data[0]

			if d.ComponentYypeFlag() {
				d.component_type = d.data[component_pos]
				component_pos++
			}
			if d.BsidFlag() {
				d.bsid = d.data[component_pos]
				component_pos++
			}
			if d.MainidFlag() {
				d.mainid = d.data[component_pos]
				component_pos++
			}
			if d.AsvcFlag() {
				d.asvc = d.data[component_pos]
				component_pos++
			}
			if d.Substream1Flag() {
				d.substream1 = d.data[component_pos]
				component_pos++
			}
			if d.Substream2Flag() {
				d.substream2 = d.data[component_pos]
				component_pos++
			}
			if d.Substream3Flag() {
				d.substream3 = d.data[component_pos]
				component_pos++
			}
		}
		return nil
	} else {
		return ErrDescriptorFormat
	}
}

func (d *Desc_7A) ComponentYypeFlag() bool {
	return d.component_flag&0b10000000 == 0x80
}
func (d *Desc_7A) BsidFlag() bool {
	return d.component_flag&0b01000000 == 0x80
}
func (d *Desc_7A) MainidFlag() bool {
	return d.component_flag&0b00100000 == 0x80
}
func (d *Desc_7A) AsvcFlag() bool {
	return d.component_flag&0b00010000 == 0x80
}

func (d *Desc_7A) Mixinfoexists() bool {
	return d.component_flag&0b00001000 == 0x80
}
func (d *Desc_7A) Substream1Flag() bool {
	return d.component_flag&0b00000100 == 0x80
}
func (d *Desc_7A) Substream2Flag() bool {
	return d.component_flag&0b00000010 == 0x80
}
func (d *Desc_7A) Substream3Flag() bool {
	return d.component_flag&0b00000001 == 0x80
}
