package mem

// Memory operation permissions
// Permissions can be joined via | (p := PermRead | PermWrite)
// Permissions can be checked via & (ok := p & PermWrite)
type Permissions int

var (
	// basic permissions

	PermRead    Permissions = 0b001
	PermWrite   Permissions = 0b010
	PermExecute Permissions = 0b100

	// popular permissions

	PermReadWrite        = PermRead | PermWrite
	PermReadExecute      = PermRead | PermExecute
	PermReadWriteExecute = PermRead | PermWrite | PermExecute
)

func (p Permissions) String() string {
	permStr := []byte{'-', '-', '-'}

	if p.Has(PermRead) {
		permStr[0] = 'r'
	}
	if p.Has(PermWrite) {
		permStr[1] = 'w'
	}
	if p.Has(PermExecute) {
		permStr[2] = 'x'
	}

	return string(permStr)
}

// Check if this permissions has permissions
func (p Permissions) Has(permissions Permissions) bool {
	return p&permissions != 0
}
