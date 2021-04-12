package proto

//获取建筑简介
type C2SGetBuildingInfo struct {
	LocationID int32  `json:"locationId"` //位置ID
	LocationX  int32  `json:"x"`          //X坐标
	LocationY  int32  `json:"y"`          //Y坐标
	Token      string `json:"token"`
}

type S2CGetBuildingInfo struct {
	S2CCommon
	Desc             string `json:"desc"`
	H5Url            string `json:"h5_url"`
	WebUrl           string `json:"web_url"`
	PassportAviable  string `json:"passportAviable"`
	ImageUrl         string `json:"imageUrl"`
	SmallType        string `json:"smallType"`
	BuildingName     string `json:"buildingName"`
	BuildingTypeName string `json:"buildingTypeName"`
}

type C2SQueryShop struct {
	KeyWord string `json:"keyWord"`
	Token   string `json:"token"`
}

//商家信息项
type ShopItem struct {
	Id              string `json:"id"`
	PositionX       int    `json:"position_x"`
	PositionY       int    `json:"position_y"`
	SmallType       string `json:"small_type"`
	ShopName        string `json:"shop_name"`     //商家名称
	BuildingName    string `json:"building_name"` //建筑名
	Desc            string `json:"desc"`
	H5Url           string `json:"h5_url"`
	WebUrl          string `json:"web_url"`
	PassportAviable string `json:"passport_aviable"` //可用通证类型
	ImageUrl        string `json:"image_url"`
}

type S2CQueryShop struct {
	S2CCommon
	Data []ShopItem `json:"data"`
}
