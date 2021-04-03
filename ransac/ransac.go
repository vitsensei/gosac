package ransac

import (
	"fmt"
	"math/rand"
	"os"
)

// This implementation is based on https://en.wikipedia.org/wiki/Random_sample_consensus
// 							----------------------------
// Do this for X time:
// 	Step 1: Randomly sample X number of data and fit it to a model M
//  Step 2: For other points, try to fit into M. Not all of them will be considered
//   		to be part of the model. Ransac assume there will be outliers and that's ok.
//  Step 3: If M has more than D number of inliers, this is a good model and will be considered
//  		to be the ultimate best model
//  Step 4: We compare model M, and the best model so far through some sort of Loss() function.
//  		This can be anything, such as RMS or Sum of all loss.

type Model interface {
	FitData(p Population, ind []int) // Take in a set of data and fit the model
	SingleLoss(p Population, ind int) float64
	Loss(p Population, ind []int) float64 // Take in one data point and calculate the loss
	Clone() Model
}

type Population struct {
	Data []Data
}

func (p Population) ToFile(path string) {
	f, err := os.Create(path)
	must(err)
	defer f.Close()

	for _, d := range p.Data {
		d.ToFile(f)
	}
}

type Data struct {
	NumericFields []float64
}

func (d Data) ToFile(f *os.File) {
	for i := range d.NumericFields {
		if i == len(d.NumericFields) - 1 {
			_, _ = f.WriteString(fmt.Sprintf("%f\n", d.NumericFields[i]))
		} else {
			_, _ = f.WriteString(fmt.Sprintf("%f, ", d.NumericFields[i]))
		}
	}
}

type Ransac struct {
	baseModel Model // Self-explained
	population Population // The whole population
	numberOfSubSample int // The number of samples considered for model fitting each iteration
	iterations int // Maximum iteration before the model was forced to stop
	thresholdLoss float64 // Only point with loss x < thresholdLoss be accepted as the
	// inliers of the model
	thresholdPoint int // Only model with number of inliers d > thresholdPoint be considered
	// to be bestModel
}

func NewRansac(iterations int, thresholdLoss float64, thresholdPoint int, numberOfSubSample int) (Ransac, error) {
	newRansac := Ransac{
		iterations: iterations,
		thresholdLoss: thresholdLoss,
		thresholdPoint: thresholdPoint,
		numberOfSubSample: numberOfSubSample,
		baseModel: nil,
		population: Population{},
	}

	return newRansac, nil
}

func (r *Ransac) SetModel(m Model) {
	r.baseModel = m
}

func (r *Ransac) SetPopulation(p Population) {
	r.population = p
}

func (r *Ransac) Run() (Model, float64) {
	var bestModel Model
	var bestLoss float64

	for i := 0; i <  r.iterations; i++ {
		// Step 1.a: Randomly select numberOfSubSample sample
		populationSize := len(r.population.Data)
		subset, indMap := generateNUniqueNumber(r.numberOfSubSample, populationSize)

		// Step 1.b: Use the selected subset to fit the model
		r.baseModel.FitData(r.population, subset)

		// Step 2: For each of the (other) points, try to calculate the loss
		for j := 0; j < populationSize; j++ {
			if indMap[j] {
				continue
			}

			// Calculate the loss for other point. If the loss is smaller than
			// the threshold, include the other point into the subset
			if r.baseModel.SingleLoss(r.population, j) < r.thresholdLoss {
				subset = append(subset, j)
			}
		}

		// If there are more than a certain amount of inliers, consider the model
		// to take the throne
		if len(subset) > r.thresholdPoint {
			// Fit the subset to the model again. This time, the subset contains
			// more inliers
			currentLoss := r.baseModel.Loss(r.population, subset)

			// Compare current model with the best model and update if necessary
			if bestModel == nil {
				bestModel = r.baseModel.Clone()
				bestLoss = currentLoss
			} else if currentLoss < bestLoss {
				bestModel = r.baseModel.Clone()
				bestLoss = currentLoss
			}
		}
	}

	return bestModel, bestLoss
}

func generateNUniqueNumber(n, max int) ([]int, map[int]bool) {
	generated := make(map[int]bool)

	for len(generated) < n {
		i := rand.Intn(max)
		if !generated[i] {
			generated[i] = true
		}
	}

	keys := make([]int, 0, len(generated))
	for k := range generated {
		keys = append(keys, k)
	}

	return keys, generated
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}