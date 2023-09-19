package service

import "testing"

func Test_m2h(t *testing.T) {
	cases := []struct {
		minutes       int
		expectedHours float32
	}{
		{
			minutes:       90,
			expectedHours: 1.5,
		},
		{
			minutes:       60,
			expectedHours: 1.0,
		},
		{
			minutes:       30,
			expectedHours: 0.5,
		},
		{
			minutes:       15,
			expectedHours: 0.25,
		},
		{
			minutes:       135,
			expectedHours: 2.25,
		},
	}

	for _, c := range cases {
		h := m2h(c.minutes)
		if h != c.expectedHours {
			t.Fatalf("expected %f, got %f", c.expectedHours, h)
		}
	}
}
