package pkg

import "testing"

func TestHeightMap_Set(t *testing.T) {
	hm := NewHeightMap(2)
	hm[0][0] = 10
	type args struct {
		x   int
		y   int
		val float64
	}
	tests := []struct {
		name string
		hm   *HeightMap
		args args
	}{
		{"test1", &hm, args{0, 0, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := tt.hm.Get(tt.args.x, tt.args.y)
			tt.hm.Set(tt.args.x, tt.args.y, tt.args.val)
			if o == tt.hm.Get(tt.args.x, tt.args.y) {
				t.Fail()
			}
		})
	}
}
