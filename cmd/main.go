package main

import "github.com/wujunwei928/edge-tts-go/internal/cmd"

func main() {
	cmd.Execute()

	//edge_tts.ListVoices("")
	//T := time.Now().UTC()
	//fmt.Println(T, T.Format("Mon Jan 02 2006 15:04:05 GMT+0000 (Coordinated Universal Time)"))

	//ssml := "白日依山尽，黄河入海流。欲穷千里目，更上一层楼。"
	//	ssml := `
	//白日登山望烽火，黄昏饮马傍交河。
	//行人刁斗风沙暗，公主琵琶幽怨多。
	//野云万里无城郭，雨雪纷纷连大漠。
	//胡雁哀鸣夜夜飞，胡儿眼泪双双落。
	//闻道玉门犹被遮，应将性命逐轻车。
	//年年战骨埋荒外，空见蒲桃入汉家。
	//`
	//	log.Println("接收到SSML(Edge):", ssml)
	//	c, err := edge_tts.NewCommunicate(
	//		ssml,
	//		edge_tts.SetVoice("zh-CN-liaoning-XiaobeiNeural"),
	//	)
	//	fmt.Println(c)
	//	fmt.Println(err)
	//	data, err := c.Stream()
	//	fmt.Println(err)
	//	os.WriteFile("zh-CN-liaoning-XiaobeiNeural.mp3", data, 0644)
}
