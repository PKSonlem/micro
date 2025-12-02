package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/generated"
	"github.com/timurzdev/mentorship-test-task/internal/service/helpers"
)

func Test_FaltToGen(t *testing.T) {
	type testcase struct {
		name   string
		in     entity.Flat
		expect generated.Flat
	}

	testcases := []testcase{
		{
			name: "full model",
			in: entity.Flat{
				ID:         1,
				HouseID:    1,
				FlatNumber: 1,
				Price:      400,
				Rooms:      3,
				Status:     "created",
			},
			expect: generated.Flat{
				Id:      1,
				HouseId: 1,
				Price:   400,
				Rooms:   3,
				Status:  "created",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			out := FlatToGen(tc.in)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func Test_FlatFromGen(t *testing.T) {
	type testcace struct {
		name   string
		in     generated.PostFlatCreateJSONBody
		expect entity.Flat
	}

	testcases := []testcace{
		{
			name: "full model",
			in: generated.PostFlatCreateJSONBody{
				HouseId: 1,
				Price:   400,
				Rooms:   helpers.ToPtr(3),
			},
			expect: entity.Flat{
				HouseID: 1,
				Price:   400,
				Rooms:   3,
			},
		},
		{
			name: "only required fields",
			in: generated.PostFlatCreateJSONBody{
				HouseId: 1,
				Price:   400,
			},
			expect: entity.Flat{
				HouseID: 1,
				Price:   400,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			out := FlatFromGenCreate(tc.in)
			assert.Equal(t, tc.expect, out)
		})
	}
}
