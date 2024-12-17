package music

import (
	"ChatZoo/common/cfg"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"os"
	"time"
)

// 可以暂停、循环播放、加速播放、调音量
// todo server.json自动生成代码 枚举
// 同时播放多个音频会混合起来播放

type IMusicPlayer interface {
	Play(typ string)
}

type _MusicPlay struct {
	musicMap map[string]*beep.Buffer
}

const (
	MusicType_None = iota
)

var SrvMusicPlayer IMusicPlayer

func init() {
	config := make(map[string]string)
	cfg.ReadAnything("./music.json", &config)
	SrvMusicPlayer = newMusicPlay(config)
}

func PlayMusic(typ string) {
	SrvMusicPlayer.Play(typ)
}

func newMusicPlay(musicType map[string]string) IMusicPlayer {
	ret := &_MusicPlay{
		musicMap: make(map[string]*beep.Buffer),
	}
	ret.init(musicType)
	return ret
}

func (m *_MusicPlay) init(musicType map[string]string) {
	for typ, filePath := range musicType {
		bf, err := getMusicBuffer(filePath)
		if err != nil {
			fmt.Printf("getMusicBuffer err:%v, typ:%v filePath:%v\n", err, typ, filePath)
			return
		}
		m.musicMap[typ] = bf
	}
}

func (m *_MusicPlay) Play(typ string) {
	buffer, ok := m.musicMap[typ]
	if !ok {
		fmt.Printf("play failed, typ illegal:%v\n", typ)
		return
	}
	shot := buffer.Streamer(0, buffer.Len())
	done := make(chan bool)
	speaker.Play(beep.Seq(shot, beep.Callback(func() { // 异步回调
		done <- true
	})))
	<-done
}

func getMusicBuffer(filePath string) (*beep.Buffer, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("os.open err:%v", err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("mp3.decode err:%v", err)
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer) // 将数据都传到Buff中保存
	streamer.Close()
	return buffer, nil
}

//////////////  教程文档 ////////////////////////

// PlayMusicFromDisk 从磁盘中加载音频文件并播放 (常用于只播放一次且文件很大的情况)
func PlayMusicFromDisk(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("os open file err ", err, " filePath ", filePath)
		return
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Println("mp3.decode ", err)
		return
	}
	defer streamer.Close() // 调用这个就播放结束
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() { // 异步回调
		done <- true
	})))
	<-done
}

// PlayMusicFromMemory 将文件放在内存中, 可用于反复播放, 文件得小
func PlayMusicFromMemory(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("os open file err ", err, " filePath ", filePath)
		return
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Println("mp3.decode ", err)
		return
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer) // 将数据都传到Buff中保存
	streamer.Close()

	for {
		shot := buffer.Streamer(0, buffer.Len())
		speaker.Play(shot)
		time.Sleep(3 * time.Second)
	}
}
