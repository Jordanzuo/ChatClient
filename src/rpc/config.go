package rpc

var (
	chatServerAddress string

	//id, name, unionId, extraMsg, sign string, partnerId, serverId int
	id string

	name string

	unionId string

	extraMsg string

	sign string

	partnerId int

	serverId int
)

func SetConfig(_chatServerAddress, _id, _name, _unionId, _extraMsg, _sign string, _partnerId, _serverId int) {
	chatServerAddress = _chatServerAddress
	id = _id
	name = _name
	unionId = _unionId
	extraMsg = _extraMsg
	sign = _sign
	partnerId = _partnerId
	serverId = _serverId
}
