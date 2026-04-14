package optimizations

import (
	"labs/charting"
	"math"
	"math/rand/v2"
)

const (
	maxIter = 100
)

func dichotomicSearch(f func(float64) float64, a, b, tol float64) float64 {
	eps := tol / 2
	for range maxIter {
		if math.Abs(a-b) < eps {
			return (a + b) / 2
		}
		p := (a+b)/2 + eps
		q := (a+b)/2 - eps

		if f(p) < f(q) {
			a = q
		} else {
			b = p
		}
	}
	return (a + b) / 2
}

// randomSearchNdim finds the extremum for a function with N variables.
// bounds is an Nx2 slice where each element is [min, max] for that dimension.
func randomSearchNdim(f func(...float64) float64, nSamples int, bounds [][]float64) (minPoint, maxPoint []float64) {
	dims := len(bounds)
	minVal := math.MaxFloat64
	maxVal := -math.MaxFloat64
	minPoint = make([]float64, dims)
	maxPoint = make([]float64, dims)

	for range nSamples {
		point := make([]float64, dims)
		for d := range dims {
			point[d] = bounds[d][0] + rand.Float64()*(bounds[d][1]-bounds[d][0])
		}

		val := f(point...)
		if val < minVal {
			minVal = val
			copy(minPoint, point)
		}
		if val > maxVal {
			maxVal = val
			copy(maxPoint, point)
		}
	}
	return minPoint, maxPoint
}

// In the lab instructions there is an error in the slope calculation, which is fixed here
//
// The error is that the slope is calculated as (f(X+eps) - f(X+eps)) / eps
//
// It will always be zero
func fastDescentNdim(f func(...float64) float64, initPos []float64, eps, tol, learningRate float64) []float64 {
	dims := len(initPos)
	pos := make([]float64, dims)
	copy(pos, initPos)

	for range maxIter {
		grad := make([]float64, dims)
		var gradNormSq float64

		// Calculate gradient vector
		for i := range dims {
			original := pos[i]
			pos[i] = original + eps
			fPlus := f(pos...)
			pos[i] = original - eps
			fMinus := f(pos...)
			pos[i] = original // restore

			grad[i] = (fPlus - fMinus) / (2 * eps)
			gradNormSq += grad[i] * grad[i]
		}

		if math.Sqrt(gradNormSq) < tol {
			return pos
		}

		// Trial step
		newPos := make([]float64, dims)
		newGrad := make([]float64, dims)
		var dotProduct float64

		for i := range dims {
			newPos[i] = pos[i] - grad[i]*learningRate
		}

		// Calculate gradient at new position for adjustment
		for i := range dims {
			original := newPos[i]
			newPos[i] = original + eps
			fPlus := f(newPos...)
			newPos[i] = original - eps
			fMinus := f(newPos...)
			newPos[i] = original

			newGrad[i] = (fPlus - fMinus) / (2 * eps)
			dotProduct += grad[i] * newGrad[i]
		}

		if dotProduct < 0 {
			learningRate *= 0.5
		} else {
			learningRate *= 1.6
		}

		copy(pos, newPos)
	}
	return pos
}

func RenderOneDimOptimization(req *charting.RenderRequest) (res *charting.RenderResponse) {
	return res.NewError("RenderOneDimOptimization not implemented")
}

func RenderTwoDimOptimization(req *charting.RenderRequest) (res *charting.RenderResponse) {
	return res.NewError("RenderTwoDimOptimization not implemented")
}
