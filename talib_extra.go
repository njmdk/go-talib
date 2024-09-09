package talib

import (
	"github.com/shopspring/decimal"
	"math"
)

func MinValue(array []float64) float64 {
	var minValue = array[0]
	for i := 0; i < len(array); i++ {
		if array[i] < minValue {
			minValue = array[i]
		}
	}
	return minValue
}
func MaxValue(array []float64) float64 {
	var minValue = array[0]
	for i := 0; i < len(array); i++ {
		if array[i] > minValue {
			minValue = array[i]
		}
	}
	return minValue
}
func ExtStochRSI(records []float64, rsiNum int, stochNum int, kPeriod int, dPeriod int) ([]float64, []float64) {
	var rsi = Rsi(records, rsiNum)
	var minRsi []float64
	var maxRsi []float64
	rsiLen := len(rsi)
	for i := 0; i < rsiLen; i++ {
		if i < stochNum*2-1 {
			minRsi = append(minRsi, 0)
			maxRsi = append(maxRsi, 0)
		} else {
			minV := MinValue(rsi[i-stochNum+1 : i+1])
			maxV := MaxValue(rsi[i-stochNum+1 : i+1])
			minRsi = append(minRsi, minV)
			maxRsi = append(maxRsi, maxV)
		}
	}
	var stoch []float64
	for i := 0; i < rsiLen; i++ {
		if maxRsi[i] == minRsi[i] {
			stoch = append(stoch, 0)
		} else {
			stochValue := decimal.NewFromFloat(100).Mul((decimal.NewFromFloat(rsi[i]).Sub(decimal.NewFromFloat(minRsi[i]))).Div(decimal.NewFromFloat(maxRsi[i]).Sub(decimal.NewFromFloat(minRsi[i]))))
			stoch = append(stoch, stochValue.InexactFloat64())
		}
	}
	k := Ma(stoch, kPeriod, SMA)
	d := Ma(k, dPeriod, SMA)

	return k, d
}
func ExtVar(inReal []float64, inTimePeriod int) []float64 {

	outReal := make([]float64, len(inReal))

	nbInitialElementNeeded := inTimePeriod - 1
	startIdx := nbInitialElementNeeded
	periodTotal1 := decimal.NewFromFloat(0.0)
	periodTotal2 := decimal.NewFromFloat(0.0)
	trailingIdx := startIdx - nbInitialElementNeeded
	i := trailingIdx
	if inTimePeriod > 1 {
		for i < startIdx {
			tempReal := inReal[i]
			//periodTotal1 += tempReal
			periodTotal1 = periodTotal1.Add(decimal.NewFromFloat(tempReal))
			//tempReal *= tempReal
			tempReal = decimal.NewFromFloat(tempReal).Mul(decimal.NewFromFloat(tempReal)).InexactFloat64()
			//periodTotal2 += tempReal
			periodTotal2 = periodTotal2.Add(decimal.NewFromFloat(tempReal))
			i++
		}
	}
	outIdx := startIdx
	for ok := true; ok; {
		tempReal := inReal[i]
		//periodTotal1 += tempReal
		periodTotal1 = periodTotal1.Add(decimal.NewFromFloat(tempReal))
		//tempReal *= tempReal
		tempReal = decimal.NewFromFloat(tempReal).Mul(decimal.NewFromFloat(tempReal)).InexactFloat64()
		//periodTotal2 += tempReal
		periodTotal2 = periodTotal2.Add(decimal.NewFromFloat(tempReal))
		//meanValue1 := periodTotal1 / float64(inTimePeriod)
		meanValue1 := periodTotal1.Div(decimal.NewFromFloat(float64(inTimePeriod)))
		//meanValue2 := periodTotal2 / float64(inTimePeriod)
		meanValue2 := periodTotal2.Div(decimal.NewFromFloat(float64(inTimePeriod)))
		tempReal = inReal[trailingIdx]
		//periodTotal1 -= tempReal
		periodTotal1 = periodTotal1.Sub(decimal.NewFromFloat(tempReal))
		//tempReal *= tempReal
		tempReal = decimal.NewFromFloat(tempReal).Mul(decimal.NewFromFloat(tempReal)).InexactFloat64()
		//periodTotal2 -= tempReal
		periodTotal2 = periodTotal2.Sub(decimal.NewFromFloat(tempReal))
		//outReal[outIdx] = meanValue2 - meanValue1*meanValue1
		outReal[outIdx] = meanValue2.Sub(meanValue1.Mul(meanValue1)).InexactFloat64()
		i++
		trailingIdx++
		outIdx++
		ok = i < len(inReal)
	}
	return outReal
}
func ExtStdDev(inReal []float64, inTimePeriod int, inNbDev float64) []float64 {

	outReal := ExtVar(inReal, inTimePeriod)

	if inNbDev != 1.0 {
		for i := 0; i < len(inReal); i++ {
			tempReal := outReal[i]
			if tempReal > 0.0 {
				outReal[i] = math.Sqrt(tempReal) * inNbDev
			} else {
				outReal[i] = 0.0
			}
		}
	} else {
		for i := 0; i < len(inReal); i++ {
			tempReal := outReal[i]
			if tempReal > 0.0 {
				outReal[i] = math.Sqrt(tempReal)
			} else {
				outReal[i] = 0.0
			}
		}
	}
	return outReal
}
func ExtBBands(inReal []float64, inTimePeriod int, inNbDevUp float64, inNbDevDn float64, inMAType MaType) ([]float64, []float64, []float64) {

	outRealUpperBand := make([]float64, len(inReal))
	outRealMiddleBand := Ma(inReal, inTimePeriod, inMAType)
	outRealLowerBand := make([]float64, len(inReal))

	tempBuffer2 := ExtStdDev(inReal, inTimePeriod, 1.0)

	if inNbDevUp == inNbDevDn {

		if inNbDevUp == 1.0 {
			for i := 0; i < len(inReal); i++ {
				tempReal := tempBuffer2[i]
				tempReal2 := outRealMiddleBand[i]
				outRealUpperBand[i] = decimal.NewFromFloat(tempReal2).Add(decimal.NewFromFloat(tempReal)).InexactFloat64()
				outRealLowerBand[i] = decimal.NewFromFloat(tempReal2).Sub(decimal.NewFromFloat(tempReal)).InexactFloat64()
			}
		} else {
			for i := 0; i < len(inReal); i++ {
				tempReal := decimal.NewFromFloat(tempBuffer2[i]).Mul(decimal.NewFromFloat(inNbDevUp))
				tempReal2 := outRealMiddleBand[i]
				outRealUpperBand[i] = decimal.NewFromFloat(tempReal2).Add(tempReal).InexactFloat64()
				outRealLowerBand[i] = decimal.NewFromFloat(tempReal2).Sub(tempReal).InexactFloat64()
			}
		}
	} else if inNbDevUp == 1.0 {
		for i := 0; i < len(inReal); i++ {
			tempReal := tempBuffer2[i]
			tempReal2 := outRealMiddleBand[i]
			outRealUpperBand[i] = decimal.NewFromFloat(tempReal2).Add(decimal.NewFromFloat(tempReal)).InexactFloat64()
			outRealLowerBand[i] = decimal.NewFromFloat(tempReal2).Sub(decimal.NewFromFloat(tempReal).Mul(decimal.NewFromFloat(inNbDevDn))).InexactFloat64()
		}
	} else if inNbDevDn == 1.0 {
		for i := 0; i < len(inReal); i++ {
			tempReal := tempBuffer2[i]
			tempReal2 := outRealMiddleBand[i]
			outRealLowerBand[i] = decimal.NewFromFloat(tempReal2).Sub(decimal.NewFromFloat(tempReal)).InexactFloat64()
			outRealUpperBand[i] = decimal.NewFromFloat(tempReal2).Add(decimal.NewFromFloat(tempReal).Mul(decimal.NewFromFloat(inNbDevUp))).InexactFloat64()
		}
	} else {
		for i := 0; i < len(inReal); i++ {
			tempReal := tempBuffer2[i]
			tempReal2 := outRealMiddleBand[i]
			outRealUpperBand[i] = decimal.NewFromFloat(tempReal2).Add(decimal.NewFromFloat(tempReal).Mul(decimal.NewFromFloat(inNbDevUp))).InexactFloat64()
			outRealLowerBand[i] = decimal.NewFromFloat(tempReal2).Sub(decimal.NewFromFloat(tempReal).Mul(decimal.NewFromFloat(inNbDevDn))).InexactFloat64()
		}
	}
	return outRealUpperBand, outRealMiddleBand, outRealLowerBand
}
