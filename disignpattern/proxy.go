package disignpattern

import "fmt"

type SensitiveInfo interface {
	GetInfo() string
}

type RealSensitiveInfo struct {
}

func (r *RealSensitiveInfo) GetInfo() string {
	return "sensitive info"
}

type ProtecedSensitiveInfoProxy struct {
	realInfo *RealSensitiveInfo
	auth     bool
}

func NewProtecedSensitiveInfoProxy() *ProtecedSensitiveInfoProxy {
	return &ProtecedSensitiveInfoProxy{
		realInfo: &RealSensitiveInfo{},
		auth:     false,
	}
}

func (p *ProtecedSensitiveInfoProxy) GetInfo() string {
	if p.auth {
		return p.realInfo.GetInfo()
	}
	return "Access denied. Please authorize first"
}

func (p *ProtecedSensitiveInfoProxy) Authorize() error {
	p.auth = true
	return nil
}

func ProxyExample() {
	info := NewProtecedSensitiveInfoProxy()
	fmt.Println(info.GetInfo())
	info.Authorize()
	fmt.Println(info.GetInfo())
}
