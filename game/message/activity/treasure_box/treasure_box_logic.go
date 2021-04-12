package treasure_box

//redis key
const ()

type TreasureBoxLogic struct {
}

//奖励
func (this *TreasureBoxLogic) Award(userId string) float32 {
	return 0.05
}
func NewTreasureBoxLogic() *TreasureBoxLogic {
	return &TreasureBoxLogic{}
}
