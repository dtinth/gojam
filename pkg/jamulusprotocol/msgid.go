//go:generate go run golang.org/x/tools/cmd/stringer -type=MsgId
package jamulusprotocol

type MsgId int

const (
	Illegal                 MsgId = 0
	Ackn                    MsgId = 1
	JittBufSize             MsgId = 10
	ReqJittBufSize          MsgId = 11
	ChannelGain             MsgId = 13
	ReqConnClientsList      MsgId = 16
	ChatText                MsgId = 18
	NetwTransportProps      MsgId = 20
	ReqNetwTransportProps   MsgId = 21
	ReqChannelInfos         MsgId = 23
	ConnClientsList         MsgId = 24
	ChannelInfos            MsgId = 25
	OpusSupported           MsgId = 26
	LicenceRequired         MsgId = 27
	VersionAndOs            MsgId = 29
	ChannelPan              MsgId = 30
	MuteStateChanged        MsgId = 31
	ClientId                MsgId = 32
	RecorderState           MsgId = 33
	ReqSplitMessSupport     MsgId = 34
	SplitMessSupported      MsgId = 35
	ClmPingMs               MsgId = 1001
	ClmPingMsWithNumClients MsgId = 1002
	ClmServerFull           MsgId = 1003
	ClmRegisterServer       MsgId = 1004
	ClmUnregisterServer     MsgId = 1005
	ClmServerList           MsgId = 1006
	ClmReqServerList        MsgId = 1007
	ClmSendEmptyMessage     MsgId = 1008
	ClmEmptyMessage         MsgId = 1009
	ClmDisconnection        MsgId = 1010
	ClmVersionAndOs         MsgId = 1011
	ClmReqVersionAndOs      MsgId = 1012
	ClmConnClientsList      MsgId = 1013
	ClmReqConnClientsList   MsgId = 1014
	ClmChannelLevelList     MsgId = 1015
	ClmRegisterServerResp   MsgId = 1016
	ClmRegisterServerEx     MsgId = 1017
	ClmRedServerList        MsgId = 1018
)
