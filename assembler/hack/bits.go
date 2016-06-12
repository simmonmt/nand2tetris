package hack

const (
	DEST_M = 1 << iota
	DEST_D = 1 << iota
	DEST_A = 1 << iota
)

const (
	JUMP_GT = 1 << iota
	JUMP_EQ = 1 << iota
	JUMP_LT = 1 << iota
)

const (
	COMP_0      = 0x2a
	COMP_1      = 0x3f
	COMP_NEG1   = 0x3a
	COMP_D      = 0xc
	COMP_A      = 0x30
	COMP_NOTD   = 0x0d
	COMP_NOTA   = 0x31
	COMP_NEGD   = 0x0f
	COMP_NEGA   = 0x33
	COMP_DPLUS1 = 0x1f
	COMP_APLUS1 = 0x37
	COMP_DMIN1  = 0x0e
	COMP_AMIN1  = 0x32
	COMP_DPLUSA = 0x02
	COMP_DMINA  = 0x13
	COMP_AMIND  = 0x07
	COMP_DANDA  = 0x00
	COMP_DORA   = 0x15

	COMP_M      = 0x70
	COMP_NOTM   = 0x71
	COMP_NEGM   = 0x73
	COMP_MPLUS1 = 0x77
	COMP_MMIN1  = 0x72
	COMP_DPLUSM = 0x42
	COMP_DMINM  = 0x53
	COMP_MMIND  = 0x47
	COMP_DANDM  = 0x40
	COMP_DORM   = 0x55
)
