package converters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/timurzdev/mentorship-test-task/internal/entity"
	"github.com/timurzdev/mentorship-test-task/internal/generated"
	"github.com/timurzdev/mentorship-test-task/internal/service/helpers"
)

func Test_HouseToGen(t *testing.T) {
	type testcase struct {
		name   string
		in     entity.House
		expect generated.House
	}

	now := time.Now()

	testcaces := []testcase{
		{
			name: "full model",
			in: entity.House{
				ID:        1,
				Address:   "address",
				Year:      0,
				Developer: helpers.ToPtr("developer"),
				CreatedAt: now,
				UpdatedAt: now,
			},
			expect: generated.House{
				Id:        1,
				Address:   "address",
				Year:      0,
				Developer: helpers.ToPtr("developer"),
				CreatedAt: helpers.ToPtr(now),
				UpdatedAt: helpers.ToPtr(now),
			},
		},
	}

	for _, tc := range testcaces {
		t.Run(tc.name, func(t *testing.T) {
			out := HouseToGen(tc.in)
			assert.Equal(t, tc.expect, out)
		})
	}
}

func Test_HouseFromGen(t *testing.T) {
	type testcase struct {
		name   string
		in     generated.PostHouseCreateJSONBody
		expect entity.House
	}

	testcaces := []testcase{
		{
			name: "full model",
			in: generated.PostHouseCreateJSONBody{
				Address:   "address",
				Year:      0,
				Developer: helpers.ToPtr("developer"),
			},
			expect: entity.House{
				Address:   "address",
				Year:      0,
				Developer: helpers.ToPtr("developer"),
			},
		},
		{
			name: "only required fields",
			in: generated.PostHouseCreateJSONBody{
				Address: "address",
				Year:    0,
			},
			expect: entity.House{
				Address: "address",
				Year:    0,
			},
		},
	}

	for _, tc := range testcaces {
		t.Run(tc.name, func(t *testing.T) {
			out := HouseFromGen(tc.in)
			assert.Equal(t, tc.expect, out)
		})
	}
}
