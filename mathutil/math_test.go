package mathutil

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestMin(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	minValue := Min(values)
	assert.Equal(1.0, minValue)
}

func TestMinRev(t *testing.T) {
	assert := assert.New(t)
	values := []float64{10.0, 9.0, 8.0, 7.0, 6.0, 5.0, 4.0, 3.0, 2.0, 1.0}
	minValue := Min(values)
	assert.Equal(1.0, minValue)
}

func TestMinEmpty(t *testing.T) {
	assert := assert.New(t)
	min := Min([]float64{})
	assert.Zero(min)
}

func TestMax(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	maxValue := Max(values)
	assert.Equal(10.0, maxValue)
}

func TestMaxEmpty(t *testing.T) {
	assert := assert.New(t)
	max := Max([]float64{})
	assert.Zero(max)
}

func TestSum(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	sum := Sum(values)
	assert.Equal(55.0, sum)
}

func TestSumInts(t *testing.T) {
	assert := assert.New(t)
	values := []int{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	sum := SumInts(values)
	assert.Equal(55, sum)
}

func TestSumDurations(t *testing.T) {
	assert := assert.New(t)
	values := []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second}
	sum := SumDurations(values)
	assert.Equal(6*time.Second, sum)
}

func TestSumEmpty(t *testing.T) {
	assert := assert.New(t)
	sum := Sum([]float64{})
	assert.Zero(sum)
}

func TestMean(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	mean := Mean(values)
	assert.Equal(5.5, mean)
}

func TestMeanDurations(t *testing.T) {
	assert := assert.New(t)
	values := []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second}
	mean := MeanDurations(values)
	assert.Equal(2*time.Second, mean)
}

func TestMeanEmpty(t *testing.T) {
	assert := assert.New(t)
	mean := Mean([]float64{})
	assert.Zero(mean)
}

func TestMedian(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	value := Median(values)
	assert.Equal(5.5, value)
}

func TestMedianOdd(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0}
	value := Median(values)
	assert.Equal(2.0, value)
}

func TestMedianEmpty(t *testing.T) {
	assert := assert.New(t)
	median := Median([]float64{})
	assert.Zero(median)
}

func TestModeSingleInput(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0}
	value := Mode(values)
	assert.Equal([]float64{1.0}, value)
}

func TestModeAllDifferent(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	value := Mode(values)
	assert.Equal([]float64{}, value)
}

func TestModeOddLength(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 1.0, 3.0, 4.0, 5.0}
	value := Mode(values)
	assert.Equal([]float64{1.0}, value)
}

func TestModeEvenLength(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 1.0, 1.0, 2.0, 3.0, 4.0}
	value := Mode(values)
	assert.Equal([]float64{1.0}, value)
}

func TestModeEvenLengthExtra(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 1.0, 1.0, 2.0, 3.0, 4.0, 5.0, 5.0, 5.0, 5.0}
	value := Mode(values)
	assert.Equal([]float64{5.0}, value)
}

func TestModeEmpty(t *testing.T) {
	assert := assert.New(t)
	empty := Mode([]float64{})
	assert.Empty(empty)
}

func TestVar(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	value := Var(values, 0)
	assert.Equal(8.25, value)
}

func TestVarEmpty(t *testing.T) {
	assert := assert.New(t)
	zero := Var([]float64{}, 0)
	assert.Zero(zero)
}

func TestVarP(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	value := VarP(values)
	assert.Equal(8.25, value)
}

func TestVarPEmpty(t *testing.T) {
	assert := assert.New(t)
	zero := VarP([]float64{})
	assert.Zero(zero)
}

func TestVarS(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0}
	value := VarS(values)
	assert.Equal(7.5, value)
}

func TestVarSEmpty(t *testing.T) {
	assert := assert.New(t)
	zero := VarS([]float64{})
	assert.Zero(zero)
}

func TestStdDevP(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	value := StdDevP(values)
	assert.InDelta(2.872, value, 0.001)
}

func TestStdDevPZero(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	value := StdDevP(values)
	assert.Equal(0.0, value)
}

func TestStdDevPEmpty(t *testing.T) {
	assert := assert.New(t)
	zero := StdDevP([]float64{})
	assert.Zero(zero)
}

func TestStdDevS(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	value := StdDevS(values)
	assert.InDelta(3.027, value, 0.001)
}

func TestStdDevSZero(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
	value := StdDevS(values)
	assert.Equal(0.0, value)
}

func TestStdDevSEmpty(t *testing.T) {
	assert := assert.New(t)
	zero := StdDevS([]float64{})
	assert.Zero(zero)
}

func TestRoundPlaces(t *testing.T) {
	assert := assert.New(t)
	value := RoundPlaces(0.55, 1)
	assert.InDelta(0.6, value, 0.01)
}

func TestRoundDown(t *testing.T) {
	assert := assert.New(t)
	value := RoundPlaces(0.54, 1)
	assert.InDelta(0.5, value, 0.01)
}

func TestRoundNegative(t *testing.T) {
	assert := assert.New(t)

	value := RoundPlaces(-0.55, 1)
	assert.InDelta(-0.6, value, 0.01)
}

func TestPercentile(t *testing.T) {
	assert := assert.New(t)
	values := []float64{2.0, 10.0, 3.0, 5.0, 6.0, 8.0, 7.0, 9.0, 1.0, 4.0}
	value := Percentile(values, 90.0)
	assert.InDelta(9.5, value, 0.0001)
}

func TestPercentileSorted(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	value := PercentileSorted(values, 90.0)
	assert.InDelta(9.5, value, 0.0001)
}

func TestPercentileNonInteger(t *testing.T) {
	assert := assert.New(t)
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	value := Percentile(values, 92.0)
	assert.InDelta(9.0, value, 0.0001)
}

func TestPercentileEmpty(t *testing.T) {
	assert := assert.New(t)
	zero := Percentile([]float64{}, 80.0)
	assert.Zero(zero)
}

func TestInEpsilon(t *testing.T) {
	assert := assert.New(t)

	assert.True(InEpsilon(0.0, 1-1))
	assert.False(InEpsilon(0.001, 1-1))
}
