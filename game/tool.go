package game

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/game/db_service"
	"game_server/game/model"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const HOUSE_MAX_NUM = 1000 //小区最大房子数

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//添加世界地图坐标
func write_world_map(fileName string, building_type string) {
	logger.Debugf("write_world_map in")

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Errorf("读取配置文件错误")
		return
	}

	var mapary []interface{}
	body := bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
	err = json.Unmarshal(body, &mapary)
	if err != nil {
		logger.Errorf(" parse json error ", err)
		return
	}

	for _, value := range mapary {
		logger.Errorf(" value:%+v", value)
		data := value.(map[string]interface{})
		// var service db_service.WorldMap
		data_map := model.WorldMap{
			PositionX:    int(data["x"].(float64)),
			PositionY:    int(data["y"].(float64)),
			SmallType:    data["type"].(string),
			BuildingName: data["building_name"].(string),
			ShopName:     data["shop_name"].(string),
			Desc:         data["des"].(string),
			H5Url:        data["h5_url"].(string),
			WebUrl:       data["web_url"].(string),
			CanSale:      0,
			IsSale:       0,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		_, err = db_service.WorldMapIns.Add(3302, &data_map)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
	}
}

//添加小区坐标
func write_house(fileName string) {
	logger.Debugf("write_house in")

	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		dataAry := strings.Split(line, ",")
		if len(dataAry) != 2 {
			continue
		}

		// var service db_service.House
		for i := 0; i < HOUSE_MAX_NUM; i++ {
			Num := fmt.Sprintf("%04d", i+1)
			data_map := model.House{
				Position_x: dataAry[0],
				Position_y: dataAry[1],
				House_seq:  Num,
			}
			db_service.HouseIns.Add(&data_map)
		}
		if err != nil {
			if err == io.EOF {
				logger.Infof("File read ok!")
				break
			} else {
				logger.Errorf("Read file error!", err)
				return
			}
		}
	}
	logger.Debugf("write_house end")
}

//生成前端需要的地图数据
func write_file(fileName string) {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	wf, err := os.Create("E:\\gowork\\src\\game_server\\data.txt")
	for {

		line, _ := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		dataAry := strings.Split(line, ",")
		if len(dataAry) != 3 || dataAry[2] == "" {
			continue
		}

		if utils.IsNum(dataAry[2]) {
			continue
		}

		m_data := make(map[string]interface{})
		x, _ := strconv.Atoi(dataAry[0])
		m_data["x"] = x
		y, _ := strconv.Atoi(dataAry[1])
		m_data["y"] = y
		m_data["type"] = dataAry[2]
		//	buf, _ := json.Marshal(m_data)
		wf.Write([]byte(line))
		wf.Write([]byte("\n"))

		data_map := model.WorldMap{
			PositionX: x,
			PositionY: y,
			SmallType: dataAry[2],
			CanSale:   0,
		}
		db_service.WorldMapIns.Add(3302, &data_map)
	}
	// wf.Close()
}

//初始化城市的小区表和世界地图坐标表
func Init_db() {
	//添加世界地图海口空地坐标
	// var service db_service.House
	fileName := "./game_building_3302.txt"
	write_world_map(fileName, "0")
	//添加世界地图海口小区坐标
	//fileName = "E:\\game_doc\\haikou_xiaoqu.txt"
	//write_house(fileName)

	//更新深圳城市测试

	//查询深圳城市指定坐标测试
	/*
		data, _ := service.GetData("3269", 8, 21, "0002")

		for _, v := range data {
			a, _ := json.Marshal(v)
			fmt.Printf("getdata=%+v\n", string(a))
		}
	*/

}
