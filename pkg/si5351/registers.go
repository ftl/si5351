package si5351

// All registers of the Si5351.
const (
	RegDeviceStatus                   = 0
	RegInterruptStatusSticky          = 1
	RegInterruptStatusMask            = 2
	RegOutputEnableControl            = 3
	RegOebPinEnableControl            = 9
	RegPLLInputSource                 = 15
	RegClk0Control                    = 16
	RegClk1Control                    = 17
	RegClk2Control                    = 18
	RegClk3Control                    = 19
	RegClk4Control                    = 20
	RegClk5Control                    = 21
	RegClk6Control                    = 22
	RegClk7Control                    = 23
	RegClk3_0DisableState             = 24
	RegClk7_4DisableState             = 25
	RegPLLAMultisynthParameters       = 26
	RegPLLBMultisynthParameters       = 34
	RegMultisynth0Parameters          = 42
	RegMultisynth1Parameters          = 50
	RegMultisynth2Parameters          = 58
	RegMultisynth3Parameters          = 66
	RegMultisynth4Parameters          = 74
	RegMultisynth5Parameters          = 82
	RegMultisynth6Parameters          = 90
	RegMultisynth7Parameters          = 91
	RegClock6_7OutputDivider          = 92
	RegClk0InitialPhaseOffset         = 165
	RegClk1InitialPhaseOffset         = 166
	RegClk2InitialPhaseOffset         = 167
	RegClk3InitialPhaseOffset         = 168
	RegClk4InitialPhaseOffset         = 169
	RegClk5InitialPhaseOffset         = 170
	RegPLLReset                       = 177
	RegCrystalInternalLoadCapacitance = 183
)
