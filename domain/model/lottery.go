package model

type Lottery struct {
	LotteryID string //
	Name      string //宝くじの名前
	Price     int64  //値段
	Number    int64  //販売個数（最大10,000枚）
}
