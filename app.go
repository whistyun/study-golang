package main

import (
	"math"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"gonum.org/v1/gonum/mat"

	"StudyGo/affine"
	"StudyGo/image"
)

var (
	screenW int = 115
	screenH int = 25
)
var (
	rotationSpeed     [3]float64
	rotationSpeedLock sync.Mutex
	screenDot         []float64 = make([]float64, screenW*screenH)
	screenThin        []float64 = make([]float64, screenW*screenH)
	screenDotLock     sync.Mutex
)

func main() {

	imgpath := "./Default.png"
	if len(os.Args) > 1 {
		imgpath = os.Args[1]
	}

	img, err := image.LoadImage(imgpath)
	if err != nil {
		panic(err)
	}

	stop := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)

	go inputRoutine(&wg)
	go logicLoop(img, stop)
	go drawLoop(img, stop)

	wg.Wait()
	println("stopped")
	close(stop)
}

func inputRoutine(wg *sync.WaitGroup) {
	var chr = make([]byte, 1)

	var rotXY, rotYZ, rotZX float64

	for {
		_, err := os.Stdin.Read(chr)

		if err != nil || chr[0] == 'c' {
			break
		}

		switch chr[0] {
		case 'w':
			rotYZ = rotYZ - math.Pi/2
		case 's':
			rotYZ = rotYZ + math.Pi/2
		case 'a':
			rotZX = rotZX + math.Pi/2
		case 'd':
			rotZX = rotZX - math.Pi/2
		case 'q':
			rotXY = rotXY - math.Pi/2
		case 'e':
			rotXY = rotXY + math.Pi/2
		}

		rotationSpeedLock.Lock()

		rotationSpeed[0] = rotXY
		rotationSpeed[1] = rotYZ
		rotationSpeed[2] = rotZX

		rotationSpeedLock.Unlock()
	}

	wg.Done()
}

func logicLoop(img *image.InkMat, stop chan struct{}) {

	var rotate [3]float64
	mixed := mat.NewDense(4, 4, nil)

	srcVec := mat.NewVecDense(4, nil)
	srcVec.SetVec(2, 0)
	srcVec.SetVec(3, 1)

	dstVec := mat.NewVecDense(4, nil)

	for {
		select {
		case <-stop:
			return
		default:
		}

		rotationSpeedLock.Lock()

		for i := 0; i < 3; i++ {
			rotate[i] += rotationSpeed[i] / 15
		}

		rotationSpeedLock.Unlock()

		scale := math.Min(math.Min(float64(screenW)/float64(img.W), float64(screenH)/float64(img.H)), 1)
		scaleMat := affine.Scale(scale, scale, 1)
		rotateMat := affine.Rotate(rotate[0], rotate[1], rotate[2])
		imgCntrMat := affine.Translate(-float64(img.W)/2, -float64(img.H)/2, 0)
		scrnCntrMat := affine.Translate(float64(screenW)/2, float64(screenH)/2, 0)

		mixed.Product(scrnCntrMat, rotateMat, scaleMat, imgCntrMat)

		screenDotLock.Lock()

		for i := 0; i < len(screenDot); i++ {
			screenDot[i] = 0
			screenThin[i] = 0
		}

		imgIdx := 0
		for i := 0; i < img.H; i++ {
			srcVec.SetVec(1, float64(i))
			for j := 0; j < img.W; j++ {
				srcVec.SetVec(0, float64(j))
				dstVec.MulVec(mixed, srcVec)

				posX := int(dstVec.At(0, 0))
				posY := int(dstVec.At(1, 0))

				if 0 <= posX && posX < screenW && 0 <= posY && posY < screenH {
					idx := posX + screenW*posY
					screenDot[idx] = (screenDot[idx]*screenThin[idx] + float64(img.P[imgIdx])) / (screenThin[idx] + 1)
					screenThin[idx]++
				}

				imgIdx++
			}
		}

		screenDotLock.Unlock()

		time.Sleep(time.Second / 15)
	}
}

func drawLoop(img *image.InkMat, stop chan struct{}) {

	grayscale := "MWN$@%#&B89EGA6mK5HRkbYT43V0JL7gpaseyxznocv?jIftr1li*=-~^`':;,. "

	var clear func()

	if runtime.GOOS == "windows" {
		clear = func() {
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
	} else {
		clear = func() {
			print("\033[H\033[2J")
		}
	}

	clear()

	charArea := make([]byte, 0, screenW*screenH)

	for {
		select {
		case <-stop:
			return
		default:
		}

		screenDotLock.Lock()

		imgIdx := 0
		charArea = charArea[:0]

		for i := 0; i < screenH; i++ {
			for j := 0; j < screenW; j++ {
				dotColor := screenDot[imgIdx]
				charArea = append(charArea, grayscale[int(dotColor)*(len(grayscale)-1)/255])
				imgIdx++
			}
			charArea = append(charArea, '\r', '\n')
		}

		clear()
		println("WASD turn up, left, down, right. EQ closewise, counter-clockwize. C close")
		println(string(charArea))

		screenDotLock.Unlock()

		time.Sleep(time.Second / 15)
	}
}
