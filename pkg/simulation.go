package pkg

type RainSim []*Particle

func (rs *RainSim) Update(hm *HeightMap, dirty bool) error {
	nRain := RainSim{}
	for _, p := range *rs {
		if !p.Stopped {
			p.Drop(hm, .5, .1, .005, .001)
		}
		if !p.Stopped || dirty {
			nRain = append(nRain, p)
		}
	}
	*rs = nRain
	return nil
}
