package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"math"
	"math/rand"
)

type Particle struct {
	position       []float64 // Parçacığın mevcut pozisyonu
	velocity       []float64 // Parçacığın mevcut hızı
	bestPosition   []float64 // Parçacığın en iyi konumu (pBest)
	currentFitness float64   // Parçacığın mevcut uygunluk değeri
	bestFitness    float64   // Parçacığın en iyi uygunluk değeri (pBest)
}

// Fonksiyonun minimize edilmesi gerektiği varsayılan fonksiyon
func objectiveFunction(x []float64) float64 {
	// Burada fonksiyonunuzu tanımlayabilirsiniz
	// Verilen formül gibi bir fonksiyonu burada kullanabilirsiniz
	// Örneğin, verilen formüldeki gibi bir fonksiyon:
	sumX := 0.0
	sumSin := 0.0
	for _, val := range x {
		sumX += math.Abs(val)
		sumSin += math.Sin(val * val)
	}
	return sumX * math.Exp(-sumSin)
}

func main() {
	rand.Seed(42) // Rastgelelik için bir tohum (seed) ayarlayın

	// PSO parametreleri
	swarmSize := 300    // Parçacık sayısı
	dimensions := 2     // Problemin boyutu (x'in boyutu)
	maxIterations := 50 // Maksimum iterasyon sayısı
	w := 0.5            // Hız faktörü
	c1 := 2.0           // Bireysel öğrenme katsayısı
	c2 := 2.0           // Sosyal öğrenme katsayısı

	// Parçacıkları oluştur
	swarm := make([]Particle, swarmSize)
	for i := 0; i < swarmSize; i++ {
		swarm[i].position = make([]float64, dimensions)
		swarm[i].velocity = make([]float64, dimensions)
		swarm[i].bestPosition = make([]float64, dimensions)

		// Parçacıkların başlangıç pozisyonları rastgele atanır (-2𝜋 ile 2𝜋 arasında)
		for j := 0; j < dimensions; j++ {
			swarm[i].position[j] = -2*math.Pi + rand.Float64()*(2*math.Pi)
			swarm[i].velocity[j] = rand.Float64()
			swarm[i].bestPosition[j] = swarm[i].position[j]
		}

		// Parçacık için başlangıç uygunluk değeri hesaplanır
		swarm[i].currentFitness = objectiveFunction(swarm[i].position)
		swarm[i].bestFitness = swarm[i].currentFitness
	}

	// PSO Algoritması
	for iter := 0; iter < maxIterations; iter++ {
		for i := 0; i < swarmSize; i++ {
			// Parçacıkların hareketi
			for j := 0; j < dimensions; j++ {
				r1, r2 := rand.Float64(), rand.Float64()

				swarm[i].velocity[j] = w*swarm[i].velocity[j] +
					c1*r1*(swarm[i].bestPosition[j]-swarm[i].position[j]) +
					c2*r2*(swarm[i].bestPosition[j]-swarm[i].position[j])

				swarm[i].position[j] += swarm[i].velocity[j]
			}

			// Parçacığın uygunluk değerini güncelle
			swarm[i].currentFitness = objectiveFunction(swarm[i].position)

			// pBest (en iyi pozisyon) güncellemesi
			if swarm[i].currentFitness < swarm[i].bestFitness {
				swarm[i].bestFitness = swarm[i].currentFitness
				copy(swarm[i].bestPosition, swarm[i].position)
			}
		}
	}

	// En iyi çözümü ve uygunluk değerini bulma
	bestFitness := math.Inf(1)
	bestPosition := make([]float64, dimensions)
	for i := 0; i < swarmSize; i++ {
		if swarm[i].bestFitness < bestFitness {
			bestFitness = swarm[i].bestFitness
			copy(bestPosition, swarm[i].bestPosition)
		}
	}

	// Sonuçları yazdır
	fmt.Printf("En iyi uygunluk değeri: %f\n", bestFitness)
	fmt.Printf("En iyi pozisyon: %v\n", bestPosition)

	// Konverjans grafiğini çizme
	p := plot.New()

	points := make(plotter.XYs, swarmSize)
	for i := range points {
		points[i].X = float64(i)
		points[i].Y = swarm[i].bestFitness
	}

	line, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}
	p.Add(line)

	p.Title.Text = "Konverjans Grafiği"
	p.X.Label.Text = "Adım Sayısı"
	p.Y.Label.Text = "Maliyet"

	if err := p.Save(6*vg.Inch, 4*vg.Inch, "konverjans_grafigi.png"); err != nil {
		panic(err)
	}
	fmt.Println("\nKonverjans grafiği oluşturuldu: konverjans_grafigi.png")
}
