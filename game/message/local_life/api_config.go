package local_life

// 本地生活请求api列表.
const (
	// LocalLifeIndexRecommend 首页推荐.
	LocalLifeIndexRecommend = "/v1/index"
	// LocalLifeApiSearchHistory 获取历史搜索和热门搜索.
	LocalLifeApiSearchHistory = "/v1/hotel/getSearchHistory"
	// LocalLifeApiDelSearchHistory 删除历史搜索记录.
	LocalLifeApiDelSearchHistory = "/v1/hotel/delSearchHistory"
	// LocalLifeApiHotelSearch 搜索酒店.
	LocalLifeApiHotelSearch = "/v1/hotel/search"
	// LocalLifeApiHotelDetail 获取酒店详情.
	// 请求格式为：/v1/hotel/store/{storeId}, storeId为酒店id.
	LocalLifeApiHotelDetail = "/v1/hotel/store"
	// LocalLifeApiHotelRoomDetail 获取酒店房间详情.
	// 请求格式为：/v1/hotel/room/{goodsId}, goodsId为房间id.
	LocalLifeApiHotelRoomDetail = "/v1/hotel/room"
	// LocalLifeApiCityList 获取城市列表.
	LocalLifeApiCityList = "/v1/common/getCityList"
	// LocalLifeApiStoreType 查询店铺分类信息.
	LocalLifeApiStoreType = "/v1/common/getStoreType"
	// LocalLifeApiIndexSelectCity 首页切换城市.
	LocalLifeApiIndexSelectCity = "/v1/common/pickCity"

	// LocalLifeApiCategorySearch 分类搜索.
	LocalLifeApiCategorySearch = "/v2/search"
	// LocalLifeApiSearchSuggest 搜索联想.
	LocalLifeApiSearchSuggest = "/v2/search/suggest"
	// LocalLifeApiTopSearchV2 历史和热门搜索.
	LocalLifeApiTopSearchV2 = "/v2/search/history"
	// LocalLifeApiSearchDelete 删除搜索历史.
	LocalLifeApiSearchDelete = "/v2/search/history"

	// LocalLifeApiStoreCategory V2版本查询店铺分类信息.
	LocalLifeApiStoreTypeV2 = "/v2/store/type"
	// LocalLifeApiStoreHotelDetail 住宿详情.
	LocalLifeApiStoreHotelDetail = "/v2/store/detail"
	// LocalLifeApiStoreRestaurantDetail 美食详情.
	LocalLifeApiStoreRestaurantDetail = "/v2/store/detail"
	// LocalLifeApiGoodsDiscountDetail 美食优惠眷.
	LocalLifeApiGoodsDiscountDetail = "/v2/package/discount/detail"
	// LocalLifeApiGoodsRestaurantDetail 美食套餐.
	LocalLifeApiGoodsRestaurantDetail = "/v2/package/restaurant/detail"

	// LocalLifeApiGoodsHotelDetail 商品－住宿(房间详情)/团购-套餐详情/代金劵－优惠详情.
	LocalLifeApiGoodsHotelDetail = "/v2/goods/detail"

	// LocalLifeApiCity 查询城市列表.
	LocalLifeApiCity = "/v2/city"
	// LocalLifeApiCitySuggest 城市联想.
	LocalLifeApiCitySuggest = "/v2/city/suggest"
	// LocalLifeApiCityPick 切换城市.
	LocalLifeApiCityPick = "/v2/city/pick"
)
