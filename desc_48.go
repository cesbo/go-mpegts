package mpegts

import (
	"fmt"
	"log"
)

type Desc_48 struct {
	ServiceProviderName string
	ServiceName         string
}

func (d *Desc_48) String() string {
	return fmt.Sprintf("0x48 service_descriptor: service_provider_name= %s, service_name: %s", d.ServiceProviderName, d.ServiceName)
}

func (d *Desc_48) Encode() (desc Descriptors) {
	// 8 is the length of the fixed part of the descriptor
	// descriptor_tag = 0x48 : 8bit
	// descriptor_length : 8bit
	// service_type: 8bit
	// service_provider_name_length: 8bit
	// service_name_length: 8bit
	desc_len := 3 + len(d.ServiceProviderName) + len(d.ServiceName) //3 is the service_type + service_provider_name_length + service_name_length
	desc = make(Descriptors, desc_len+2)
	desc[0] = 0x48
	desc[1] = byte(desc_len)
	desc[2] = 0x01
	desc[3] = byte(len(d.ServiceProviderName))
	copy(desc[4:], []byte(d.ServiceProviderName))
	desc[4+len(d.ServiceProviderName)] = byte(len(d.ServiceName))
	copy(desc[5+len(d.ServiceProviderName):], []byte(d.ServiceName))
	return desc
}

func (d *Desc_48) Decode(desc Descriptors) error {
	if len(desc) >= 8 && desc[0] == 0x48 {
		d.ServiceProviderName = string(desc[4 : 4+desc[3]])
		d.ServiceName = string(desc[5+desc[3] : 5+desc[3]+desc[4+desc[3]]])
		return nil
	} else {
		log.Printf("Desc_48: invalid format")
		return ErrDescriptorFormat
	}
}
