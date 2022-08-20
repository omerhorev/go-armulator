package mem

import "fmt"

type SegmentationFaultError struct {
	Address           int64
	AccessPermissions Permissions
}

func NewSegmentationFaultError(address int64, permissions Permissions) *SegmentationFaultError {
	return &SegmentationFaultError{
		Address:           address,
		AccessPermissions: permissions,
	}
}

func (s *SegmentationFaultError) Error() string {
	return fmt.Sprintf("bad %s access to address 0x%x", s.AccessPermissions, s.Address)
}
