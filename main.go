package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"
)

const (
	V_source    = 10.0 // 電源電圧 (V)
	Resistance  = 50.0 // 抵抗値 (Ω)
	Capacitance = 0.05 // コンデンサ容量 (F)
	dt          = 0.1  // 時間ステップ (秒)
	MaxSteps    = 30   // アニメーションの総フレーム数
	Width       = 300  // 画像の横幅 (ピクセル)
	Height      = 300  // 画像の縦幅 (ピクセル)
)

func main() {
	fmt.Println("⚡ Cgo不要版：RC回路過渡現象＆電場変化シミュレーションを開始します...")

	var chargeQ float64
	var g gif.GIF // GIFアニメーション用のオブジェクト

	// タイムステップごとに電気回路と電場をシミュレート
	for step := 0; step < MaxSteps; step++ {
		// --- 1. 動電気：コンデンサの過渡現象計算 ---
		v_capacitor := chargeQ / Capacitance
		current := (V_source - v_capacitor) / Resistance
		chargeQ += current * dt // 電荷の蓄積

		// --- 2. 静電気：この瞬間の電場を画像（1フレーム）として描画 ---
		// パレット（0: 背景黒, 1〜255: 充電が進むほど鮮やかになる赤のグラデーション）
		palette := []color.Color{color.RGBA{0, 0, 0, 255}}
		for i := 1; i <= 255; i++ {
			// 充電比率（0.0〜1.0）に応じて赤の強さを変える
			ratio := v_capacitor / V_source
			r := uint8(float64(i) * ratio)
			palette = append(palette, color.RGBA{R: r, G: 30, B: 30, A: 255})
		}

		img := image.NewPaletted(image.Rect(0, 0, Width, Height), palette)

		// 擬似的な極板（左右）からの距離に応じた電場の強さを各ピクセルで計算
		for x := 0; x < Width; x++ {
			for y := 0; y < Height; y++ {
				// 中心からの距離
				dx := float64(x - Width/2)
				dy := float64(y - Height/2)
				dist := math.Sqrt(dx*dx + dy*dy)

				if dist < 10 {
					continue // 中心部は除外
				}

				// 電場は距離の2乗に反比例し、電荷量(chargeQ)に比例する
				// 画面に綺麗に収めるための擬似的な係数を掛けています
				fieldStrength := (chargeQ * 500000) / (dist * dist)
				if fieldStrength > 1.0 {
					fieldStrength = 1.0
				}

				// 電場の強さに応じてパレットの色（赤の濃さ）を割り当てる
				colorIndex := uint8(fieldStrength * 255)
				if colorIndex > 0 {
					img.SetColorIndex(x, y, colorIndex)
				}
			}
		}

		// フレームをアニメーションに追加
		g.Image = append(g.Image, img)
		g.Delay = append(g.Delay, 8) // 80ミリ秒/フレームの速度

		fmt.Printf("⏱️ 時刻: %.1f秒 | 電圧: %.2f V | 回路電流: %.4f A (フレーム %d/%d 作成中)\n",
			float64(step)*dt, v_capacitor, current, step+1, MaxSteps)
	}

	// 3. 完成したアニメーションをGIFファイルとして保存
	f, err := os.Create("simulation.gif")
	if err != nil {
		fmt.Println("ファイル作成エラー:", err)
		return
	}
	defer f.Close()

	err = gif.EncodeAll(f, &g)
	if err != nil {
		fmt.Println("GIFエンコードエラー:", err)
		return
	}

	fmt.Println("\n✅ 成功！すべての動向を含んだアニメーション 'simulation.gif' が生成されました！")
}