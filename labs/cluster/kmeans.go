package cluster

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

const (
	KmeansGraphID = "kmeans"
	KmeansChartID = "kmeans"

	VariableCenroidsChoiseOptionID = "centroids-choise-method"
	VariableNumCentriodsID         = "centroids-num"
)

type scatterPoint struct {
	X, Y float64
}

func kmeans(points *Points, num_centroids, maxIter uint32) ([]int, []scatterPoint, error) {
	centroids := make([]scatterPoint, num_centroids)
	for i := range num_centroids {
		rand_idx := rand.Intn(len(points.X))
		centroids[i] = scatterPoint{X: points.X[rand_idx], Y: points.Y[rand_idx]}
	}

	var labels []int
	for i := range maxIter {
		labels = assignPointsToCentroids(points, centroids)
		newCentroids := updateCentroids(points, labels, num_centroids)

		if fmt.Sprint(newCentroids) == fmt.Sprint(centroids) {
			break
		}
		centroids = newCentroids
		if i == maxIter-1 {
			return nil, nil, errors.New("k-means: centroids haven't coverged")
		}
	}

	return labels, centroids, nil
}

func assignPointsToCentroids(points *Points, centroids []scatterPoint) []int {
	labels := make([]int, len(points.X))
	for i := range points.X {
		best, minDist := 0, math.MaxFloat64
		for j, c := range centroids {
			if d := eclidianDist(scatterPoint{X: points.X[i], Y: points.Y[i]}, c); d < minDist {
				minDist, best = d, j
			}
		}
		labels[i] = best
	}
	return labels
}

func updateCentroids(points *Points, labels []int, k uint32) []scatterPoint {
	sums := make([]scatterPoint, k)
	counts := make([]int, k)
	for i := range points.X {
		sums[labels[i]].X += points.X[i]
		sums[labels[i]].Y += points.Y[i]
		counts[labels[i]]++
	}
	centroids := make([]scatterPoint, k)
	for i := range centroids {
		if counts[i] > 0 {
			centroids[i] = scatterPoint{sums[i].X / float64(counts[i]), sums[i].Y / float64(counts[i])}
		}
	}
	return centroids
}

func eclidianDist(a, b scatterPoint) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}
