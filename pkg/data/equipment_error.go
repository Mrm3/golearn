package data

var ErrLackOfRequiredFieldThreeParId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldThreeParId",
	Message:            "Missing required field ThreeParId",
	DescriptionChinese: "缺少必要的字段 ThreeParId，请查阅文档。",
}

var ErrLackOfRequiredFieldIpv4AddrManagement = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldIpv4AddrManagement",
	Message:            "Missing required field Ipv4AddrManagement",
	DescriptionChinese: "缺少必要的字段 Ipv4AddrManagement，请查阅文档。",
}

var ErrLackOfRequiredFieldIpv4AddrSSH = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldIpv4AddrSSH",
	Message:            "Missing required field Ipv4AddrSSH",
	DescriptionChinese: "缺少必要的字段 Ipv4AddrSSH，请查阅文档。",
}

var ErrLackOfRequiredFieldIpv4AddrController = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldIpv4AddrController",
	Message:            "Missing required field Ipv4AddrController",
	DescriptionChinese: "缺少必要的字段 Ipv4AddrController，请查阅文档。",
}

var ErrLackOfRequiredFieldUsername = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldUsername",
	Message:            "Missing required field Username",
	DescriptionChinese: "缺少必要的字段 Username，请查阅文档。",
}

var ErrLackOfRequiredFieldPassword = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldPassword",
	Message:            "Missing required field Password",
	DescriptionChinese: "缺少必要的字段 Password，请查阅文档。",
}

var ErrLackOfRequiredFieldInCharge = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldInCharge",
	Message:            "Missing required field InCharge",
	DescriptionChinese: "缺少必要的字段 InCharge，请查阅文档。",
}

var ErrLackOfRequiredFieldContact = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldContact",
	Message:            "Missing required field Contact",
	DescriptionChinese: "缺少必要的字段 Contact，请查阅文档。",
}

var ErrLackOfRequiredFieldLaunchTime = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldLaunchTime",
	Message:            "Missing required field LaunchTime",
	DescriptionChinese: "缺少必要的字段 LaunchTime，请查阅文档。",
}

// conflict
var ErrConflictThreeParId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ConflictThreeParId",
	Message:            "Repetitive value for ThreeParId",
	DescriptionChinese: "ThreeParId 已存在，请查阅文档。",
}

//others
var ErrInvalidLaunchTime = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidLaunchTime",
	Message:            "Parse err of LaunchTime",
	DescriptionChinese: "字段 LaunchTime解析错误，请查阅文档。",
}

var ErrInvalidThreeParId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidThreeParId",
	Message:            "The specified ThreeParId is invalid",
	DescriptionChinese: "指定的 ThreeParId 不存在",
}
